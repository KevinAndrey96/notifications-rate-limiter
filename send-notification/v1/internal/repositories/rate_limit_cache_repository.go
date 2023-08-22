// Package repositories contains all logic related to repositories
package repositories

import (
	"fmt"
	"time"

	"modak/send-notification/v1/internal/infraestructure"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// RateLimitCacheRepository struct for this repository
type RateLimitCacheRepository struct {
	client    infraestructure.DynamoAPI
	tableName string
}

// SetNotificationSentTimestamp Save in database a record to identify that this user was notified in that timestamp
func (r *RateLimitCacheRepository) SetNotificationSentTimestamp(
	notificationType, email, timestamp, uuid string,
	ttl int64,
) error {
	partitionKey := fmt.Sprintf("%s#%s", notificationType, email)
	sortKey := fmt.Sprintf("%s#%s", timestamp, uuid)

	input := &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item: map[string]*dynamodb.AttributeValue{
			"pk": {
				S: aws.String(partitionKey),
			},
			"sk": {
				S: aws.String(sortKey),
			},
			"ttl": {
				N: aws.String(fmt.Sprintf("%d", ttl)),
			},
		},
	}

	_, err := r.client.PutItem(input)

	return err
}

// CountNotificationsWithinInterval check the number of notifications that one user had in the given interval of time
func (r *RateLimitCacheRepository) CountNotificationsWithinInterval(
	notificationType, email string,
	intervalInMinutes int,
) (int, error) {
	startTimestamp := time.Now().Add(-time.Duration(intervalInMinutes) * time.Minute).Unix()

	partitionKey := fmt.Sprintf("%s#%s", notificationType, email)

	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		KeyConditionExpression: aws.String("pk = :pk AND sk >= :startRange"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":pk": {
				S: aws.String(partitionKey),
			},
			":startRange": {
				S: aws.String(fmt.Sprintf("%d#-", startTimestamp)),
			},
		},
	}

	result, err := r.client.Query(input)
	if err != nil {
		return 0, err
	}

	return int(*result.Count), nil
}

// NewRateLimitCacheRepository new instance of this repository
func NewRateLimitCacheRepository(
	client infraestructure.DynamoAPI,
	tableName string,
) *RateLimitCacheRepository {
	return &RateLimitCacheRepository{
		client:    client,
		tableName: tableName,
	}
}
