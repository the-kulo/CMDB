// model/vm.go
package model

import (
	"time"
)

// VM 虚拟机模型
type VM struct {
	ID             int64             `json:"-"`
	VMID           string            `json:"vm_id"`
	ResourceID     string            `json:"resource_id"`
	Name           string            `json:"name"`
	Location       string            `json:"location"`
	Type           string            `json:"type"`
	Status         string            `json:"status"`
	Owner          string            `json:"owner"`
	SubscriptionID string            `json:"subscription_id"`
	Tags           map[string]string `json:"tags"`
	LastSyncAt     time.Time         `json:"last_sync_at"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
}
