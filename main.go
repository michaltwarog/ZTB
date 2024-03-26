package main

import "database-tester/types"

type StorageManager interface {
	InsertNote(note types.Note) (id int, err error)
	InsertUser(user types.User) (id int, err error)
	PatchtNote(note types.Note) (id int, err error)
	PatchUser(user types.User) (id int, err error)
	GetNote(id int) (found bool, err error)
	GetUser(id int) (found bool, err error)
	GetNotes(userID int, pageToken types.Token) (notes []types.Note, nextPageToken types.Token, err error)
}

func main() {

}
