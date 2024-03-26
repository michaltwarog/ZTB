package dynamodb

import (
	"context"
	localtypes "database-tester/types"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

type DynamoManager struct {
	Client *dynamodb.Client
}

func NewDynamoManager(profile string) (*DynamoManager, error) {
	client, err := newClient(profile)
	if err != nil {
		return nil, err
	}
	return &DynamoManager{Client: client}, nil
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
	var tableDesc *types.TableDescription
	table, err := db.Client.CreateTable(context.TODO(), input)
	if err != nil {
		log.Printf("Failed to create table %v with error: %v\n", tableName, err)
	} else {
		waiter := dynamodb.NewTableExistsWaiter(db.Client)
		err = waiter.Wait(context.TODO(), &dynamodb.DescribeTableInput{
			TableName: aws.String(tableName)}, 5*time.Minute)
		if err != nil {
			log.Printf("Failed to wait on create table %v with error: %v\n", tableName, err)
		}
		tableDesc = table.TableDescription
	}

	fmt.Println(tableDesc)

	return err
}

func (db DynamoManager) InsertNote(note localtypes.Note) (id string, err error) {
	note.ID = uuid.New().String()
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
func (db DynamoManager) GetNote(id string) (note localtypes.Note, found bool, err error) {
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
	found = true
	return
}
