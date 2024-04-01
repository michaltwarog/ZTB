package main

import (
	"database-tester/types"
	"encoding/json"
	"os"

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

	users, notes, updatedNotes, err := readDataFromFiles()

	{
		fmt.Println("Users:", users[0])
		fmt.Println("Number of users:", len(users))

		fmt.Println("Notes:", notes[0])
		fmt.Println("Number of notes:", len(notes))

		fmt.Println("Updated Notes:", updatedNotes[0])
		fmt.Println("Number of updated notes:", len(updatedNotes))
		if err != nil {
			fmt.Println("Error reading files:", err)
			return
		}
	}
}

func readDataFromFiles() ([]types.User, []types.Note, []types.Note, error) {
	users, err := readUsersFromFile("cmd/generator/users.json")
	if err != nil {
		return nil, nil, nil, err
	}

	notes, err := readNotesFromFile("cmd/generator/notes.json")
	if err != nil {
		return nil, nil, nil, err
	}

	updatedNotes, err := readNotesFromFile("cmd/generator/updatedNotes.json")
	if err != nil {
		return nil, nil, nil, err
	}

	return users, notes, updatedNotes, nil

}

func readUsersFromFile(filename string) ([]types.User, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var users []types.User
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&users); err != nil {
		return nil, err
	}

	return users, nil
}

func readNotesFromFile(filename string) ([]types.Note, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var notes []types.Note
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&notes); err != nil {
		return nil, err
	}

	return notes, nil
}
