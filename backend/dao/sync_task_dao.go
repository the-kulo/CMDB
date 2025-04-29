// dao/sync_task_dao.go
package dao

import (
	"database/sql"
	"time"

	"CMDB/model"
)

// SyncTaskDAO 同步任务数据访问对象
type SyncTaskDAO struct {
	db *sql.DB
}

// NewSyncTaskDAO 创建新的SyncTaskDAO实例
func NewSyncTaskDAO(db *sql.DB) *SyncTaskDAO {
	return &SyncTaskDAO{db: db}
}

// CreateSyncTask 创建同步任务记录
func (dao *SyncTaskDAO) CreateSyncTask(taskType string) (int64, error) {
	query := `
        INSERT INTO sync_tasks (task_type, status, start_time, created_at)
        VALUES (?, ?, ?, ?)
    `
	
	now := time.Now()
	result, err := dao.db.Exec(query, taskType, "RUNNING", now, now)
	if err != nil {
		return 0, err
	}
	
	return result.LastInsertId()
}

// UpdateSyncTaskStatus 更新同步任务状态
func (dao *SyncTaskDAO) UpdateSyncTaskStatus(taskID int64, status string, itemCount int, errorMsg string) error {
	query := `
        UPDATE sync_tasks
        SET status = ?, end_time = ?, item_count = ?, error_msg = ?
        WHERE id = ?
    `
	
	_, err := dao.db.Exec(query, status, time.Now(), itemCount, errorMsg, taskID)
	return err
}

// GetSyncTaskByID 根据ID获取同步任务
func (dao *SyncTaskDAO) GetSyncTaskByID(taskID int64) (*model.SyncTask, error) {
	query := `
        SELECT id, task_type, status, start_time, end_time, item_count, error_msg, created_at
        FROM sync_tasks
        WHERE id = ?
    `
	
	task := &model.SyncTask{}
	err := dao.db.QueryRow(query, taskID).Scan(
		&task.ID,
		&task.TaskType,
		&task.Status,
		&task.StartTime,
		&task.EndTime,
		&task.ItemCount,
		&task.ErrorMsg,
		&task.CreatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	
	return task, nil
}

// ListSyncTasks 列出同步任务
func (dao *SyncTaskDAO) ListSyncTasks(limit int) ([]*model.SyncTask, error) {
	query := `
        SELECT id, task_type, status, start_time, end_time, item_count, error_msg, created_at
        FROM sync_tasks
        ORDER BY id DESC
        LIMIT ?
    `
	
	rows, err := dao.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var tasks []*model.SyncTask
	for rows.Next() {
		task := &model.SyncTask{}
		err := rows.Scan(
			&task.ID,
			&task.TaskType,
			&task.Status,
			&task.StartTime,
			&task.EndTime,
			&task.ItemCount,
			&task.ErrorMsg,
			&task.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	
	return tasks, nil
}