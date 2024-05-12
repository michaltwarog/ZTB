package dynamodb

import (
	"database-tester/types"
	"time"
)

func (db DynamoManager) GetUserStats(userID string) (stats types.UserStats, err error) {

	notes, err := db.GetNotes(userID)
	if err != nil {
		return
	}

	stats.NotesCount = len(notes)
	for _, note := range notes {
		if note.IsShared {
			stats.SharedCount++
		}
		if note.DateOfCreation > stats.LatestNoteDate {
			stats.LatestNoteDate = note.DateOfCreation
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

	latestModificationTime := ""
	maxUnmodifiedTime := ""
	var avgModificationTime time.Duration

	for _, note := range notes {
		if note.DateOfModification != note.DateOfCreation && note.DateOfModification != "" {
			if latestModificationTime == "" || note.DateOfModification > latestModificationTime {
				latestModificationTime = note.DateOfModification
			}
			if maxUnmodifiedTime == "" || note.DateOfModification < maxUnmodifiedTime {
				maxUnmodifiedTime = note.DateOfModification
			}

			layout := "2006-01-02"
			doc, _ := time.Parse(layout, note.DateOfModification)
			dom, _ := time.Parse(layout, note.DateOfCreation)
			avgModificationTime = avgModificationTime + (doc.Sub(dom))
		}
	}

	return
}
