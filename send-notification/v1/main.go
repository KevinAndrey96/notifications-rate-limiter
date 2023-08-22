// Package main have the logic necessary to deploy the main handler
package main

import (
	"modak/send-notification/v1/internal/di"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	handler, err := di.Initialize()
	if err != nil {
		panic("fatal err: " + err.Error())
	}
	lambda.Start(handler.Handle)
}
