// model/database.go
package model

import (
	"time"
)

// Database 数据库模型
type Database struct {
	ID             int64             `json:"-"`
	DatabaseID     string            `json:"database_id"`
	ResourceID     string            `json:"resource_id"`
	Name           string            `json:"name"`
	Location       string            `json:"location"`
	Server         string            `json:"server"`
	DBType         string            `json:"db_type"`
	Version        string            `json:"version"`
	Status         string            `json:"status"`
	Owner          string            `json:"owner"`
	SubscriptionID string            `json:"subscription_id"`
	Tags           map[string]string `json:"tags"`
	LastSyncAt     time.Time         `json:"last_sync_at"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
}

// SyncTask 同步任务模型
type SyncTask struct {
	ID          int64     `json:"id"`
	TaskType    string    `json:"task_type"`
	Status      string    `json:"status"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	ItemCount   int       `json:"item_count"`
	ErrorMsg    string    `json:"error_msg"`
	CreatedAt   time.Time `json:"created_at"`
}