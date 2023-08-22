// Package uc contains all the main logic related to use case layer
package uc

import (
	"errors"
	"testing"

	"modak/send-notification/v1/internal"

	"github.com/stretchr/testify/assert"
)

// MockRateLimitRulesRepository mock for repository with rate limit rules
type MockRateLimitRulesRepository struct {
	GetByTypeFunc func(notificationType string) (*internal.RateLimitRule, error)
}

// GetByType mock for the method that get the rules about rate limit
func (m *MockRateLimitRulesRepository) GetByType(notificationType string) (*internal.RateLimitRule, error) {
	return m.GetByTypeFunc(notificationType)
}

// MockRateLimitCacheRepository mock for repository with the cache of notifications
type MockRateLimitCacheRepository struct {
	SetNotificationSentTimestampFunc     func(notificationType, email, timestamp, uuid string, ttl int64) error
	CountNotificationsWithinIntervalFunc func(notificationType, email string, intervalInMinutes int) (int, error)
}

// SetNotificationSentTimestamp Mock for the method that save into the cache
func (m *MockRateLimitCacheRepository) SetNotificationSentTimestamp(
	notificationType,
	email,
	timestamp,
	uuid string,
	ttl int64,
) error {
	return m.SetNotificationSentTimestampFunc(notificationType, email, timestamp, uuid, ttl)
}

// CountNotificationsWithinInterval Mock for the method that count the number of notifications sent to a user
func (m *MockRateLimitCacheRepository) CountNotificationsWithinInterval(
	notificationType,
	email string,
	intervalInMinutes int,
) (int, error) {
	return m.CountNotificationsWithinIntervalFunc(notificationType, email, intervalInMinutes)
}

// TestValidateRateLimitUC_Handle Test for this method
func TestValidateRateLimitUC_Handle(t *testing.T) {
	rule := &internal.RateLimitRule{
		NotificationsLimit: 5,
		IntervalInMinutes:  10,
	}

	notification := internal.Notification{
		Type:      "test",
		Recipient: "test@example.com",
		Message:   "Hello",
	}

	tests := []struct {
		name          string
		rulesRepoFunc func() *MockRateLimitRulesRepository
		cacheRepoFunc func() *MockRateLimitCacheRepository
		want          bool
		wantErr       bool
	}{
		{
			name: "successful notification send",
			rulesRepoFunc: func() *MockRateLimitRulesRepository {
				return &MockRateLimitRulesRepository{
					GetByTypeFunc: func(notificationType string) (*internal.RateLimitRule, error) {
						return rule, nil
					},
				}
			},
			cacheRepoFunc: func() *MockRateLimitCacheRepository {
				return &MockRateLimitCacheRepository{
					CountNotificationsWithinIntervalFunc: func(
						notificationType,
						email string,
						intervalInMinutes int,
					) (int, error) {
						return 3, nil
					},
					SetNotificationSentTimestampFunc: func(
						notificationType,
						email,
						timestamp,
						uuid string,
						ttl int64,
					) error {
						return nil
					},
				}
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "error getting rule by type",
			rulesRepoFunc: func() *MockRateLimitRulesRepository {
				return &MockRateLimitRulesRepository{
					GetByTypeFunc: func(notificationType string) (*internal.RateLimitRule, error) {
						return nil, errors.New("database error")
					},
				}
			},
			cacheRepoFunc: func() *MockRateLimitCacheRepository {
				return &MockRateLimitCacheRepository{}
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "rule is nil",
			rulesRepoFunc: func() *MockRateLimitRulesRepository {
				return &MockRateLimitRulesRepository{
					GetByTypeFunc: func(notificationType string) (*internal.RateLimitRule, error) {
						return nil, nil
					},
				}
			},
			cacheRepoFunc: func() *MockRateLimitCacheRepository {
				return &MockRateLimitCacheRepository{}
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "notifications limit is zero",
			rulesRepoFunc: func() *MockRateLimitRulesRepository {
				return &MockRateLimitRulesRepository{
					GetByTypeFunc: func(notificationType string) (*internal.RateLimitRule, error) {
						return &internal.RateLimitRule{
							NotificationsLimit: 0,
						}, nil
					},
				}
			},
			cacheRepoFunc: func() *MockRateLimitCacheRepository {
				return &MockRateLimitCacheRepository{}
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "canSend returns false",
			rulesRepoFunc: func() *MockRateLimitRulesRepository {
				return &MockRateLimitRulesRepository{
					GetByTypeFunc: func(notificationType string) (*internal.RateLimitRule, error) {
						return &internal.RateLimitRule{
							NotificationsLimit: 5,
							IntervalInMinutes:  10,
						}, nil
					},
				}
			},
			cacheRepoFunc: func() *MockRateLimitCacheRepository {
				return &MockRateLimitCacheRepository{
					CountNotificationsWithinIntervalFunc: func(
						notificationType,
						email string,
						intervalInMinutes int,
					) (int, error) {
						return 5, nil
					},
				}
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "error from CountNotificationsWithinInterval",
			rulesRepoFunc: func() *MockRateLimitRulesRepository {
				return &MockRateLimitRulesRepository{
					GetByTypeFunc: func(notificationType string) (*internal.RateLimitRule, error) {
						return &internal.RateLimitRule{
							NotificationsLimit: 5,
							IntervalInMinutes:  10,
						}, nil
					},
				}
			},
			cacheRepoFunc: func() *MockRateLimitCacheRepository {
				return &MockRateLimitCacheRepository{
					CountNotificationsWithinIntervalFunc: func(
						notificationType,
						email string,
						intervalInMinutes int,
					) (int, error) {
						return 0, errors.New("cache retrieval error")
					},
				}
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "error from SetNotificationSentTimestamp",
			rulesRepoFunc: func() *MockRateLimitRulesRepository {
				return &MockRateLimitRulesRepository{
					GetByTypeFunc: func(notificationType string) (*internal.RateLimitRule, error) {
						return &internal.RateLimitRule{
							NotificationsLimit: 5,
							IntervalInMinutes:  10,
						}, nil
					},
				}
			},
			cacheRepoFunc: func() *MockRateLimitCacheRepository {
				return &MockRateLimitCacheRepository{
					CountNotificationsWithinIntervalFunc: func(
						notificationType,
						email string,
						intervalInMinutes int,
					) (int, error) {
						return 3, nil
					},
					SetNotificationSentTimestampFunc: func(
						notificationType,
						email,
						timestamp,
						uuid string,
						ttl int64,
					) error {
						return errors.New("cache update error")
					},
				}
			},
			want:    false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rulesRepo := tt.rulesRepoFunc()
			cacheRepo := tt.cacheRepoFunc()

			ucInstance := NewValidateRateLimitUC(rulesRepo, cacheRepo)
			result, err := ucInstance.Handle(notification)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.want, result)
		})
	}
}
