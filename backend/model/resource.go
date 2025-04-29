// model/resource.go
package model

import (
	"time"
)

// Resource 资源基本模型
type Resource struct {
	ID             int64             `json:"-"`
	ResourceID     string            `json:"resource_id"`
	Name           string            `json:"name"`
	Location       string            `json:"location"`
	ResourceType   string            `json:"resource_type"`
	Owner          string            `json:"owner"`
	Status         string            `json:"status"`
	SubscriptionID string            `json:"subscription_id"`
	Tags           map[string]string `json:"tags"`
	LastSyncAt     time.Time         `json:"last_sync_at"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
}