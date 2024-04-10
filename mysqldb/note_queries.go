package mysqldb

import (
	"database-tester/types"
	"database/sql"
	"fmt"
)

func (m *MySQLManager) GetNote(id string) (note types.Note, err error) {
	note = types.Note{}
	err = m.db.QueryRow("SELECT id, title, content, date_of_creation, date_of_modification, is_shared, id_user FROM NOTE WHERE id = ?", id).Scan(
		&note.ID,
		&note.Title,
		&note.Content,
		&note.DateOfCreation,
		&note.DateOfModification,
		&note.IsShared,
		&note.UserID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return note, nil
		}
		return note, fmt.Errorf("error getting note: %v", err)
	}
	return note, nil
}

func (m *MySQLManager) GetNotes(userID string) ([]types.Note, error) {
	var notes []types.Note
	query := "SELECT id, title, content, date_of_creation, date_of_modification, is_shared, id_user FROM NOTE WHERE id_user = ?"

	rows, err := m.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("error executing get notes by user ID query: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var note types.Note
		err := rows.Scan(&note.ID, &note.Title, &note.Content, &note.DateOfCreation, &note.DateOfModification, &note.IsShared, &note.UserID)
		if err != nil {
			return nil, fmt.Errorf("error scanning note: %v", err)
		}
		notes = append(notes, note)
	}

	return notes, nil
}

func (m *MySQLManager) InsertNote(note types.Note) (id string, err error) {
	stmt, err := m.db.Prepare("INSERT INTO NOTE (id, title, content, date_of_creation, date_of_modification, is_shared, id_user) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return "", fmt.Errorf("error preparing insert statement: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(note.ID, note.Title, note.Content, note.DateOfCreation, note.DateOfModification, note.IsShared, note.UserID)
	if err != nil {
		return "", fmt.Errorf("error executing insert statement: %v", err)
	}

	return note.ID, nil
}

func (m *MySQLManager) PatchNote(note types.Note) (id string, err error) {
	query := "UPDATE NOTE SET title = ?, content = ?, date_of_modification = ?, is_shared = ?, id_user = ? WHERE id = ?"
	_, err = m.db.Exec(query, note.Title, note.Content, note.DateOfModification, note.IsShared, note.UserID, note.ID)
	if err != nil {
		return "", fmt.Errorf("error updating note: %v", err)
	}

	return note.ID, nil
}

func (m *MySQLManager) DeleteNote(note types.Note) (id string, err error) {
	query := "DELETE FROM NOTE WHERE id = ?"
	_, err = m.db.Exec(query, note.ID)
	if err != nil {
		return "", fmt.Errorf("error deleting note: %v", err)
	}

	return note.ID, nil
}
