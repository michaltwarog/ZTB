package main

import (
	"database-tester/dynamodb"
	"database-tester/types"

	"fmt"
)

type StorageManager interface {
	InsertNote(note types.Note) (id string, err error)
	GetNote(id string) (note types.Note, err error)
	GetNotes(userID string) (notes []types.Note, err error)
	PatchNote(note types.Note) (id string, err error)
	DeleteNote(note types.Note) (id string, err error)

	InsertUser(user types.User) (id string, err error)
	GetUser(id string) (user types.User, err error)
	PatchUser(user types.User) (id string, err error)
	DeleteUser(user types.User) (id string, err error)
}

func main() {

	dynamoClient, err := dynamodb.NewDynamoManager("local")

	if err != nil {
		fmt.Println("Error creating dynamodb client")
		return
	}

	fmt.Println("Test Users")
	dynamodb.TestUsers(dynamoClient)

	fmt.Println("Test Notes")
	dynamodb.TestNotes(dynamoClient)

}
