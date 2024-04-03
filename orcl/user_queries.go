package oracledb

import (
	"database/sql"
	"fmt"
	"log"

	"database-tester/types"

	"github.com/google/uuid"
)

func (om *OracleManager) InsertUser(user types.User) (string, error) {
	query := `INSERT INTO "USER" (id, first_name, last_name, email, username, is_admin, is_enabled)
			  VALUES (:1, :2, :3, :4, :5, :6, :7)`

	userID := uuid.NewString()

	_, err := om.DB.Exec(query, userID, user.FirstName, user.LastName, user.Email, user.Username, user.IsAdmin, user.IsEnabled)
	if err != nil {
		log.Printf("Failed to insert user: %v", err)
		return "", err
	}

	return userID, nil
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

func TestUsers(oracleClient *OracleManager) {
	// Generate a new UUID for the user
	userID := uuid.NewString()

	// Create a new user with a unique ID
	user := types.User{
		ID:        userID,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Username:  "johndoe",
		IsAdmin:   false,
		IsEnabled: true,
	}

	// Insert the user
	_, err := oracleClient.InsertUser(user)
	if err != nil {
		fmt.Println("Error inserting user:", err)
		return
	}
	fmt.Println("User inserted successfully with ID:", userID)

	// Retrieve the user
	retrievedUser, err := oracleClient.GetUser(userID)
	if err != nil {
		fmt.Println("Error retrieving user:", err)
		return
	}
	fmt.Println("Retrieved User:", retrievedUser)

	// Patch the user
	retrievedUser.FirstName = "Jane" // Modify some attributes
	retrievedUser.LastName = "Doe"
	_, err = oracleClient.PatchUser(retrievedUser)
	if err != nil {
		fmt.Println("Error patching user:", err)
		return
	}
	fmt.Println("User patched successfully")

	// Delete the user
	_, err = oracleClient.DeleteUser(retrievedUser)
	if err != nil {
		fmt.Println("Error deleting user:", err)
		return
	}
	fmt.Println("User deleted successfully")
}
