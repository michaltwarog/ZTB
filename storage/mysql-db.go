package storage

import (
	"database-tester/types"
	"database/sql"
	"fmt"
)

type MySQLStorageManager struct {
	db *sql.DB
}

func NewMySQLStorageManager(db *sql.DB) *MySQLStorageManager {
	return &MySQLStorageManager{db: db}
}

func (m *MySQLStorageManager) InsertNote(note types.Note) (id int, err error) {
	stmt, err := m.db.Prepare("INSERT INTO NOTE(title, content, date_of_creation, date_of_modification, is_shared, id_user) VALUES(?, ?, ?, ?, ?, ?)")
	if err != nil {
		return 0, fmt.Errorf("error preparing insert statement: %v", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(note.Title, note.Content, note.Date_Of_Creation, note.Date_Of_Modification, note.Is_Shared, note.ID_User)
	if err != nil {
		return 0, fmt.Errorf("error executing insert statement: %v", err)
	}

	id64, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error getting last insert ID: %v", err)
	}
	id = int(id64)
	return int(id), nil
}

func (m *MySQLStorageManager) InsertUser(user types.User) (id int, err error) {
	stmt, err := m.db.Prepare("INSERT INTO USER(first_name, last_name, email, username, is_admin, is_enabled) VALUES(?, ?, ?, ?, ?, ?)")
	if err != nil {
		return 0, fmt.Errorf("error preparing insert statement: %v", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(user.First_Name, user.Last_Name, user.Email, user.Username, user.Is_Admin, user.Is_Enabled)
	if err != nil {
		return 0, fmt.Errorf("error executing insert statement: %v", err)
	}

	id64, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error getting last insert ID: %v", err)
	}
	id = int(id64)
	return int(id), nil
}
func (m *MySQLStorageManager) GetUser(id int) (user types.User, found bool, err error) {
	user = types.User{}
	err = m.db.QueryRow("SELECT id, first_name, last_name, email, username, is_admin, is_enabled FROM USER WHERE id = ?", id).Scan(
		&user.ID,
		&user.First_Name,
		&user.Last_Name,
		&user.Email,
		&user.Username,
		&user.Is_Admin,
		&user.Is_Enabled,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, false, nil // Nie znaleziono użytkownika, ale brak błędu
		}
		return user, false, fmt.Errorf("error getting user: %v", err)
	}
	return user, true, nil
}

func (m *MySQLStorageManager) GetNote(id int) (note types.Note, found bool, err error) {
	note = types.Note{}
	err = m.db.QueryRow("SELECT id, title, content, date_of_creation, date_of_modification, is_shared, id_user FROM note WHERE id = ?", id).Scan(
		&note.ID_Note,
		&note.Title,
		&note.Content,
		&note.Date_Of_Creation,
		&note.Date_Of_Modification,
		&note.Is_Shared,
		&note.ID_User,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return note, false, nil // Nie znaleziono notatki, ale brak błędu
		}
		return note, false, fmt.Errorf("error getting note: %v", err)
	}
	return note, true, nil
}
