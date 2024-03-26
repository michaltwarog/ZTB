package types

type Note struct {
	ID                 string
	Title              string
	Content            string
	DateOfCreation     string
	DateofModification string
	IsShared           bool
	UserID             int
}
