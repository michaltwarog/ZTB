package main

import (
	"database-tester/types"
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	_ "github.com/sijms/go-ora/v2"
)

func GetORCLDBStorageManager(dbParams map[string]string) *ORCLDBStorageManager {
	db, err := sql.Open("oracle", fmt.Sprintf("oracle://%s:%s@%s:%s/%s",
		dbParams["username"],
		dbParams["password"],
		dbParams["server"],
		dbParams["port"],
		dbParams["service"],
	))
	if err != nil {
		panic(fmt.Errorf("error in sql.Open: %w", err))
	}

	err = db.Ping()
	if err != nil {
		panic(fmt.Errorf("error pinging db: %w", err))
	}
	print("Connected to Oracle DB\n")
	return &ORCLDBStorageManager{DB: db}
}

func convertDateStringToOracle(dateString string) interface{} {
	return fmt.Sprintf("TO_DATE('%s', 'YYYY-MM-DD HH24:MI:SS')", dateString)
}

// func convertToOracleType(value interface{}, fieldType reflect.StructField) interface{} {
// 	switch fieldType.Type.Kind() {
// 	case reflect.String:
// 		if fieldType.Name == "Date_Of_Creation" || fieldType.Name == "Date_Of_Modification" || fieldType.Name == "Date_Of_Upload" {
// 			return convertDateStringToOracle(value.(string))
// 		}
// 		return value
// 	case reflect.Int, reflect.Int32, reflect.Int64, reflect.Float32, reflect.Float64:
// 		return value
// 	default:
// 		return value
// 	}
// }

func convertToOracleType(value interface{}, fieldType reflect.StructField) interface{} {
	switch fieldType.Type.Kind() {
	case reflect.String:
		if fieldType.Name == "Date_Of_Creation" || fieldType.Name == "Date_Of_Modification" {
			return convertDateStringToOracle(value.(string))
		}
		return value
	case reflect.Int, reflect.Int32, reflect.Int64, reflect.Float32, reflect.Float64:
		return value
	default:
		return value
	}
}

func Insert(db *sql.DB, tableName string, data interface{}) (sql.Result, error) {

	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	var columns []string
	var placeholders []string
	var values []interface{}

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		dbTag := field.Tag.Get("db")
		if dbTag == "" {
			dbTag = strings.ToUpper(field.Name)
		}
		if dbTag != "ID" {
			columns = append(columns, dbTag)
			placeholders = append(placeholders, fmt.Sprintf(":%d", i+1))
			convertedValue := convertToOracleType(val.Field(i).Interface(), field)
			values = append(values, convertedValue)
		}
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		tableName,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "))

	fmt.Println(query)
	return db.Exec(query, values...)
}

func Patch(db *sql.DB, tableName string, id int, data interface{}) error {
	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	var setClauses []string
	var values []interface{}

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		dbTag := field.Tag.Get("db")
		if dbTag == "" {
			dbTag = strings.ToUpper(field.Name)
		}
		if dbTag != "ID" {
			setClause := fmt.Sprintf("%s = :%s", dbTag, dbTag)
			setClauses = append(setClauses, setClause)
			convertedValue := convertToOracleType(val.Field(i).Interface(), field)
			values = append(values, convertedValue)
		}
	}

	query := fmt.Sprintf("UPDATE %s SET %s WHERE ID = %d",
		tableName,
		strings.Join(setClauses, ", "),
		id,
	)

	_, err := db.Exec(query, values...)
	return err
}

func Delete(db *sql.DB, tableName string, id int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE ID = :id", tableName)
	_, err := db.Exec(query, sql.Named("id", id))
	return err
}

func Get(db *sql.DB, tableName string, id int) (map[string]interface{}, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE ID = %d FETCH FIRST 1 ROWS ONLY", tableName, id)

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("error fetching row: %w", err)
		}
		return nil, fmt.Errorf("no rows were returned")
	}

	cols, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("error getting columns: %w", err)
	}

	colVals := make([]interface{}, len(cols))
	colPointers := make([]interface{}, len(cols))
	for i := range colVals {
		colPointers[i] = &colVals[i]
	}

	if err := rows.Scan(colPointers...); err != nil {
		return nil, fmt.Errorf("error scanning row: %w", err)
	}

	result := make(map[string]interface{})
	for i, colName := range cols {
		var val interface{}
		byteVal, ok := colVals[i].([]byte)
		if ok {
			val = string(byteVal)
		} else {
			val = colVals[i]
		}
		result[colName] = val
	}

	return result, nil
}

type ORCLDBStorageManager struct {
	DB *sql.DB
}

func getInsertedID(row *sql.Row) (int, error) {
	var newID int
	if err := row.Scan(&newID); err != nil {
		return 0, err
	}
	return newID, nil
}

func (sm *ORCLDBStorageManager) InsertUser(user types.User) (int, error) {
	result, err := Insert(sm.DB, `"USER"`, user)
	if err != nil {
		return 0, fmt.Errorf("error inserting user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error getting last insert ID: %w", err)
	}

	return int(id), nil
}

func (sm *ORCLDBStorageManager) InsertNote(note types.Note) (int, error) {
	tableName := "NOTE"
	result, err := Insert(sm.DB, tableName, note)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (sm *ORCLDBStorageManager) PatchUser(user types.User) (int, error) {
	if user.ID == 0 {
		return 0, fmt.Errorf("user ID is required")
	}
	err := Patch(sm.DB, `"USER"`, user.ID, user)
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}

func (sm *ORCLDBStorageManager) PatchNote(note types.Note) (int, error) {
	if note.ID == 0 {
		return 0, fmt.Errorf("note ID is required")
	}
	err := Patch(sm.DB, "NOTE", note.ID, note)
	if err != nil {
		return 0, err
	}
	return note.ID, nil
}

func (sm *ORCLDBStorageManager) DeleteUser(id int) error {
	tableName := "USER"
	return Delete(sm.DB, tableName, id)
}

func (sm *ORCLDBStorageManager) DeleteNote(id int) error {
	tableName := "NOTE"
	return Delete(sm.DB, tableName, id)
}

func (sm *ORCLDBStorageManager) GetUser(id int) (bool, error) {
	tableName := `"USER"`
	data, err := Get(sm.DB, tableName, id)
	if err != nil {
		fmt.Println("Error fetching user:", err)
		return false, err
	}

	if username, ok := data["USERNAME"]; ok && username != "" {
		fmt.Println("Found user:", data)
		return true, nil
	}
	return false, nil
}

func (sm *ORCLDBStorageManager) GetNote(id int) (bool, error) {
	tableName := "NOTE"
	data, err := Get(sm.DB, tableName, id)
	if err != nil {
		fmt.Println("Error fetching note:", err)
		return false, err
	}

	if title, ok := data["TITLE"]; ok && title != "" {
		fmt.Println("Found note:", data)
		return true, nil
	}
	return false, nil
}

func (sm *ORCLDBStorageManager) GetNotes() ([]types.Note, error) {
	return nil, nil
}

func (sm *ORCLDBStorageManager) Close() error {
	return sm.DB.Close()
}
