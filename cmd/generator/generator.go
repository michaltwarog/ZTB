package main

import (
	"database-tester/types"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/google/uuid"
)

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func generateRandomUser() types.User {
	return types.User{
		ID:        uuid.New().String(),
		FirstName: generateRandomString(8),
		LastName:  generateRandomString(8),
		Email:     generateRandomString(10) + "@example.com",
		Username:  generateRandomString(8),
		IsAdmin:   rand.Float32() < 0.5,
		IsEnabled: true,
	}
}

func generateRandomNote(userID string) types.Note {
	return types.Note{
		ID:                 uuid.New().String(),
		Title:              generateRandomString(15),
		Content:            generateRandomString(50),
		DateOfCreation:     time.Now().Format(time.RFC3339),
		DateOfModification: time.Now().Format(time.RFC3339),
		IsShared:           rand.Float32() < 0.5,
		UserID:             userID,
	}
}

func generateHyperUserWithNotes() (types.User, []types.Note) {
	user := generateRandomUser()
	notes := make([]types.Note, 0)
	for i := 0; i < 100000; i++ {
		if i%10000 == 0 {
			fmt.Println("Generated", i, "notes")
		}
		notes = append(notes, generateRandomNote(user.ID))
	}
	return user, notes
}

func main() {

	// Open files for appending data
	usersFile, err := os.OpenFile("users.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening users.json:", err)
		os.Exit(1)
	}
	defer usersFile.Close()

	notesFile, err := os.OpenFile("notes.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening notes.json:", err)
		os.Exit(1)
	}
	defer notesFile.Close()

	updatedNotesFile, err := os.OpenFile("updatedNotes.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening updatedNotes.json:", err)
		os.Exit(1)
	}
	defer updatedNotesFile.Close()

	userIDsFile, err := os.OpenFile("userIDs.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening userIDs.json:", err)
		os.Exit(1)
	}
	defer userIDsFile.Close()

	noteIDsFile, err := os.OpenFile("noteIDs.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening noteIDs.json:", err)
		os.Exit(1)
	}
	defer noteIDsFile.Close()
	var users []types.User
	var notes []types.Note
	var updatedNotes []types.Note

	{
		hyperUser, notes := generateHyperUserWithNotes()
		users = append(users, hyperUser)
		notes = append(notes, notes...)
		if _, err := userIDsFile.WriteString(hyperUser.ID + "\n"); err != nil {
			fmt.Println("Error writing user ID to file:", err)
		}
		for _, note := range notes {
			if _, err := noteIDsFile.WriteString(note.ID + "\n"); err != nil {
				fmt.Println("Error writing note ID to file:", err)
			}
		}
	}
	start := time.Now()
	nUsers := 100
	nNotes := 100
	updatedNote := types.Note{}
	// Generate and append users and user IDs
	for i := 0; i < nUsers; i++ {

		if i%100 == 0 {
			fmt.Println("Percent done:", float32(i)/float32(nUsers)*100)
			fmt.Printf("Time elapsed: %.0f seconds\n", time.Since(start).Seconds())
			remaining := float64(nUsers-i) * time.Since(start).Seconds() / float64(i)
			fmt.Printf("Time remaining: %.0f seconds\n", remaining)
		}
		user := generateRandomUser()
		users = append(users, user)

		if _, err := userIDsFile.WriteString(user.ID + "\n"); err != nil {
			fmt.Println("Error writing user ID to file:", err)
		}
		// Generate and append notes and note IDs
		for i := 0; i < nNotes; i++ {
			note := generateRandomNote(user.ID)
			notes = append(notes, note)
			{
				updatedNote = note
				updatedNote.Title = "Updated " + updatedNote.Title
				updatedNotes = append(updatedNotes, updatedNote)
			}
			if _, err := noteIDsFile.WriteString(note.ID + "\n"); err != nil {
				fmt.Println("Error writing note ID to file:", err)
			}
		}

	}

	usersJSON, err := json.Marshal(users)
	if err != nil {
		fmt.Println("Error marshaling users:", err)
		return
	}
	if _, err := usersFile.Write(usersJSON); err != nil {
		fmt.Println("Error writing users to file:", err)
	}

	notesJSON, err := json.Marshal(notes)
	if err != nil {
		fmt.Println("Error marshaling notes:", err)
		return
	}
	if _, err := notesFile.Write(notesJSON); err != nil {
		fmt.Println("Error writing notes to file:", err)
	}

	updatedNotesJSON, err := json.Marshal(updatedNotes)
	if err != nil {
		fmt.Println("Error marshaling updatedNotes:", err)
		return
	}
	if _, err := updatedNotesFile.Write(updatedNotesJSON); err != nil {
		fmt.Println("Error writing notes to file:", err)
	}

	fmt.Println("Data appended to files")
}
