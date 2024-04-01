package types

type Note struct {
	ID                 string
	Title              string
	Content            string
	DateOfCreation     string
	DateOfModification string
	IsShared           bool
	UserID             string
}
