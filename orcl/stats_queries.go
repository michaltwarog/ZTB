package oracledb

import (
	"database-tester/types"
	"database/sql"
	"log"
)

func (om *OracleManager) GetUserStats(userID string) (types.UserStats, error) {
	var stats types.UserStats

	var firstName, lastName, email string
	var totalNotes, sharedNotes sql.NullInt64
	var latestNote sql.NullString

	query := `SELECT
		u.first_name,
		u.last_name,
		u.email,
		COUNT(n.id) AS total_notes,
		SUM(CASE WHEN n.is_shared = 1 THEN 1 ELSE 0 END) AS shared_notes,
		TO_CHAR(MAX(n.date_of_creation), 'YYYY-MM-DD HH24:MI:SS') AS latest_note
		FROM "USER" u
		LEFT JOIN note n ON u.id = n.id_user
		WHERE u.id = :1
		GROUP BY
			u.first_name,
			u.last_name,
			u.email
	`

	row := om.DB.QueryRow(query, userID)
	err := row.Scan(&firstName, &lastName, &email, &totalNotes, &sharedNotes, &latestNote)
	if err != nil {
		log.Printf("Failed to retrieve user stats: %v", err)
		return types.UserStats{}, err
	}

	stats.NotesCount = int(totalNotes.Int64)
	stats.SharedCount = int(sharedNotes.Int64)
	stats.LatestNoteDate = latestNote.String

	return stats, nil
}

func (om *OracleManager) GetUserModifiedNotesStats(userID string) (types.ModifiedNotesStats, error) {
	var stats types.ModifiedNotesStats

	var firstName, lastName, email string
	var avgModificationTime, maxUnmodifiedTime string
	var latestModification sql.NullTime

	query := `SELECT
		u.first_name,
		u.last_name,
		u.email,
		TO_CHAR(AVG(n.date_of_modification - n.date_of_creation), 'YYYY-MM-DD HH24:MI:SS') AS avg_modification_time,
		TO_CHAR(MAX(CASE WHEN n.date_of_modification IS NOT NULL THEN n.date_of_modification - n.date_of_creation ELSE NULL END), 'DD HH24:MI:SS') AS max_unmodified_time,
		TO_CHAR(MAX(n.date_of_modification), 'YYYY-MM-DD HH24:MI:SS') AS latest_modification
		FROM "USER" u
		JOIN note n ON u.id = n.id_user
		WHERE u.id = :1
		GROUP BY
			u.first_name,
			u.last_name,
			u.email
	`

	row := om.DB.QueryRow(query, userID)
	err := row.Scan(&firstName, &lastName, &email, &avgModificationTime, &maxUnmodifiedTime, &latestModification)
	if err != nil {
		log.Printf("Failed to retrieve modified notes stats: %v", err)
		return types.ModifiedNotesStats{}, err
	}

	stats.AvgModificationTime = avgModificationTime
	stats.MaxUnmodifiedTime = maxUnmodifiedTime
	if latestModification.Valid {
		stats.LatestModificationTime = latestModification.Time.Format("2006-01-02 15:04:05")
	} else {
		stats.LatestModificationTime = ""
	}

	return stats, nil
}
