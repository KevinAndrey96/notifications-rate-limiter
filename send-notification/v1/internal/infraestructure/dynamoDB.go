package infraestructure

import "github.com/aws/aws-sdk-go/service/dynamodb"

// DynamoAPI interface for DynamoDB methods.
type DynamoAPI interface {
	GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error)
	PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error)
	Query(input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error)
}

// DynamoProvider interface for Dynamo client.
type DynamoProvider interface {
	DynamoClient() (DynamoAPI, error)
}

// DynamoConfig struct with config for Dynamo.
type DynamoConfig struct{}

// Dynamo attributes required for DynamoProvider.
type Dynamo struct {
	client  *dynamodb.DynamoDB
	session SessionProvider
	config  *DynamoConfig
}

// DynamoClient create a new client for DynamoDB.
func (d *Dynamo) DynamoClient() (DynamoAPI, error) {
	if d.client == nil {
		dynamoSession, err := d.session.Session()
		if err != nil {
			return nil, err
		}
		d.client = dynamodb.New(dynamoSession)
	}

	return d.client, nil
}

// NewDynamoProvider instantiate new NewDynamoProvider.
func NewDynamoProvider(session SessionProvider, config *DynamoConfig) DynamoProvider {
	return &Dynamo{
		session: session,
		config:  config,
	}
}
