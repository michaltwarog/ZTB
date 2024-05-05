package dynamodb

import "database-tester/types"

func (db DynamoManager) GetUserStats(userID string) (stats types.UserStats, err error) {
	// Query the notes table for the user
	notes, err := db.GetNotes(userID)
	if err != nil {
		return
	}
	// Count the notes
	stats.TotalNotes = len(notes)
	// Count the notes that have been updated
	for _, note := range notes {
		if note.ModifiedAt != note.CreatedAt {
			stats.ModifiedNotes++
		}
	}
	// Count the notes that have been deleted
	for _, note := range notes {
		if note.DeletedAt != 0 {
			stats.DeletedNotes++
		}
	}
	return
}
func (db DynamoManager) GetUserModifiedNotesStats(userID string) (stats types.ModifiedNotesStats, err error) {
	// Query the notes table for the user
	notes, err := db.GetNotes(userID)
	if err != nil {
		return
	}
	// Count the notes that have been updated
	for _, note := range notes {
		if note.ModifiedAt != note.CreatedAt {
			stats.ModifiedNotes++
		}
	}
	return
}
