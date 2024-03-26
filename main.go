package main

import (
	"database-tester/dynamodb"
	"database-tester/types"

	"fmt"
)

type StorageManager interface {
	InsertNote(note types.Note) (id int, err error)
	InsertUser(user types.User) (id int, err error)
	PatchNote(note types.Note) (id int, err error)
	PatchUser(user types.User) (id int, err error)
	GetNote(id string) (note types.Note, found bool, err error)
	GetUser(id int) (found bool, err error)
	GetNotes(userID int, pageToken types.Token) (notes []types.Note, nextPageToken types.Token, err error)
}

func main() {

	dynamoClient, err := dynamodb.NewDynamoManager("local")

	if err != nil {
		fmt.Println("Error creating dynamodb client")
		return
	}

	// err = dynamoClient.CreateDynamoDBTable(dynamodb.NotesTableName, dynamodb.TableInput)

	// if err != nil {
	// 	fmt.Println("Error creating dynamodb table")
	// 	fmt.Println(err)
	// 	return
	// }

	// Create a new note
	note := types.Note{
		ID:                 "2",
		Title:              "asd",
		Content:            "asd",
		DateOfCreation:     "123",
		DateofModification: "123",
		IsShared:           true,
		UserID:             1,
	}

	// Insert the note
	id, err := dynamoClient.InsertNote(note)
	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Println(id)

	note, found, err := dynamoClient.GetNote(id)
	if err != nil {
		fmt.Print(err)
		return
	}

	if found {
		fmt.Println(note)
	} else {
		fmt.Println("Note not found")
	}

	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Println("Hello, World!")
	fmt.Println("Hello, World!")
	fmt.Println("Hello, World!")
	fmt.Println("Hello, World!")
	fmt.Println("Hello, World!")

}
