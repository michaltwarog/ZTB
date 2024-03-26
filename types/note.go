package types

type Note struct {
	ID                   int
	Title                string
	Content              string
	Date_Of_Creation     string
	Date_Of_Modification string
	Is_Shared            bool
	ID_User              int
}
