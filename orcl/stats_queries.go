package oracledb

import (
	"database-tester/types"
	"database/sql"
	"fmt"
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
	fmt.Println("stats", stats)
	return stats, nil
}

func (om *OracleManager) GetUserModifiedNotesStats(userID string) (types.ModifiedNotesStats, error) {
	var stats types.ModifiedNotesStats

	var firstName, lastName, email string
	var avgModificationTime, maxUnmodifiedTime, latestModification string // Change latestModification to string

	query := `SELECT
        u.first_name,
        u.last_name,
        u.email,
        COALESCE(TO_CHAR(
            AVG(
                EXTRACT(DAY FROM (n.date_of_modification - n.date_of_creation)) * 86400 +
                EXTRACT(HOUR FROM (n.date_of_modification - n.date_of_creation)) * 3600 +
                EXTRACT(MINUTE FROM (n.date_of_modification - n.date_of_creation)) * 60 +
                EXTRACT(SECOND FROM (n.date_of_modification - n.date_of_creation))
            ), 'FM9999999990.00'), 'No data') AS avg_modification_time,
        COALESCE(TO_CHAR(
            MAX(
                CASE
                    WHEN n.date_of_modification IS NOT NULL THEN
                        EXTRACT(DAY FROM (n.date_of_modification - n.date_of_creation)) * 86400 +
                        EXTRACT(HOUR FROM (n.date_of_modification - n.date_of_creation)) * 3600 +
                        EXTRACT(MINUTE FROM (n.date_of_modification - n.date_of_creation)) * 60 +
                        EXTRACT(SECOND FROM (n.date_of_modification - n.date_of_creation))
                    ELSE NULL
                END
            ), 'FM9999999990.00'), 'No data') AS max_unmodified_time,
        COALESCE(TO_CHAR(MAX(n.date_of_modification), 'YYYY-MM-DD HH24:MI:SS'), 'No data') AS latest_modification
    FROM "USER" u
    LEFT JOIN note n ON u.id = n.id_user
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
	stats.LatestModificationTime = latestModification // Directly use the string
	fmt.Println("stats", stats)
	return stats, nil
}
