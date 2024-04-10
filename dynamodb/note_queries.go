package dynamodb

import (
	"context"
	localtypes "database-tester/types"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func (db DynamoManager) InsertNote(note localtypes.Note) (id string, err error) {
	marshalledNote, err := attributevalue.MarshalMap(note)
	if err != nil {
		fmt.Println("Error marshalling note")
		fmt.Print(err)
		return
	}
	_, err = db.Client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(NotesTableName), Item: marshalledNote,
	})

	if err != nil {
		fmt.Println("Error inserting note")
		fmt.Print(err)
		return
	}

	return note.ID, err
}
func (db DynamoManager) GetNote(id string) (note localtypes.Note, err error) {
	result, err := db.Client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(NotesTableName),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		fmt.Println("Error getting note")
		fmt.Print(err)
		return
	}
	if result.Item == nil {
		return
	}
	err = attributevalue.UnmarshalMap(result.Item, &note)
	if err != nil {
		fmt.Println("Error unmarshalling note")
		fmt.Print(err)
		return
	}
	return
}

func (db DynamoManager) GetNotes(userID string) (notes []localtypes.Note, err error) {
	result, err := db.Client.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String(NotesTableName),
		IndexName:              aws.String("UserIDIndex"), // Specify the index name
		KeyConditionExpression: aws.String("UserID = :userID"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":userID": &types.AttributeValueMemberS{Value: userID},
		},
	})
	if err != nil {
		fmt.Println("Error getting notes")
		fmt.Println(err)
		return
	}
	for _, item := range result.Items {
		note := localtypes.Note{}
		err = attributevalue.UnmarshalMap(item, &note)
		if err != nil {
			fmt.Println("Error unmarshalling note")
			fmt.Println(err)
			return
		}
		notes = append(notes, note)
	}
	return
}

func (db DynamoManager) PatchNote(note localtypes.Note) (id string, err error) {
	marshalledNote, err := attributevalue.MarshalMap(note)
	if err != nil {
		fmt.Println("Error marshalling note")
		fmt.Print(err)
		return
	}
	_, err = db.Client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(NotesTableName), Item: marshalledNote,
	})

	if err != nil {
		fmt.Println("Error patching note")
		fmt.Print(err)
		return
	}

	return note.ID, err
}

func (db DynamoManager) DeleteNote(note localtypes.Note) (id string, err error) {
	_, err = db.Client.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String(NotesTableName),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: note.ID},
		},
	})
	if err != nil {
		fmt.Println("Error deleting note")
		fmt.Print(err)
		return
	}
	return note.ID, err
}

// to nie powinno byc tutaj jak co
func TestNotes(dynamoClient *DynamoManager) {
	// Create a new note
	note := localtypes.Note{
		ID:                 "2",
		Title:              "asd",
		Content:            "asd",
		DateOfCreation:     "123",
		DateOfModification: "123",
		IsShared:           true,
		UserID:             "2",
	}

	note2 := localtypes.Note{
		ID:                 "3",
		Title:              "asd",
		Content:            "asd",
		DateOfCreation:     "123",
		DateOfModification: "123",
		IsShared:           true,
		UserID:             "2",
	}

	// Insert the first note
	id, err := dynamoClient.InsertNote(note)
	if err != nil {
		fmt.Println("Error inserting the first note:", err)
		return
	}
	fmt.Println("Inserted first note with ID:", id)

	// Retrieve the first note
	retrievedNote, err := dynamoClient.GetNote(id)
	if err != nil {
		fmt.Println("Error retrieving the first note:", err)
		return
	}
	fmt.Println("Retrieved first note:", retrievedNote)

	// Patch the first note
	retrievedNote.Title = "new title"
	patchedNoteID, err := dynamoClient.PatchNote(retrievedNote)
	if err != nil {
		fmt.Println("Error patching the first note:", err)
		return
	}
	fmt.Println("Patched first note with ID:", patchedNoteID)

	// Retrieve the patched note
	patchedNote, err := dynamoClient.GetNote(patchedNoteID)
	if err != nil {
		fmt.Println("Error retrieving the patched note:", err)
		return
	}
	fmt.Println("Retrieved patched note:", patchedNote)

	// Insert the second note
	secondNoteID, err := dynamoClient.InsertNote(note2)
	if err != nil {
		fmt.Println("Error inserting the second note:", err)
		return
	}
	fmt.Println("Inserted second note with ID:", secondNoteID)

	// Retrieve notes by UserID
	notes, err := dynamoClient.GetNotes("2")
	if err != nil {
		fmt.Println("Error retrieving notes by UserID:", err)
		return
	}
	fmt.Println("Retrieved notes by UserID:", notes)

	// Delete the patched note
	_, err = dynamoClient.DeleteNote(patchedNote)
	if err != nil {
		fmt.Println("Error deleting the patched note:", err)
		return
	}
	fmt.Println("Deleted patched note with ID:", patchedNoteID)

	// Retrieve notes by UserID again
	notesAfterDeletion, err := dynamoClient.GetNotes("2")
	if err != nil {
		fmt.Println("Error retrieving notes by UserID after deletion:", err)
		return
	}
	fmt.Println("Retrieved notes by UserID after deletion:", notesAfterDeletion)
}
