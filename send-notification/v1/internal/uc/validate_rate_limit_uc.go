// Package uc contains all the main logic related to use case layer
package uc

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"modak/send-notification/v1/internal"

	"github.com/google/uuid"
)

// RateLimitRulesRepositoryInterface struct for this repository related to rules
type RateLimitRulesRepositoryInterface interface {
	GetByType(notificationType string) (*internal.RateLimitRule, error)
}

// RateLimitCacheRepositoryInterface struct for this repository related to cache
type RateLimitCacheRepositoryInterface interface {
	SetNotificationSentTimestamp(
		notificationType, email, timestamp, uuid string,
		ttl int64,
	) error
	CountNotificationsWithinInterval(notificationType, email string, intervalInMinutes int) (int, error)
}

// ValidateRateLimitUC struct for this use case
type ValidateRateLimitUC struct {
	rateLimitRulesRepository RateLimitRulesRepositoryInterface
	rateLimitCacheRepository RateLimitCacheRepositoryInterface
}

// Handle main method with the logic to validate the rules of rate limit
func (uc *ValidateRateLimitUC) Handle(notification internal.Notification) (bool, error) {
	// Get the rules for the current notification
	rule, err := uc.rateLimitRulesRepository.GetByType(notification.Type)
	if err != nil {
		return false, &internal.GeneralError{
			Code:          internal.CodeGeneralError,
			ID:            internal.IDGeneralError,
			Message:       "Error getting from rule repository (GetByType)",
			StatusCode:    http.StatusInternalServerError,
			OriginalError: err,
		}
	}

	// If the notification rule does not exist we return an alert error
	if rule == nil {
		return false, &internal.GeneralError{
			Code:          internal.CodeNotificationError,
			ID:            internal.IDNotificationTypeNotImplemented,
			Message:       fmt.Sprintf("Notification type '%s' not implemented", notification.Type),
			StatusCode:    http.StatusInternalServerError,
			OriginalError: err,
		}
	}

	// If the notification rule about limit is zero we can't send any notification due to rate limit
	if rule.NotificationsLimit <= 0 {
		return false, nil
	}

	// We check if the notification can be sent
	isAllowedToSendNotification, err := uc.CanSend(notification, *rule)
	if err != nil {
		return false, err
	}

	// Check if is possible to send the notification
	if isAllowedToSendNotification {
		currentTimestamp := time.Now().Unix()

		// Update the timestamp in the cache to know that this user already received a message
		err := uc.rateLimitCacheRepository.SetNotificationSentTimestamp(
			notification.Type,
			notification.Recipient,
			strconv.FormatInt(currentTimestamp, 10),
			fmt.Sprintf("%s", uuid.New()),
			currentTimestamp+int64(rule.IntervalInMinutes*int(time.Minute/time.Second)),
		)
		if err != nil {
			return false, &internal.GeneralError{
				Code:          internal.CodeGeneralError,
				ID:            internal.IDGeneralError,
				Message:       "Error saving in rule repository (SetNotificationSentTimestamp)",
				StatusCode:    http.StatusInternalServerError,
				OriginalError: err,
			}
		}

		// Notification sent successfully
		return true, nil
	}

	// Rate limit exceeded, no errors, no notifications sent
	return false, nil
}

// CanSend check if the notification can be sent following the rules of rate limit
func (uc *ValidateRateLimitUC) CanSend(notification internal.Notification, rule internal.RateLimitRule) (bool, error) {
	// If we are within the interval, we need to check how many notifications were sent in this interval
	count, err := uc.rateLimitCacheRepository.CountNotificationsWithinInterval(
		notification.Type,
		notification.Recipient,
		rule.IntervalInMinutes,
	)
	if err != nil {
		return false, &internal.GeneralError{
			Code:          internal.CodeGeneralError,
			ID:            internal.IDGeneralError,
			Message:       "Error getting from cache repository (CountNotificationsWithinInterval)",
			StatusCode:    http.StatusInternalServerError,
			OriginalError: err,
		}
	}

	// If the count of notifications sent is less than the allowed limit, we can send another one
	// Otherwise, we have exceeded the rate limit
	return count < rule.NotificationsLimit, nil
}

// NewValidateRateLimitUC new instance of this use case
func NewValidateRateLimitUC(
	rateLimitRulesRepository RateLimitRulesRepositoryInterface,
	rateLimitCacheRepository RateLimitCacheRepositoryInterface,
) *ValidateRateLimitUC {
	return &ValidateRateLimitUC{
		rateLimitRulesRepository: rateLimitRulesRepository,
		rateLimitCacheRepository: rateLimitCacheRepository,
	}
}
