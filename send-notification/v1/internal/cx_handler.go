// Package internal contains all the main logic
package internal

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"modak/send-notification/v1/internal/infraestructure"

	"github.com/aws/aws-lambda-go/events"
)

// ValidateRateLimitUCInterface interface for this use case validate rate limit
type ValidateRateLimitUCInterface interface {
	Handle(notification Notification) (bool, error)
}

// SendNotificationUCInterface interface for this use case validate rate limit
type SendNotificationUCInterface interface {
	Handle(notification Notification) error
}

// Handler declaration of handler struct used in this file
type Handler struct {
	validateRateLimitUC ValidateRateLimitUCInterface
	sendNotificationUC  SendNotificationUCInterface
	logger              infraestructure.LoggerInterface
}

// Handle main method controller to execute this lambda function
func (h *Handler) Handle(event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Init logger with light ECS specification
	logger := h.logger.WithFields(
		"@timestamp", time.Now().Format(time.RFC3339),
		"file", "cx_handler",
		"method", "Handle",
	)

	var requestBody RequestBody

	err := json.Unmarshal([]byte(event.Body), &requestBody)
	if err != nil {
		logger.Errorf("error: ", err)

		return responseError(err)
	}

	var sent, failed []Notification

	// Create channels to handle concurrency
	sentChannel := make(chan Notification, len(requestBody.Notifications))
	failedChannel := make(chan Notification, len(requestBody.Notifications))
	errorsChannel := make(chan error, len(requestBody.Notifications))

	// Process notifications concurrently
	for _, notification := range requestBody.Notifications {
		go func(notification Notification) {
			canSend, err := h.validateRateLimitUC.Handle(notification)
			if err != nil {
				errorsChannel <- err

				return
			}

			if !canSend {
				failedChannel <- notification

				return
			}

			err = h.sendNotificationUC.Handle(notification)
			if err != nil {
				errorsChannel <- err

				return
			}
			sentChannel <- notification
		}(notification)
	}

	// Collect results
	for i := 0; i < len(requestBody.Notifications); i++ {
		select {
		case notification := <-sentChannel:
			sent = append(sent, notification)
		case notification := <-failedChannel:
			failed = append(failed, notification)
		case err := <-errorsChannel:
			logger.Errorf("error: ", err)

			return responseError(err)
		}
	}

	responseBody := ResponseBody{
		Sent:   sent,
		Failed: failed,
	}

	jsonData, err := json.Marshal(responseBody)
	if err != nil {
		logger.Errorf("error: ", err)

		return responseError(err)
	}

	logger.Infof("Notifications processed successfully. Sent %d, Failed %d", len(sent), len(failed))

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(jsonData),
	}, nil
}

// responseError return response according error type
func responseError(err error) (events.APIGatewayProxyResponse, error) {
	var lambdaError error

	var httpStatusCode int

	errors := new(ErrorsJSONAPI)

	switch e := err.(type) {
	case *GeneralError:
		errors.Add(ErrorJSONAPI{
			Status: strconv.Itoa(e.StatusCode),
			Code:   e.Code,
			ID:     e.ID,
			Title:  GeneralErrorTitle,
			Detail: err.Error(),
		})

		httpStatusCode = e.StatusCode
	default:
		lambdaError = e

		errors.Add(ErrorJSONAPI{
			Status: strconv.Itoa(http.StatusInternalServerError),
			Code:   CodeGeneralError,
			ID:     IDGeneralError,
			Title:  GeneralErrorTitle,
			Detail: e.Error(),
		})

		httpStatusCode = http.StatusInternalServerError
	}

	errorsResponse, _ := json.Marshal(errors)

	return events.APIGatewayProxyResponse{
		StatusCode: httpStatusCode,
		Body:       string(errorsResponse),
	}, lambdaError
}

// NewHandler Initialize Handle
func NewHandler(
	validateRateLimitUC ValidateRateLimitUCInterface,
	sendNotificationUC SendNotificationUCInterface,
	logger infraestructure.LoggerInterface,
) *Handler {
	return &Handler{
		validateRateLimitUC: validateRateLimitUC,
		sendNotificationUC:  sendNotificationUC,
		logger:              logger,
	}
}
