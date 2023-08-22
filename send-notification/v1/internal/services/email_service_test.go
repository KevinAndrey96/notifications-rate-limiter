// Package services contains all logic related to services
package services

import (
	"errors"
	"testing"

	"modak/send-notification/v1/internal/infraestructure"

	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/stretchr/testify/assert"
)

// mockSESAPI mock for SES API
type mockSESAPI struct {
	SendEmailFunc func(input *ses.SendEmailInput) (*ses.SendEmailOutput, error)
}

// SendEmail mock for this method to send email
func (m *mockSESAPI) SendEmail(input *ses.SendEmailInput) (*ses.SendEmailOutput, error) {
	return m.SendEmailFunc(input)
}

// TestEmailService_Send test for this method
func TestEmailService_Send(t *testing.T) {
	type fields struct {
		client infraestructure.SESAPI
	}

	type args struct {
		recipient string
		subject   string
		message   string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				client: &mockSESAPI{
					SendEmailFunc: func(input *ses.SendEmailInput) (*ses.SendEmailOutput, error) {
						return &ses.SendEmailOutput{}, nil
					},
				},
			},
			args: args{
				recipient: "test@example.com",
				subject:   "Test Subject",
				message:   "Test Message",
			},
			wantErr: false,
		},
		{
			name: "error sending email",
			fields: fields{
				client: &mockSESAPI{
					SendEmailFunc: func(input *ses.SendEmailInput) (*ses.SendEmailOutput, error) {
						return nil, assert.AnError
					},
				},
			},
			args: args{
				recipient: "test@example.com",
				subject:   "Test Subject",
				message:   "Test Message",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &EmailService{
				client: tt.fields.client,
			}
			err := s.Send(tt.args.recipient, tt.args.subject, tt.args.message)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestNewEmailService test for this service
func TestNewEmailService(t *testing.T) {
	type args struct {
		client infraestructure.SESAPI
	}

	tests := []struct {
		name          string
		args          args
		wantSendError bool
	}{
		{
			name: "success",
			args: args{
				client: &mockSESAPI{
					SendEmailFunc: func(input *ses.SendEmailInput) (*ses.SendEmailOutput, error) {
						return &ses.SendEmailOutput{}, nil
					},
				},
			},
			wantSendError: false,
		},
		{
			name: "send error",
			args: args{
				client: &mockSESAPI{
					SendEmailFunc: func(input *ses.SendEmailInput) (*ses.SendEmailOutput, error) {
						return nil, errors.New("send error")
					},
				},
			},
			wantSendError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewEmailService(tt.args.client)
			err := service.Send("test@example.com", "subject", "message")
			if (err != nil) != tt.wantSendError {
				t.Errorf("EmailService.Send() error = %v, wantSendError %v", err, tt.wantSendError)
			}
		})
	}
}
