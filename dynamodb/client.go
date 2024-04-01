package dynamodb

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type DynamoManager struct {
	Client *dynamodb.Client
}

func NewDynamoManager() (*DynamoManager, error) {
	client, err := newClient("local")
	if err != nil {
		return nil, err
	}
	manager := &DynamoManager{Client: client}
	// manager.createTables()
	return manager, nil
}

// newclient constructs a new dynamodb client using a default configuration
// and a provided profile name (created via aws configure cmd).
func newClient(profile string) (*dynamodb.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("localhost"),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: "http://localhost:8000"}, nil
			})),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID: "abcd", SecretAccessKey: "a1b2c3", SessionToken: "",
				Source: "Mock credentials used above for local instance",
			},
		}),
	)
	if err != nil {
		return nil, err
	}

	c := dynamodb.NewFromConfig(cfg)
	return c, nil
}

func (db DynamoManager) CreateDynamoDBTable(
	tableName string, input *dynamodb.CreateTableInput,
) error {
	_, err := db.Client.CreateTable(context.TODO(), input)
	if err != nil {
		log.Printf("Failed to create table %v with error: %v\n", tableName, err)
	} else {
		waiter := dynamodb.NewTableExistsWaiter(db.Client)
		err = waiter.Wait(context.TODO(), &dynamodb.DescribeTableInput{
			TableName: aws.String(tableName)}, 5*time.Minute)
		if err != nil {
			log.Printf("Failed to wait on create table %v with error: %v\n", tableName, err)
		}
	}

	return err
}

func (db DynamoManager) createTables() {
	err := db.CreateDynamoDBTable(NotesTableName, NotesTableInput)
	if err != nil {
		fmt.Println("Error creating notes table")
		fmt.Println(err)
		panic(err)
		return
	}

	err = db.CreateDynamoDBTable(UsersTableName, UsersTableInput)

	if err != nil {
		fmt.Println("Error creating users table")
		fmt.Println(err)
		panic(err)
		return
	}
}
