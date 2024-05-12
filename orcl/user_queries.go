package oracledb

import (
	"database/sql"
	"log"

	"database-tester/types"
)

func (om *OracleManager) InsertUser(user types.User) (string, error) {
	query := `INSERT INTO "USER" (id, first_name, last_name, email, username, is_admin, is_enabled)
			  VALUES (:1, :2, :3, :4, :5, :6, :7)`

	_, err := om.DB.Exec(query, user.ID, user.FirstName, user.LastName, user.Email, user.Username, user.IsAdmin, user.IsEnabled)
	if err != nil {
		log.Printf("Failed to insert user: %v", err)
		return "", err
	}

	return user.ID, nil
}

func (om *OracleManager) GetUser(userID string) (types.User, error) {
	var user types.User

	query := `SELECT * FROM "USER" WHERE id = :1`

	row := om.DB.QueryRow(query, userID)

	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Username, &user.IsAdmin, &user.IsEnabled)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No user found with ID %s", userID)
			return types.User{}, err
		}
		log.Printf("Failed to retrieve user: %v", err)
		return types.User{}, err
	}
	// log.Printf("User found: %v", user)
	return user, nil
}

func (om *OracleManager) PatchUser(user types.User) (string, error) {
	query := `UPDATE "USER" SET first_name = :1, last_name = :2, email = :3, username = :4, is_admin = :5, is_enabled = :6 WHERE id = :7`

	_, err := om.DB.Exec(query, user.FirstName, user.LastName, user.Email, user.Username, user.IsAdmin, user.IsEnabled, user.ID)
	if err != nil {
		log.Printf("Failed to update user: %v", err)
		return "", err
	}

	return user.ID, nil
}

func (om *OracleManager) DeleteUser(user types.User) (string, error) {
	query := `DELETE FROM "USER" WHERE id = :1`

	_, err := om.DB.Exec(query, user.ID)
	if err != nil {
		log.Printf("Failed to delete user: %v", err)
		return "", err
	}

	return user.ID, nil
}
