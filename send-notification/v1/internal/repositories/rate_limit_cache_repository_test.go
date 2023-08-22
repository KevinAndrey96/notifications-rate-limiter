// Package repositories contains all logic related to repositories
package repositories

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// mockDynamoAPI mock for dynamoAPI
type mockDynamoAPI struct {
	PutItemFunc func(*dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error)
	QueryFunc   func(*dynamodb.QueryInput) (*dynamodb.QueryOutput, error)
	GetItemFunc func(*dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error)
}

// PutItem insert a new item into dynamoDB
func (m *mockDynamoAPI) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	return m.PutItemFunc(input)
}

// Query get elements from dynamo given a query
func (m *mockDynamoAPI) Query(input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	return m.QueryFunc(input)
}

// GetItem get only one element from dynamoDB
func (m *mockDynamoAPI) GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	return m.GetItemFunc(input)
}

// TestRateLimitCacheRepository_SetNotificationSentTimestamp test for this method
func TestRateLimitCacheRepository_SetNotificationSentTimestamp(t *testing.T) {
	tests := []struct {
		name    string
		mock    *mockDynamoAPI
		wantErr bool
	}{
		{
			name: "success",
			mock: &mockDynamoAPI{
				PutItemFunc: func(*dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
					return &dynamodb.PutItemOutput{}, nil
				},
			},
			wantErr: false,
		},
		{
			name: "error on put",
			mock: &mockDynamoAPI{
				PutItemFunc: func(*dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
					return nil, errors.New("error on put")
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRateLimitCacheRepository(tt.mock, "test-table")
			err := r.SetNotificationSentTimestamp("testType", "test@email.com", "1234567890", "testUUID", 1234567890)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetNotificationSentTimestamp() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestRateLimitCacheRepository_CountNotificationsWithinInterval test for this method
func TestRateLimitCacheRepository_CountNotificationsWithinInterval(t *testing.T) {
	tests := []struct {
		name    string
		mock    *mockDynamoAPI
		want    int
		wantErr bool
	}{
		{
			name: "success",
			mock: &mockDynamoAPI{
				QueryFunc: func(*dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
					return &dynamodb.QueryOutput{Count: aws.Int64(5)}, nil
				},
			},
			want:    5,
			wantErr: false,
		},
		{
			name: "error on query",
			mock: &mockDynamoAPI{
				QueryFunc: func(*dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
					return nil, errors.New("error on query")
				},
			},
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRateLimitCacheRepository(tt.mock, "test-table")
			got, err := r.CountNotificationsWithinInterval("testType", "test@email.com", 10)
			if (err != nil) != tt.wantErr {
				t.Errorf("CountNotificationsWithinInterval() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CountNotificationsWithinInterval() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestNewRateLimitCacheRepository test for this repository
func TestNewRateLimitCacheRepository(t *testing.T) {
	client := &mockDynamoAPI{}
	tableName := "test-table"
	r := NewRateLimitCacheRepository(client, tableName)
	if r == nil || r.client != client || r.tableName != tableName {
		t.Errorf("NewRateLimitCacheRepository() did not initialize correctly")
	}
}
