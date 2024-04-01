package types

type Note struct {
	ID                 string `json:"id"`
	Title              string `json:"title"`
	Content            string `json:"content"`
	DateOfCreation     string `json:"date_of_creation"`
	DateOfModification string `json:"date_of_modification"`
	IsShared           bool   `json:"is_shared"`
	UserID             string `json:"user_id"`
}
