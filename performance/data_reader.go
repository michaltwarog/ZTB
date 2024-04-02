package performance

import (
	"database-tester/types"
	"encoding/json"
	"fmt"
	"os"
)

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

	{
		fmt.Println("Users:", users[0])
		fmt.Println("Number of users:", len(users))

		fmt.Println("Notes:", notes[0])
		fmt.Println("Number of notes:", len(notes))

		fmt.Println("Updated Notes:", updatedNotes[0])
		fmt.Println("Number of updated notes:", len(updatedNotes))

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
