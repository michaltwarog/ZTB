package mysqldb

import (
	"database-tester/types"
	"database/sql"
	"fmt"
)

func (m *MySQLManager) GetUser(id string) (user types.User, err error) {
	user = types.User{}
	err = m.db.QueryRow("SELECT id, first_name, last_name, email, username, is_admin, is_enabled FROM USER WHERE id = ?", id).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Username,
		&user.IsAdmin,
		&user.IsEnabled,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, nil // User not found but no error
		}
		return user, fmt.Errorf("error getting user: %v", err)
	}
	return user, nil
}

func (m *MySQLManager) InsertUser(user types.User) (id string, err error) {
	stmt, err := m.db.Prepare("INSERT INTO USER (id, first_name, last_name, email, username, is_admin, is_enabled) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return "", fmt.Errorf("error preparing insert statement: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.ID, user.FirstName, user.LastName, user.Email, user.Username, user.IsAdmin, user.IsEnabled)
	if err != nil {
		return "", fmt.Errorf("error executing insert statement: %v", err)
	}

	return user.ID, nil
}

func (m *MySQLManager) PatchUser(user types.User) (id string, err error) {
	query := "UPDATE USER SET first_name = ?, last_name = ?, email = ?, username = ?, is_admin = ?, is_enabled = ? WHERE id = ?"
	_, err = m.db.Exec(query, user.FirstName, user.LastName, user.Email, user.Username, user.IsAdmin, user.IsEnabled, user.ID)
	if err != nil {
		return "", fmt.Errorf("error updating user: %v", err)
	}

	return user.ID, nil
}

func (m *MySQLManager) DeleteUser(user types.User) (id string, err error) {
	query := "DELETE FROM USER WHERE id = ?"
	_, err = m.db.Exec(query, user.ID)
	if err != nil {
		return "", fmt.Errorf("error deleting user: %v", err)
	}

	return user.ID, nil
}
