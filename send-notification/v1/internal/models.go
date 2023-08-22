// Package internal contains all the main logic
package internal

// RequestBody struct for request body
type RequestBody struct {
	Notifications []Notification `json:"notifications"`
}

// ResponseBody struct for response body
type ResponseBody struct {
	Sent   []Notification `json:"sent"`
	Failed []Notification `json:"failed"`
}

// Notification model for notification sent
type Notification struct {
	Type      string `json:"type"`
	Recipient string `json:"recipient"`
	Message   string `json:"message"`
}

// RateLimitRule model for rate limit rules stored in database
type RateLimitRule struct {
	PK                 string `dynamodbav:"pk"`
	NotificationsLimit int    `dynamodbav:"notifications_limit"`
	IntervalInMinutes  int    `dynamodbav:"interval_in_minutes"`
}
