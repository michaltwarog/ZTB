package types

type User struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	Username  string
	IsAdmin   bool
	IsEnabled bool
}
