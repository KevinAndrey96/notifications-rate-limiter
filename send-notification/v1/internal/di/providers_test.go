// Package di have all the injections dependency logic
package di

import (
	"os"
	"reflect"
	"testing"

	"modak/send-notification/v1/internal/infraestructure"
	"modak/send-notification/v1/internal/repositories"
	"modak/send-notification/v1/internal/services"
	"modak/send-notification/v1/internal/uc"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

// Test_newAWSSessionProvider tests for this provider
func Test_newAWSSessionProvider(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		want infraestructure.SessionProvider
	}{
		{
			name: "ok",
			want: infraestructure.NewSessionProvider(&infraestructure.SessionConfig{}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := newAWSSessionProvider(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newAWSSessionProvider() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test_newDynamoDBProvider mock for this provider
func Test_newDynamoDBProvider(t *testing.T) {
	t.Parallel()

	type args struct {
		awsSession infraestructure.SessionProvider
	}

	tests := []struct {
		name string
		args args
		want func(a args) infraestructure.DynamoAPI
	}{
		{
			name: "success",
			args: args{
				awsSession: infraestructure.NewSessionProvider(&infraestructure.SessionConfig{Region: "us-east-1"}),
			},
			want: func(a args) infraestructure.DynamoAPI {
				dynamoProvider := infraestructure.NewDynamoProvider(a.awsSession, &infraestructure.DynamoConfig{})
				dynamoClient, err := dynamoProvider.DynamoClient()
				if err != nil {
					panic(err)
				}
				return dynamoClient
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := newDynamoDBProvider(
				infraestructure.NewSessionProvider(&infraestructure.SessionConfig{Region: "us-east-1"}),
			); reflect.DeepEqual(got, tt.want) {
				t.Errorf(
					"newDynamoDBProvider() = %v, want %v",
					got,
					tt.want(tt.args),
				)
			}
		})
	}
}

// mockSESProvider mock for ses provider
type mockSESProvider struct{}

// SendEmail mock for the method SendEmail
func (m *mockSESProvider) SendEmail(input *ses.SendEmailInput) (*ses.SendEmailOutput, error) {
	return &ses.SendEmailOutput{}, nil
}

// Test_newEmailServiceProvider tests for this provider
func Test_newEmailServiceProvider(t *testing.T) {
	t.Parallel()

	type args struct {
		sesProvider infraestructure.SESAPI
	}

	tests := []struct {
		name string
		args args
		want func(a args) uc.EmailServiceInterface
	}{
		{
			name: "success",
			args: args{
				sesProvider: &mockSESProvider{},
			},
			want: func(a args) uc.EmailServiceInterface {
				return services.NewEmailService(a.sesProvider)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := newEmailServiceProvider(tt.args.sesProvider); !reflect.DeepEqual(got, tt.want(tt.args)) {
				t.Errorf("newEmailServiceProvider() = %v, want %v", got, tt.want(tt.args))
			}
		})
	}
}

// Test_newLoggerProvider test for New logger Provider
func Test_newLoggerProvider(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		want infraestructure.LoggerInterface
	}{
		{
			name: "new logger success",
			want: infraestructure.NewLogrusProvider().Logger(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := newLoggerProvider()
			if reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf("newLoggerProvider() = %v, want %v", reflect.TypeOf(got), reflect.TypeOf(tt.want))
			}
		})
	}
}

// Test_newRateLimitCacheRepositoryProvider tests for this provider
func Test_newRateLimitCacheRepositoryProvider(t *testing.T) {
	t.Parallel()

	err := os.Setenv(
		"DYNAMODB_NOTIFICATION_RATE_LIMIT_CACHE_TABLE_NAME",
		"prod-notification-rate-limit-cache",
	)
	if err != nil {
		panic(err)
	}

	type args struct {
		dynamoProvider infraestructure.DynamoAPI
	}

	tests := []struct {
		name string
		args args
		want func(a args) uc.RateLimitCacheRepositoryInterface
	}{
		{
			name: "success",
			args: args{
				dynamoProvider: newDynamoDBProvider(
					infraestructure.NewSessionProvider(&infraestructure.SessionConfig{}),
				),
			},
			want: func(a args) uc.RateLimitCacheRepositoryInterface {
				return repositories.NewRateLimitCacheRepository(
					a.dynamoProvider,
					os.Getenv("DYNAMODB_NOTIFICATION_RATE_LIMIT_CACHE_TABLE_NAME"),
				)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := newRateLimitCacheRepositoryProvider(
				tt.args.dynamoProvider,
			); !reflect.DeepEqual(got, tt.want(tt.args)) {
				t.Errorf("newRateLimitCacheRepositoryProvider() = %v, want %v", got, tt.want(tt.args))
			}
		})
	}
}

// Test_newRateLimitRulesRepositoryProvider Tests for this provider
func Test_newRateLimitRulesRepositoryProvider(t *testing.T) {
	t.Parallel()

	err := os.Setenv(
		"DYNAMODB_NOTIFICATION_RATE_LIMIT_RULES_TABLE_NAME",
		"prod-notification-rate-limit-rules",
	)
	if err != nil {
		panic(err)
	}

	type args struct {
		dynamoProvider infraestructure.DynamoAPI
	}

	tests := []struct {
		name string
		args args
		want func(a args) uc.RateLimitRulesRepositoryInterface
	}{
		{
			name: "success",
			args: args{
				dynamoProvider: newDynamoDBProvider(
					infraestructure.NewSessionProvider(&infraestructure.SessionConfig{}),
				),
			},
			want: func(a args) uc.RateLimitRulesRepositoryInterface {
				return repositories.NewRateLimitRulesRepository(
					a.dynamoProvider,
					os.Getenv("DYNAMODB_NOTIFICATION_RATE_LIMIT_RULES_TABLE_NAME"),
				)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := newRateLimitRulesRepositoryProvider(
				tt.args.dynamoProvider,
			); !reflect.DeepEqual(got, tt.want(tt.args)) {
				t.Errorf("newRateLimitRulesRepositoryProvider() = %v, want %v", got, tt.want(tt.args))
			}
		})
	}
}

// mockSessionProvider mock for session provider
type mockSessionProvider struct{}

// Session mock for the method Session
func (m *mockSessionProvider) Session() (client.ConfigProvider, error) {
	return session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
	})
}

// Test_newSESProvider test for this method
func Test_newSESProvider(t *testing.T) {
	t.Parallel()

	type args struct {
		awsSession infraestructure.SessionProvider
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "success",
			args: args{
				awsSession: &mockSessionProvider{},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := newSESProvider(tt.args.awsSession)
			if (got != nil) != tt.want {
				t.Errorf("newSESProvider() implemented SESAPI = %v, want %v", (got != nil), tt.want)
			}
		})
	}
}
