// Package internal contains all the main logic
package internal

import (
	"errors"
	"net/http"
	"testing"

	"modak/send-notification/v1/internal/infraestructure"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

type mockValidateRateLimitUC struct {
	handleFunc func(notification Notification) (bool, error)
}

func (m *mockValidateRateLimitUC) Handle(notification Notification) (bool, error) {
	return m.handleFunc(notification)
}

type mockSendNotificationUC struct {
	handleFunc func(notification Notification) error
}

func (m *mockSendNotificationUC) Handle(notification Notification) error {
	return m.handleFunc(notification)
}

type mockLogger struct{}

func (m *mockLogger) Infof(format string, args ...interface{})  {}
func (m *mockLogger) Errorf(format string, args ...interface{}) {}
func (m *mockLogger) WithFields(fields ...interface{}) infraestructure.LoggerInterface {
	return m
}

func TestHandler_Handle(t *testing.T) {
	tests := []struct {
		name           string
		eventBody      string
		validateRateUC ValidateRateLimitUCInterface
		sendNotifUC    SendNotificationUCInterface
		wantStatusCode int
		wantErr        bool
	}{
		{
			name:           "unmarshal error",
			eventBody:      "{",
			validateRateUC: &mockValidateRateLimitUC{},
			sendNotifUC:    &mockSendNotificationUC{},
			wantStatusCode: http.StatusInternalServerError,
			wantErr:        true,
		},
		{
			name:      "validate rate limit error",
			eventBody: `{"notifications":[{"type":"test","recipient":"test@example.com","message":"Hello"}]}`,
			validateRateUC: &mockValidateRateLimitUC{
				handleFunc: func(notification Notification) (bool, error) {
					return false, errors.New("rate limit error")
				},
			},
			sendNotifUC:    &mockSendNotificationUC{},
			wantStatusCode: http.StatusInternalServerError,
			wantErr:        true,
		},
		{
			name:      "validate rate limit returns canSend=false",
			eventBody: `{"notifications":[{"type":"test","recipient":"test@example.com","message":"Hello"}]}`,
			validateRateUC: &mockValidateRateLimitUC{
				handleFunc: func(notification Notification) (bool, error) {
					return false, nil
				},
			},
			sendNotifUC:    &mockSendNotificationUC{},
			wantStatusCode: http.StatusOK,
			wantErr:        false,
		},
		{
			name:      "send notification error",
			eventBody: `{"notifications":[{"type":"test","recipient":"test@example.com","message":"Hello"}]}`,
			validateRateUC: &mockValidateRateLimitUC{
				handleFunc: func(notification Notification) (bool, error) {
					return true, nil
				},
			},
			sendNotifUC: &mockSendNotificationUC{
				handleFunc: func(notification Notification) error {
					return errors.New("send notification error")
				},
			},
			wantStatusCode: http.StatusInternalServerError,
			wantErr:        true,
		}, {
			name:      "successful notification send",
			eventBody: `{"notifications":[{"type":"test","recipient":"test@example.com","message":"Hello"}]}`,
			validateRateUC: &mockValidateRateLimitUC{
				handleFunc: func(notification Notification) (bool, error) {
					return true, nil
				},
			},
			sendNotifUC: &mockSendNotificationUC{
				handleFunc: func(notification Notification) error {
					return nil
				},
			},
			wantStatusCode: http.StatusOK,
			wantErr:        false,
		},
		{
			name:      "general error",
			eventBody: `{"notifications":[{"type":"test","recipient":"test@example.com","message":"Hello"}]}`,
			validateRateUC: &mockValidateRateLimitUC{
				handleFunc: func(notification Notification) (bool, error) {
					return false, nil
				},
			},
			sendNotifUC: &mockSendNotificationUC{
				handleFunc: func(notification Notification) error {
					return &GeneralError{
						Code:       CodeGeneralError,
						ID:         IDGeneralError,
						Message:    "An unexpected error occurred",
						StatusCode: http.StatusInternalServerError,
					}
				},
			},
			wantStatusCode: http.StatusOK,
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHandler(tt.validateRateUC, tt.sendNotifUC, &mockLogger{})
			event := events.APIGatewayProxyRequest{
				Body: tt.eventBody,
			}
			resp, err := h.Handle(event)
			assert.Equal(t, tt.wantStatusCode, resp.StatusCode)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
