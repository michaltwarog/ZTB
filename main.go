package main

import (
	"database-tester/types"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var oracleDB = map[string]string{
	"service":  "xe",
	"username": "system",
	"server":   "localhost",
	"port":     "1521",
	"password": "abc123",
}

type StorageManager interface {
	InsertNote(note types.Note) (id int, err error)
	InsertUser(user types.User) (id int, err error)
	PatchtNote(note types.Note) (id int, err error)
	PatchUser(user types.User) (id int, err error)
	GetNote(id int) (found bool, err error)
	GetUser(id int) (found bool, err error)
	GetNotes(userID int, pageToken types.Token) (notes []types.Note, nextPageToken types.Token, err error)
	DeleteNote(id int) (err error)
	DeleteUser(id int) (err error)
}

func main() {
	db := GetORCLDBStorageManager(oracleDB)
	defer func() {
		err := db.Close()
		if err != nil {
			fmt.Println("Can't close connection: ", err)
		}
	}()

	// id, err := db.GetUser(0)

	// if err != nil {
	// 	fmt.Println("Error fetching user:", err)
	// } else {
	// 	fmt.Println("Fetched user with ID:", id)
	// }

	// user := types.User{
	// 	First_Name: "Johdna",
	// 	Last_Name:  "Daofea",
	// 	Email:      "johaadn.doe@exampled.com",
	// 	Username:   "jodshndoed1",
	// 	Is_Admin:   false,
	// 	Is_Enabled: true,
	// }

	// id, err := db.InsertUser(user)
	// if err != nil {
	// 	fmt.Println("Error inserting user:", err)
	// } else {
	// 	fmt.Println("Inserted user with ID:", id)
	// }

	note := types.Note{
		Title:            "My first note",
		Content:          "This is the content of my first note",
		Date_Of_Creation: "2021-09-01 12:00:00",
		Is_Shared:        false,
		ID_User:          0,
	}

	id, err := db.InsertNote(note)
	if err != nil {
		fmt.Println("Error inserting note:", err)
	} else {
		fmt.Println("Inserted note with ID:", id)
	}

}
