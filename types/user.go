package types

type User struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	IsAdmin   bool   `json:"is_admin"`
	IsEnabled bool   `json:"is_enabled"`
}
