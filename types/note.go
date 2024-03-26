package types

type Note struct {
	ID                 int
	Title              string
	Content            string
	DateOfCreation     string
	DateofModification string
	IsShared           bool
	UserID             int
}
