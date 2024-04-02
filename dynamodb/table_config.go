package dynamodb

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const NotesTableName string = "NotesTable"
const UsersTableName string = "UsersTable"

var NotesTableInput *dynamodb.CreateTableInput = &dynamodb.CreateTableInput{
	AttributeDefinitions: []types.AttributeDefinition{
		{
			AttributeName: aws.String("ID"),
			AttributeType: types.ScalarAttributeTypeS,
		},
		{
			AttributeName: aws.String("UserID"),
			AttributeType: types.ScalarAttributeTypeS,
		},
	},
	KeySchema: []types.KeySchemaElement{
		{
			AttributeName: aws.String("ID"),
			KeyType:       types.KeyTypeHash,
		},
	},
	GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{
		{
			IndexName: aws.String("UserIDIndex"),
			KeySchema: []types.KeySchemaElement{
				{
					AttributeName: aws.String("UserID"),
					KeyType:       types.KeyTypeHash,
				},
			},
			Projection: &types.Projection{
				ProjectionType: types.ProjectionTypeAll,
			},
		},
	},
	TableName:   aws.String(NotesTableName),
	BillingMode: types.BillingModePayPerRequest,
}

var UsersTableInput *dynamodb.CreateTableInput = &dynamodb.CreateTableInput{
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
	TableName:   aws.String(UsersTableName),
	BillingMode: types.BillingModePayPerRequest,
}
