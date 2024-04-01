package dynamodb

import (
	"context"
	"fmt"

	localtypes "database-tester/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

func (db DynamoManager) InsertUser(user localtypes.User) (id string, err error) {
	user.ID = uuid.New().String()
	marshalledUser, err := attributevalue.MarshalMap(user)
	if err != nil {
		fmt.Println("Error marshalling user")
		fmt.Print(err)
		return
	}
	_, err = db.Client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(UsersTableName),
		Item:      marshalledUser,
	})

	if err != nil {
		fmt.Println("Error inserting user")
		fmt.Print(err)
		return
	}

	return user.ID, nil
}

func (db DynamoManager) GetUser(id string) (user localtypes.User, err error) {
	result, err := db.Client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(UsersTableName),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		fmt.Println("Error getting user")
		fmt.Print(err)
		return
	}
	if result.Item == nil {
		return
	}
	err = attributevalue.UnmarshalMap(result.Item, &user)
	if err != nil {
		fmt.Println("Error unmarshalling user")
		fmt.Print(err)
		return
	}
	return
}

func (db DynamoManager) DeleteUser(user localtypes.User) (id string, err error) {
	_, err = db.Client.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String(UsersTableName),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: user.ID},
		},
	})
	if err != nil {
		fmt.Println("Error deleting user")
		fmt.Print(err)
		return
	}
	return user.ID, nil
}

func (db DynamoManager) PatchUser(user localtypes.User) (id string, err error) {
	marshalledUser, err := attributevalue.MarshalMap(user)
	if err != nil {
		fmt.Println("Error marshalling user")
		fmt.Print(err)
		return
	}
	_, err = db.Client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(UsersTableName),
		Item:      marshalledUser,
	})

	if err != nil {
		fmt.Println("Error patching user")
		fmt.Print(err)
		return
	}

	return user.ID, nil
}

// to nie powinno byc tutaj jak co
func TestUsers(dynamoClient *DynamoManager) {
	// Insert a new user
	user := localtypes.User{
		ID:        "1",
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Username:  "johndoe",
		IsAdmin:   false,
		IsEnabled: true,
	}

	userID, err := dynamoClient.InsertUser(user)
	if err != nil {
		fmt.Println("Error inserting user")
		fmt.Println(err)
		return
	}
	fmt.Println("User ID:", userID)

	// Get a user by ID
	retrievedUser, err := dynamoClient.GetUser(userID)
	if err != nil {
		fmt.Println("Error getting user")
		fmt.Println(err)
		return
	}
	fmt.Println("Retrieved User:", retrievedUser)

	// Patch a user
	retrievedUser.FirstName = "Turbinator" // Modify some attributes
	retrievedUser.LastName = "Masakrator"

	patchedUserID, err := dynamoClient.PatchUser(retrievedUser)
	if err != nil {
		fmt.Println("Error patching user")
		fmt.Println(err)
		return
	}
	fmt.Println("Patched User ID:", patchedUserID)

	// Retrieve the patched user
	patchedUser, err := dynamoClient.GetUser(patchedUserID)
	if err != nil {
		fmt.Println("Error getting patched user")
		fmt.Println(err)
		return
	}
	fmt.Println("Patched User:", patchedUser)

	// Delete a user
	deletedUserID, err := dynamoClient.DeleteUser(retrievedUser)
	if err != nil {
		fmt.Println("Error deleting user")
		fmt.Println(err)
		return
	}
	fmt.Println("Deleted User ID:", deletedUserID)

	// Retrieve the deleted user
	deletedUser, err := dynamoClient.GetUser(patchedUserID)
	if err != nil {
		fmt.Println("Error getting deleted user")
		fmt.Println(err)
		return
	}
	fmt.Println("deleted User:", deletedUser)
}
