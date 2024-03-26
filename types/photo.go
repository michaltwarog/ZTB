package types

type Photo struct {
	ID           int
	Name         string
	Path         string
	DateOfUpload string
	IsShared     bool
	UserID       int
}
