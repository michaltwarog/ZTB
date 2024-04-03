package oracledb

import (
	"database/sql"
	"fmt"
	"log"

	"database-tester/types"

	"github.com/google/uuid"
)

func (om *OracleManager) InsertNote(note types.Note) (string, error) {
	// Ensure the note ID is included in the list of columns for insertion.
	query := `INSERT INTO "NOTE" (id, title, content, created_at, is_shared, id_user)
			  VALUES (:1, :2, :3, TO_DATE(:4, 'YYYY-MM-DD HH24:MI:SS'), :5, :6)`

	_, err := om.DB.Exec(query, note.ID, note.Title, note.Content, note.DateOfCreation, note.IsShared, note.UserID)
	if err != nil {
		log.Printf("Failed to insert note: %v", err)
		return "", err
	}

	return note.ID, nil
}

func (om *OracleManager) GetNote(noteID string) (types.Note, error) {
	var note types.Note

	query := "SELECT * FROM NOTE WHERE id = :1"

	row := om.DB.QueryRow(query, noteID)
	err := row.Scan(&note.ID, &note.Title, &note.Content, &note.DateOfCreation)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No note found with ID %s", noteID)
			return types.Note{}, err
		}
		log.Printf("Failed to retrieve note: %v", err)
		return types.Note{}, err
	}

	return note, nil
}

func (om *OracleManager) PatchNote(note types.Note) (string, error) {
	query := `UPDATE NOTE SET title = :1, content = :2, is_shared = :3 WHERE id = :4`

	_, err := om.DB.Exec(query, note.Title, note.Content, note.IsShared, note.ID)
	if err != nil {
		log.Printf("Failed to update note: %v", err)
		return "", err
	}

	return note.ID, nil
}

func (om *OracleManager) DeleteNote(note types.Note) (string, error) {
	query := `DELETE FROM NOTE WHERE id = :1`

	_, err := om.DB.Exec(query, note.ID)
	if err != nil {
		log.Printf("Failed to delete note: %v", err)
		return "", err
	}

	return note.ID, nil
}

func (om *OracleManager) GetNotes(userID string) ([]types.Note, error) {
	var notes []types.Note

	query := `SELECT id, title, content, created_at, is_shared FROM "NOTE" WHERE user_id = :1`

	rows, err := om.DB.Query(query, userID)
	if err != nil {
		log.Printf("Failed to retrieve notes: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var note types.Note
		if err := rows.Scan(&note.ID, &note.Title, &note.Content, &note.DateOfCreation, &note.IsShared); err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, err
		}
		notes = append(notes, note)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over rows: %v", err)
		return nil, err
	}

	return notes, nil
}

func TestNotes(oracleClient *OracleManager) {
	// Generate a new UUID for each note
	noteID := uuid.NewString()

	// Create a new note with a unique ID
	note := types.Note{
		ID:                 noteID,
		Title:              "Test Note",
		Content:            "This is a test note",
		DateOfCreation:     "2023-01-01 12:00:00", // Use the appropriate date format
		DateOfModification: "2023-01-01 12:00:00",
		IsShared:           true,
		UserID:             "1", // Assume this user ID exists
	}

	// Insert the note
	_, err := oracleClient.InsertNote(note)
	if err != nil {
		fmt.Println("Error inserting the note:", err)
		return
	}
	fmt.Println("Note inserted successfully with ID:", noteID)

	// Retrieve the note
	retrievedNote, err := oracleClient.GetNote(noteID)
	if err != nil {
		fmt.Println("Error retrieving the note:", err)
		return
	}
	fmt.Println("Retrieved note:", retrievedNote)

	// Patch the note
	retrievedNote.Title = "Updated Title"
	_, err = oracleClient.PatchNote(retrievedNote)
	if err != nil {
		fmt.Println("Error patching the note:", err)
		return
	}
	fmt.Println("Note patched successfully")

	// Assuming the patch does not change the ID, retrieve the patched note
	patchedNote, err := oracleClient.GetNote(noteID)
	if err != nil {
		fmt.Println("Error retrieving the patched note:", err)
		return
	}
	fmt.Println("Retrieved patched note:", patchedNote)

	// Delete the note
	_, err = oracleClient.DeleteNote(patchedNote)
	if err != nil {
		fmt.Println("Error deleting the note:", err)
		return
	}
	fmt.Println("Note deleted successfully")
}
