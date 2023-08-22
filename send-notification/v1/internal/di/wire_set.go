//go:build wireinject
// +build wireinject

// Package di contains all the logic related to di
package di

import (
	"modak/send-notification/v1/internal"
	"modak/send-notification/v1/internal/uc"

	"github.com/google/wire"
)

var stdSet = wire.NewSet(
	newAWSSessionProvider,
	newLoggerProvider,
	newDynamoDBProvider,
	newSESProvider,
	internal.NewHandler,
	newRateLimitRulesRepositoryProvider,
	newRateLimitCacheRepositoryProvider,
	newEmailServiceProvider,

	uc.NewValidateRateLimitUC,
	wire.Bind(new(internal.ValidateRateLimitUCInterface), new(*uc.ValidateRateLimitUC)),
	uc.NewSendNotificationUC,
	wire.Bind(new(internal.SendNotificationUCInterface), new(*uc.SendNotificationUC)),
)
