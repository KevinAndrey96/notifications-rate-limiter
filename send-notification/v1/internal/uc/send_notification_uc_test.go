// Package uc contains all the main logic related to use case layer
package uc

import (
	"errors"
	"testing"

	"modak/send-notification/v1/internal"
)

// mockEmailService Mock for email service
type mockEmailService struct {
	SendFunc func(recipient, subject, message string) error
}

// Send Mock for method send of email service
func (m *mockEmailService) Send(recipient, subject, message string) error {
	return m.SendFunc(recipient, subject, message)
}

// TestSendNotificationUC_Handle test for this method
func TestSendNotificationUC_Handle(t *testing.T) {
	type fields struct {
		emailService EmailServiceInterface
	}

	type args struct {
		notification internal.Notification
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
				emailService: &mockEmailService{
					SendFunc: func(recipient, subject, message string) error {
						return nil
					},
				},
			},
			args: args{
				notification: internal.Notification{
					Type:      "News",
					Recipient: "test@example.com",
					Message:   "Notification about News",
				},
			},
			wantErr: false,
		},
		{
			name: "send email error",
			fields: fields{
				emailService: &mockEmailService{
					SendFunc: func(recipient, subject, message string) error {
						return errors.New("send error")
					},
				},
			},
			args: args{
				notification: internal.Notification{
					Type:      "News",
					Recipient: "test@example.com",
					Message:   "Notification about News",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ucInstance := &SendNotificationUC{
				EmailService: tt.fields.emailService,
			}
			err := ucInstance.Handle(tt.args.notification)
			if (err != nil) != tt.wantErr {
				t.Errorf("SendNotificationUC.Handle() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestNewSendNotificationUC Test for this method
func TestNewSendNotificationUC(t *testing.T) {
	type args struct {
		emailService EmailServiceInterface
	}

	tests := []struct {
		name string
		args args
		want *SendNotificationUC
	}{
		{
			name: "success",
			args: args{
				emailService: &mockEmailService{},
			},
			want: &SendNotificationUC{
				EmailService: &mockEmailService{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSendNotificationUC(tt.args.emailService); got.EmailService == nil {
				t.Errorf("NewSendNotificationUC().EmailService is nil, want not nil")
			}
		})
	}
}
