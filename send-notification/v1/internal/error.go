// Package internal have all the main logic
package internal

// List for code errors related to notifications
const (
	// CodeGeneralError Unexpected errors code
	CodeGeneralError string = "CODE_GENERAL_ERROR"
	// IDGeneralError Unexpected errors ID
	IDGeneralError string = "ID_GENERAL_ERROR"
	// GeneralErrorTitle title for all errors
	GeneralErrorTitle string = "Error"
	// CodeNotificationError this code represents a problem with the notification
	CodeNotificationError string = "CODE_NOTIFICATION_ERROR"
	// IDNotificationTypeNotImplemented this identifier is used when a type is not implemented
	IDNotificationTypeNotImplemented string = "ID_NOTIFICATION_NOT_IMPLEMENTED"
	// IDNotificationEmailNotSent this identifier is used when an email was not sent
	IDNotificationEmailNotSent string = "ID_NOTIFICATION_EMAIL_NOT_SENT"
)

// GeneralError for unexpected errors
type GeneralError struct {
	Code          string
	ID            string
	Message       string
	StatusCode    int
	OriginalError error
}

// Error get the error message
func (e *GeneralError) Error() string {
	return e.Message
}

// ErrorJSONAPI struct base from error response
type ErrorJSONAPI struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Code   string `json:"code"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

// ErrorsJSONAPIProvider interface to add or get errors
type ErrorsJSONAPIProvider interface {
	Add(jsonError ErrorJSONAPI) *ErrorsJSONAPI
	Get() *ErrorsJSONAPI
}

// ErrorsJSONAPI list of errors
type ErrorsJSONAPI struct {
	Errors []ErrorJSONAPI `json:"errors"`
}

// Add an error
func (e *ErrorsJSONAPI) Add(jsonError ErrorJSONAPI) *ErrorsJSONAPI {
	e.Errors = append(e.Errors, jsonError)

	return e
}

// Get an error
func (e *ErrorsJSONAPI) Get() *ErrorsJSONAPI {
	return e
}
