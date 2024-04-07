package oracledb

import (
	"database/sql"
	"log"

	"database-tester/types"
)

func (om *OracleManager) InsertNote(note types.Note) (string, error) {
	query := `INSERT INTO NOTE (id, title, content, date_of_creation, is_shared, id_user)
			  VALUES (:1, :2, :3, TO_TIMESTAMP_TZ(:4, 'YYYY-MM-DD"T"HH24:MI:SS.FF TZH:TZM'), :5, :6)`

	_, err := om.DB.Exec(query, note.ID, note.Title, note.Content, note.DateOfCreation, note.IsShared, note.UserID)
	if err != nil {
		log.Printf("Failed to insert note: %v", err)
		return "", err
	}

	return note.ID, nil
}

func (om *OracleManager) GetNote(noteID string) (types.Note, error) {
	var note types.Note
	var dateOfModification sql.NullString

	query := "SELECT id, title, content, date_of_creation, date_of_modification, is_shared, id_user FROM NOTE WHERE id = :1"

	row := om.DB.QueryRow(query, noteID)
	err := row.Scan(&note.ID, &note.Title, &note.Content, &note.DateOfCreation, &dateOfModification, &note.IsShared, &note.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No note found with ID %s", noteID)
			return types.Note{}, err
		}
		log.Printf("Failed to retrieve note: %v", err)
		return types.Note{}, err
	}

	if dateOfModification.Valid {
		note.DateOfModification = dateOfModification.String
	} else {
		note.DateOfModification = ""
	}

	return note, nil
}

func (om *OracleManager) PatchNote(note types.Note) (string, error) {
	query := `UPDATE NOTE SET title = :1, content = :2,date_of_modification = TO_TIMESTAMP_TZ(:3, 'YYYY-MM-DD"T"HH24:MI:SS.FF TZH:TZM'), is_shared = :4 WHERE id = :5`

	_, err := om.DB.Exec(query, note.Title, note.Content, note.DateOfModification, note.IsShared, note.ID)
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

	query := `SELECT id, title, content, date_of_creation, date_of_modification, is_shared FROM "NOTE" WHERE id_user = :1`

	rows, err := om.DB.Query(query, userID)
	if err != nil {
		log.Printf("Failed to retrieve notes: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var note types.Note
		var dateOfModification sql.NullString

		if err := rows.Scan(&note.ID, &note.Title, &note.Content, &note.DateOfCreation, &dateOfModification, &note.IsShared); err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, err
		}

		if dateOfModification.Valid {
			note.DateOfModification = dateOfModification.String
		} else {
			note.DateOfModification = ""
		}

		notes = append(notes, note)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over rows: %v", err)
		return nil, err
	}

	return notes, nil
}
