// Package di have all the injections dependency logic
package di

import (
	"os"

	"modak/send-notification/v1/internal/infraestructure"
	"modak/send-notification/v1/internal/repositories"
	"modak/send-notification/v1/internal/services"
	"modak/send-notification/v1/internal/uc"
)

// newAWSSessionProvider provider to aws session
func newAWSSessionProvider() infraestructure.SessionProvider {
	return infraestructure.NewSessionProvider(&infraestructure.SessionConfig{})
}

// newLoggerProvider provider for logger
func newLoggerProvider() infraestructure.LoggerInterface {
	return infraestructure.NewLogrusProvider().Logger()
}

// newDynamoDBProvider dynamo db provider
func newDynamoDBProvider(awsSession infraestructure.SessionProvider) infraestructure.DynamoAPI {
	dynamoProvider := infraestructure.NewDynamoProvider(awsSession, &infraestructure.DynamoConfig{})

	dynamoClient, err := dynamoProvider.DynamoClient()
	if err != nil {
		panic(err)
	}

	return dynamoClient
}

// newSESProvider creates and returns an Amazon SES client.
func newSESProvider(awsSession infraestructure.SessionProvider) infraestructure.SESAPI {
	sesProvider := infraestructure.NewSESProvider(awsSession, &infraestructure.SESConfig{})

	sesClient, err := sesProvider.SESClient()
	if err != nil {
		panic(err)
	}

	return sesClient
}

// newRateLimitRulesRepositoryProvider provider for this repository
func newRateLimitRulesRepositoryProvider(
	dynamoProvider infraestructure.DynamoAPI,
) uc.RateLimitRulesRepositoryInterface {
	return repositories.NewRateLimitRulesRepository(
		dynamoProvider,
		os.Getenv("DYNAMODB_NOTIFICATION_RATE_LIMIT_RULES_TABLE_NAME"),
	)
}

// newRateLimitCacheRepositoryProvider provider for this repository
func newRateLimitCacheRepositoryProvider(
	dynamoProvider infraestructure.DynamoAPI,
) uc.RateLimitCacheRepositoryInterface {
	return repositories.NewRateLimitCacheRepository(
		dynamoProvider,
		os.Getenv("DYNAMODB_NOTIFICATION_RATE_LIMIT_CACHE_TABLE_NAME"),
	)
}

// newEmailServiceProvider provider for this service
func newEmailServiceProvider(
	sesProvider infraestructure.SESAPI,
) uc.EmailServiceInterface {
	return services.NewEmailService(
		sesProvider,
	)
}
