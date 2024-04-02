package main

import (
	"database-tester/storage"
	"database-tester/types"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type StorageManager interface {
	InsertNote(note types.Note) (id int, err error)
	InsertUser(user types.User) (id int, err error)
	PatchNote(note types.Note) (id int, err error)
	PatchUser(user types.User) (id int, err error)
	GetNote(id int) (found bool, err error)
	GetUser(id int) (found bool, err error)
}

func main() {
	// Ustaw dane połączeniowe
	dsn := "root:admin@tcp(localhost:3306)/mysql_ztb?charset=utf8"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Test połączenia
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	storageManager := storage.NewMySQLStorageManager(db)

	//user := types.User{
	//	ID:         1,
	//	First_Name: "Michal",
	//	Last_Name:  "Twarog",
	//	Email:      "mtwarog@gmail.com",
	//	Username:   "michalT",
	//	Is_Admin:   false,
	//	Is_Enabled: true,
	//}
	//if id, err := storageManager.InsertUser(user); err != nil {
	//	fmt.Printf("Error inserting user: %s", err)
	//} else {
	//	fmt.Printf("Note inserted with ID: %d", id)
	//}

	//note := types.Note{
	//	ID_Note:              1,
	//	Title:                "Notatka1",
	//	Content:              "contentcontentcontent",
	//	Date_Of_Creation:     "2024-11-25 11:55:12",
	//	Date_Of_Modification: "2024-11-25 11:55:12",
	//	Is_Shared:            true,
	//	ID_User:              1,
	//}
	//
	//if id, err := storageManager.InsertNote(note); err != nil {
	//	fmt.Printf("Error inserting user: %s", err)
	//} else {
	//	fmt.Printf("Note inserted with ID: %d", id)
	//}

	userID := 1
	user, found, err := storageManager.GetUser(userID)
	if err != nil {
		fmt.Printf("Error getting user: %s\n", err)
		return
	}
	if found {
		fmt.Printf("User found: %+v\n", user)
	} else {
		fmt.Println("User not found")
	}
}
