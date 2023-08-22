// Package uc contains all the main logic related to use case layer
package uc

import (
	"net/http"

	"modak/send-notification/v1/internal"
)

// EmailServiceInterface interface for this service
type EmailServiceInterface interface {
	Send(recipient, subject, message string) error
}

// SendNotificationUC struct for this use case
type SendNotificationUC struct {
	EmailService EmailServiceInterface
}

// Handle main method with the logic to send notifications
func (uc *SendNotificationUC) Handle(notification internal.Notification) error {
	// send notification via email
	err := uc.EmailService.Send(
		notification.Recipient,
		notification.Type,
		notification.Message,
	)
	if err != nil {
		return &internal.GeneralError{
			Code:          internal.CodeNotificationError,
			ID:            internal.IDNotificationEmailNotSent,
			Message:       "Error sending email notification",
			StatusCode:    http.StatusInternalServerError,
			OriginalError: err,
		}
	}

	return err
}

// NewSendNotificationUC new instance of this use case
func NewSendNotificationUC(EmailService EmailServiceInterface) *SendNotificationUC {
	return &SendNotificationUC{
		EmailService: EmailService,
	}
}
