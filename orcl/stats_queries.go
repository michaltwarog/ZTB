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
	var latestNote, latestSharedNoteCreationDate sql.NullString
	var avgDelayForSharedNoteCreation sql.NullFloat64

	query := `SELECT
        u.first_name,
        u.last_name,
        u.email,
        COUNT(n.id) AS total_notes,
        SUM(CASE WHEN n.is_shared = 1 THEN 1 ELSE 0 END) AS shared_notes,
        TO_CHAR(MAX(n.date_of_creation), 'YYYY-MM-DD HH24:MI:SS') AS latest_note,
        COALESCE(TO_CHAR(MAX(sub.latest_shared_note), 'YYYY-MM-DD HH24:MI:SS'), 'No data') AS latest_shared_note_creation_date,
        AVG(sub.shared_note_creation_delay) AS avg_delay_for_shared_note_creation
    FROM "USER" u
    LEFT JOIN note n ON u.id = n.id_user
    LEFT JOIN (
        SELECT
            n1.id_user,
            n1.date_of_creation AS latest_shared_note,
            EXTRACT(DAY FROM (n1.date_of_creation - MIN(n2.date_of_creation))) * 1440 +
            EXTRACT(HOUR FROM (n1.date_of_creation - MIN(n2.date_of_creation))) * 60 +
            EXTRACT(MINUTE FROM (n1.date_of_creation - MIN(n2.date_of_creation))) AS shared_note_creation_delay
        FROM note n1
        INNER JOIN note n2 ON n1.id_user = n2.id_user AND n1.is_shared = 1 AND n2.is_shared = 1
        WHERE n1.date_of_creation > n2.date_of_creation
        GROUP BY n1.id_user, n1.date_of_creation
    ) sub ON u.id = sub.id_user
    WHERE u.id = :1
    GROUP BY
        u.first_name,
        u.last_name,
        u.email
    `

	row := om.DB.QueryRow(query, userID)
	err := row.Scan(&firstName, &lastName, &email, &totalNotes, &sharedNotes, &latestNote, &latestSharedNoteCreationDate, &avgDelayForSharedNoteCreation)
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
	var avgModificationTime, maxUnmodifiedTime, latestModification string

	query := `SELECT
    u.first_name,
    u.last_name,
    u.email,
    COALESCE(TO_CHAR(
        AVG(sub.avg_mod_seconds), 'FM9999999990.00'
    ), 'No data') AS avg_modification_time,
    COALESCE(TO_CHAR(
        MAX(sub.max_unmod_seconds), 'FM9999999990.00'
    ), 'No data') AS max_unmodified_time,
    COALESCE(TO_CHAR(MAX(sub.latest_mod), 'YYYY-MM-DD HH24:MI:SS'), 'No data') AS latest_modification
FROM "USER" u
LEFT JOIN (
    SELECT
        n.id_user,
        EXTRACT(DAY FROM (n.date_of_modification - n.date_of_creation)) * 86400 +
        EXTRACT(HOUR FROM (n.date_of_modification - n.date_of_creation)) * 3600 +
        EXTRACT(MINUTE FROM (n.date_of_modification - n.date_of_creation)) * 60 +
        EXTRACT(SECOND FROM (n.date_of_modification - n.date_of_creation)) AS avg_mod_seconds,
        CASE
            WHEN n.date_of_modification IS NOT NULL THEN
                EXTRACT(DAY FROM (n.date_of_modification - n.date_of_creation)) * 86400 +
                EXTRACT(HOUR FROM (n.date_of_modification - n.date_of_creation)) * 3600 +
                EXTRACT(MINUTE FROM (n.date_of_modification - n.date_of_creation)) * 60 +
                EXTRACT(SECOND FROM (n.date_of_modification - n.date_of_creation))
            ELSE NULL
        END AS max_unmod_seconds,
        n.date_of_modification AS latest_mod
    FROM note n
    GROUP BY n.id_user, n.date_of_modification, n.date_of_creation
) sub ON u.id = sub.id_user
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
	stats.LatestModificationTime = latestModification
	// fmt.Println("stats", stats)
	return stats, nil
}
