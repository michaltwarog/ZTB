package oracledb

import (
	"database/sql"
	"fmt"

	_ "github.com/sijms/go-ora/v2"
)

var oracleDB = map[string]string{
	"service":  "xe",
	"username": "system",
	"server":   "localhost",
	"port":     "1521",
	"password": "abc123",
}

type OracleManager struct {
	DB *sql.DB
}

// NewOracleManager establishes a connection to the Oracle database and returns an OracleManager.
func NewOracleManager() (*OracleManager, error) {
	db, err := sql.Open("oracle", fmt.Sprintf("oracle://%s:%s@%s:%s/%s",
		oracleDB["username"],
		oracleDB["password"],
		oracleDB["server"],
		oracleDB["port"],
		oracleDB["service"],
	))
	if err != nil {
		panic(fmt.Errorf("error in sql.Open: %w", err))
	}

	err = db.Ping()
	if err != nil {
		panic(fmt.Errorf("error pinging db: %w", err))
	}
	print("Connected to Oracle DB\n")
	return &OracleManager{DB: db}, nil
}
