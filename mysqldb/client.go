package mysqldb

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var mysqlDB = map[string]string{
	"username": "root",
	"password": "admin",
	"database": "mysql_ztb",
	"host":     "localhost",
	"port":     "3306",
}

type MySQLManager struct {
	db *sql.DB
}

func NewMySQLManager() (*MySQLManager, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8",
		mysqlDB["username"],
		mysqlDB["password"],
		mysqlDB["host"],
		mysqlDB["port"],
		mysqlDB["database"],
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(fmt.Errorf("error in sql.Open: %w", err))
	}

	err = db.Ping()
	if err != nil {
		panic(fmt.Errorf("error pinging db: %w", err))
	}

	fmt.Println("Connected to MySQL DB")
	return &MySQLManager{db: db}, nil
}
