// Package services contains all logic related to services
package services

import (
	"modak/send-notification/v1/internal/infraestructure"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
)

// EmailSource address information about sender
const EmailSource = "kahs_kevin@hotmail.com"

// EmailService struct for this service
type EmailService struct {
	client infraestructure.SESAPI
}

// Send sends an email using Amazon SES
func (s *EmailService) Send(recipient, subject, message string) error {
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{
				aws.String(recipient),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Text: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(message),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String(EmailSource),
	}

	_, err := s.client.SendEmail(input)

	return err
}

// NewEmailService creates a new instance of the email service
func NewEmailService(client infraestructure.SESAPI) *EmailService {
	return &EmailService{
		client: client,
	}
}
