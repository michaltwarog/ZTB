package mysqldb

import (
	"database-tester/types"
	"database/sql"
	"fmt"
)

func (m *MySQLManager) GetUserStats(userID string) (stats types.UserStats, err error) {
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
    DATE_FORMAT(MAX(n.date_of_creation), '%Y-%m-%d %H:%i:%s') AS latest_note,
    COALESCE(DATE_FORMAT(MAX(sub.latest_shared_note), '%Y-%m-%d %H:%i:%s'), 'No data') AS latest_shared_note_creation_date,
    AVG(sub.shared_note_creation_delay) AS avg_delay_for_shared_note_creation
FROM USER u
LEFT JOIN note n ON u.id = n.id_user
LEFT JOIN (
    SELECT
        n1.id_user,
        n1.date_of_creation AS latest_shared_note,
        TIMESTAMPDIFF(MINUTE, MIN(n2.date_of_creation), n1.date_of_creation) AS shared_note_creation_delay
    FROM note n1
    INNER JOIN note n2 ON n1.id_user = n2.id_user AND n1.is_shared = 1 AND n2.is_shared = 1
    WHERE n1.date_of_creation > n2.date_of_creation
    GROUP BY n1.id_user, n1.date_of_creation
) sub ON u.id = sub.id_user
WHERE u.id = ?
GROUP BY
    u.first_name,
    u.last_name,
    u.email;`

	err = m.db.QueryRow(query, userID).Scan(&firstName, &lastName, &email, &totalNotes, &sharedNotes, &latestNote, &latestSharedNoteCreationDate, &avgDelayForSharedNoteCreation)
	if err != nil {
		if err == sql.ErrNoRows {
			return stats, nil
		}
		return stats, fmt.Errorf("error getting user stats: %v", err)
	}

	stats.NotesCount = int(totalNotes.Int64)
	stats.SharedCount = int(sharedNotes.Int64)
	stats.LatestNoteDate = latestNote.String
	return stats, nil
}

func (m *MySQLManager) GetUserModifiedNotesStats(userID string) (stats types.ModifiedNotesStats, err error) {
	var firstName, lastName, email string
	var avgModificationTime, maxUnmodifiedTime, latestModification string

	query := `SELECT
    u.first_name,
    u.last_name,
    u.email,
    COALESCE(FORMAT(
        AVG(sub.avg_mod_seconds), 2
    ), 'No data') AS avg_modification_time,
    COALESCE(FORMAT(
        MAX(sub.max_unmod_seconds), 2
    ), 'No data') AS max_unmodified_time,
    COALESCE(DATE_FORMAT(MAX(sub.latest_mod), '%Y-%m-%d %H:%i:%s'), 'No data') AS latest_modification
FROM USER u
LEFT JOIN (
    SELECT
        n.id_user,
        TIMESTAMPDIFF(SECOND, n.date_of_creation, n.date_of_modification) AS avg_mod_seconds,
        CASE
            WHEN n.date_of_modification IS NOT NULL THEN
                TIMESTAMPDIFF(SECOND, n.date_of_creation, n.date_of_modification)
            ELSE NULL
        END AS max_unmod_seconds,
        n.date_of_modification AS latest_mod
    FROM note n
    GROUP BY n.id_user, n.date_of_modification, n.date_of_creation
) sub ON u.id = sub.id_user
WHERE u.id = ?
GROUP BY
    u.first_name,
    u.last_name,
    u.email;
`

	err = m.db.QueryRow(query, userID).Scan(&firstName, &lastName, &email, &avgModificationTime, &maxUnmodifiedTime, &latestModification)

	if err != nil {
		if err == sql.ErrNoRows {
			return stats, nil
		}
		return stats, fmt.Errorf("error getting modified notes stats: %v", err)
	}

	stats.AvgModificationTime = avgModificationTime
	stats.MaxUnmodifiedTime = maxUnmodifiedTime
	stats.LatestModificationTime = latestModification

	return stats, nil
}
