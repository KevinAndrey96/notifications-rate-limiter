//go:build wireinject
// +build wireinject

// Package di contains all the logic related to di
package di

import (
	"modak/send-notification/v1/internal"

	"github.com/google/wire"
)

// Initialize method to initialize wire
func Initialize() (*internal.Handler, error) {
	wire.Build(stdSet)
	return &internal.Handler{}, nil
}
