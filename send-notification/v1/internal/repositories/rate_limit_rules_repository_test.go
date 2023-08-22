// Package repositories contains all logic related to repositories
package repositories

import (
	"errors"
	"reflect"
	"testing"

	"modak/send-notification/v1/internal"
	"modak/send-notification/v1/internal/infraestructure"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// TestRateLimitRulesRepository_GetByType test for this method
func TestRateLimitRulesRepository_GetByType(t *testing.T) {
	type fields struct {
		client    infraestructure.DynamoAPI
		tableName string
	}

	type args struct {
		notificationType string
	}

	rule := internal.RateLimitRule{
		PK:                 "TYPE#testType",
		NotificationsLimit: 5,
		IntervalInMinutes:  10,
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *internal.RateLimitRule
		wantErr bool
		mock    func(f fields, a args)
	}{
		{
			name: "error fetching data",
			fields: fields{
				client:    &mockDynamoAPI{},
				tableName: "rate-limit-rules",
			},
			args: args{
				notificationType: "testType",
			},
			wantErr: true,
			mock: func(f fields, a args) {
				_ = &dynamodb.GetItemInput{
					TableName: aws.String(f.tableName),
					Key: map[string]*dynamodb.AttributeValue{
						"pk": {
							S: aws.String("TYPE#" + a.notificationType),
						},
					},
				}
				f.client.(*mockDynamoAPI).GetItemFunc = func(*dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
					return nil, errors.New("error fetching data")
				}
			},
		},
		{
			name: "success",
			fields: fields{
				client:    &mockDynamoAPI{},
				tableName: "rate-limit-rules",
			},
			args: args{
				notificationType: "testType",
			},
			want: &rule,
			mock: func(f fields, a args) {
				_ = &dynamodb.GetItemInput{
					TableName: aws.String(f.tableName),
					Key: map[string]*dynamodb.AttributeValue{
						"pk": {
							S: aws.String("TYPE#" + a.notificationType),
						},
					},
				}
				r, _ := dynamodbattribute.MarshalMap(rule)
				f.client.(*mockDynamoAPI).GetItemFunc = func(*dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
					return &dynamodb.GetItemOutput{Item: r}, nil
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.fields, tt.args)
			r := &RateLimitRulesRepository{
				client:    tt.fields.client,
				tableName: tt.fields.tableName,
			}
			got, err := r.GetByType(tt.args.notificationType)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetByType() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestNewRateLimitRulesRepository tests for this repository
func TestNewRateLimitRulesRepository(t *testing.T) {
	rateLimitRulesRepository := NewRateLimitRulesRepository(&mockDynamoAPI{}, "rate-limit-rules")

	type args struct {
		client    infraestructure.DynamoAPI
		tableName string
	}

	tests := []struct {
		name string
		args args
		want *RateLimitRulesRepository
	}{
		{
			name: "success",
			args: args{
				client:    &mockDynamoAPI{},
				tableName: "rate-limit-rules",
			},
			want: rateLimitRulesRepository,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRateLimitRulesRepository(tt.args.client, tt.args.tableName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRateLimitRulesRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}
