package dynamodb

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const NotesTableName string = "NotesTableNameTable"

var TableInput *dynamodb.CreateTableInput = &dynamodb.CreateTableInput{
	AttributeDefinitions: []types.AttributeDefinition{
		{
			AttributeName: aws.String("ID"),
			AttributeType: types.ScalarAttributeTypeS,
		},
	},
	KeySchema: []types.KeySchemaElement{
		{
			AttributeName: aws.String("ID"),
			KeyType:       types.KeyTypeHash,
		},
	},
	TableName:   aws.String(NotesTableName),
	BillingMode: types.BillingModePayPerRequest,
}
