// Package repositories contains all logic related to repositories
package repositories

import (
	"modak/send-notification/v1/internal"
	"modak/send-notification/v1/internal/infraestructure"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// RateLimitRulesRepository struct for this repository
type RateLimitRulesRepository struct {
	client    infraestructure.DynamoAPI
	tableName string
}

// GetByType get the records in database given a valid type
func (r *RateLimitRulesRepository) GetByType(notificationType string) (*internal.RateLimitRule, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"pk": {
				S: aws.String("TYPE#" + notificationType),
			},
		},
	}

	result, err := r.client.GetItem(input)
	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, nil
	}

	var rule internal.RateLimitRule

	err = dynamodbattribute.UnmarshalMap(result.Item, &rule)
	if err != nil {
		return nil, err
	}

	return &rule, nil
}

// NewRateLimitRulesRepository instance of a new repository
func NewRateLimitRulesRepository(
	client infraestructure.DynamoAPI,
	tableName string,
) *RateLimitRulesRepository {
	return &RateLimitRulesRepository{
		client:    client,
		tableName: tableName,
	}
}
