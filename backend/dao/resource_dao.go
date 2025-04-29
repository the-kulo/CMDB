package dao

import (
	"database/sql"
	"time"

	"CMDB/model"
)

type ResourceDAO struct {
	db *sql.DB
}

func NewResourceDAO(db *sql.DB) *ResourceDAO {
	return &ResourceDAO{db: db}
}

// BeginTx 开始一个事务
func (dao *ResourceDAO) BeginTx() (*sql.Tx, error) {
	return dao.db.Begin()
}

// UpsertResource 使用 MySQL 的 ON DUPLICATE KEY UPDATE 实现 Upsert
func (dao *ResourceDAO) UpsertResource(resource *model.Resource) error {
	query := `
        INSERT INTO resources (resource_id, name, location, resource_type, owner, status, subscription_id, last_sync_at)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)
        ON DUPLICATE KEY UPDATE
            name = VALUES(name),
            location = VALUES(location),
            resource_type = VALUES(resource_type),
            owner = VALUES(owner),
            status = VALUES(status),
            subscription_id = VALUES(subscription_id),
            last_sync_at = VALUES(last_sync_at)
    `

	now := time.Now()
	_, err := dao.db.Exec(
		query,
		resource.ResourceID,
		resource.Name,
		resource.Location,
		resource.ResourceType,
		resource.Owner,
		resource.Status,
		resource.SubscriptionID,
		now,
	)

	return err
}

// UpsertResourceTx 在事务中执行资源的 Upsert 操作
func (dao *ResourceDAO) UpsertResourceTx(tx *sql.Tx, resource *model.Resource) error {
	query := `
        INSERT INTO resources (resource_id, name, location, resource_type, owner, status, subscription_id, last_sync_at)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)
        ON DUPLICATE KEY UPDATE
            name = VALUES(name),
            location = VALUES(location),
            resource_type = VALUES(resource_type),
            owner = VALUES(owner),
            status = VALUES(status),
            subscription_id = VALUES(subscription_id),
            last_sync_at = VALUES(last_sync_at)
    `

	now := time.Now()
	_, err := tx.Exec(
		query,
		resource.ResourceID,
		resource.Name,
		resource.Location,
		resource.ResourceType,
		resource.Owner,
		resource.Status,
		resource.SubscriptionID,
		now,
	)

	return err
}

// UpsertResourceTags 批量更新资源标签
func (dao *ResourceDAO) UpsertResourceTags(resourceID string, tags map[string]string) error {
	// 先删除该资源的所有标签
	_, err := dao.db.Exec("DELETE FROM resource_tags WHERE resource_id = ?", resourceID)
	if err != nil {
		return err
	}

	// 如果没有标签，直接返回
	if len(tags) == 0 {
		return nil
	}

	// 批量插入新标签
	stmt, err := dao.db.Prepare("INSERT INTO resource_tags (resource_id, tag_key, tag_value) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for key, value := range tags {
		_, err = stmt.Exec(resourceID, key, value)
		if err != nil {
			return err
		}
	}

	return nil
}

// UpsertResourceTagsTx 在事务中批量更新资源标签
func (dao *ResourceDAO) UpsertResourceTagsTx(tx *sql.Tx, resourceID string, tags map[string]string) error {
	// 先删除该资源的所有标签
	_, err := tx.Exec("DELETE FROM resource_tags WHERE resource_id = ?", resourceID)
	if err != nil {
		return err
	}

	// 如果没有标签，直接返回
	if len(tags) == 0 {
		return nil
	}

	// 批量插入新标签
	for key, value := range tags {
		_, err = tx.Exec("INSERT INTO resource_tags (resource_id, tag_key, tag_value) VALUES (?, ?, ?)",
			resourceID, key, value)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetResourceByID 根据ID获取资源
func (dao *ResourceDAO) GetResourceByID(resourceID string) (*model.Resource, error) {
	query := `
        SELECT r.resource_id, r.name, r.location, r.resource_type, r.owner, r.status, r.subscription_id, r.last_sync_at, r.created_at, r.updated_at
        FROM resources r
        WHERE r.resource_id = ?
    `

	var resource model.Resource
	err := dao.db.QueryRow(query, resourceID).Scan(
		&resource.ResourceID,
		&resource.Name,
		&resource.Location,
		&resource.ResourceType,
		&resource.Owner,
		&resource.Status,
		&resource.SubscriptionID,
		&resource.LastSyncAt,
		&resource.CreatedAt,
		&resource.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// 获取资源标签
	resource.Tags, err = dao.getResourceTags(resourceID)
	if err != nil {
		return nil, err
	}

	return &resource, nil
}

// GetResourcesByType 根据类型获取资源列表
func (dao *ResourceDAO) GetResourcesByType(resourceType string) ([]*model.Resource, error) {
	query := `
        SELECT r.resource_id, r.name, r.location, r.resource_type, r.owner, r.status, r.subscription_id, r.last_sync_at, r.created_at, r.updated_at
        FROM resources r
        WHERE r.resource_type = ?
    `

	rows, err := dao.db.Query(query, resourceType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resources []*model.Resource
	for rows.Next() {
		var resource model.Resource
		err := rows.Scan(
			&resource.ResourceID,
			&resource.Name,
			&resource.Location,
			&resource.ResourceType,
			&resource.Owner,
			&resource.Status,
			&resource.SubscriptionID,
			&resource.LastSyncAt,
			&resource.CreatedAt,
			&resource.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// 获取资源标签
		resource.Tags, err = dao.getResourceTags(resource.ResourceID)
		if err != nil {
			return nil, err
		}

		resources = append(resources, &resource)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return resources, nil
}

// GetAllResources 获取所有资源
func (dao *ResourceDAO) GetAllResources() ([]*model.Resource, error) {
	query := `
        SELECT r.resource_id, r.name, r.location, r.resource_type, r.owner, r.status, r.subscription_id, r.last_sync_at, r.created_at, r.updated_at
        FROM resources r
    `

	rows, err := dao.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resources []*model.Resource
	for rows.Next() {
		var resource model.Resource
		err := rows.Scan(
			&resource.ResourceID,
			&resource.Name,
			&resource.Location,
			&resource.ResourceType,
			&resource.Owner,
			&resource.Status,
			&resource.SubscriptionID,
			&resource.LastSyncAt,
			&resource.CreatedAt,
			&resource.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// 获取资源标签
		resource.Tags, err = dao.getResourceTags(resource.ResourceID)
		if err != nil {
			return nil, err
		}

		resources = append(resources, &resource)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return resources, nil
}

// getResourceTags 获取资源的所有标签
func (dao *ResourceDAO) getResourceTags(resourceID string) (map[string]string, error) {
	query := `
        SELECT tag_key, tag_value
        FROM resource_tags
        WHERE resource_id = ?
    `

	rows, err := dao.db.Query(query, resourceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tags := make(map[string]string)
	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			return nil, err
		}
		tags[key] = value
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tags, nil
}

// GetResourcesByLocation 根据位置获取资源
func (dao *ResourceDAO) GetResourcesByLocation(location string) ([]*model.Resource, error) {
	query := `
        SELECT r.resource_id, r.name, r.location, r.resource_type, r.owner, r.status, r.subscription_id, r.last_sync_at, r.created_at, r.updated_at
        FROM resources r
        WHERE r.location = ?
    `

	rows, err := dao.db.Query(query, location)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resources []*model.Resource
	for rows.Next() {
		var resource model.Resource
		err := rows.Scan(
			&resource.ResourceID,
			&resource.Name,
			&resource.Location,
			&resource.ResourceType,
			&resource.Owner,
			&resource.Status,
			&resource.SubscriptionID,
			&resource.LastSyncAt,
			&resource.CreatedAt,
			&resource.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// 获取资源标签
		resource.Tags, err = dao.getResourceTags(resource.ResourceID)
		if err != nil {
			return nil, err
		}

		resources = append(resources, &resource)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return resources, nil
}

// GetResourcesByTag 根据标签获取资源
func (dao *ResourceDAO) GetResourcesByTag(key string, value string) ([]*model.Resource, error) {
	query := `
        SELECT r.resource_id, r.name, r.location, r.resource_type, r.owner, r.status, r.subscription_id, r.last_sync_at, r.created_at, r.updated_at
        FROM resources r
        JOIN resource_tags t ON r.resource_id = t.resource_id
        WHERE t.tag_key = ? AND t.tag_value = ?
    `

	rows, err := dao.db.Query(query, key, value)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resources []*model.Resource
	for rows.Next() {
		var resource model.Resource
		err := rows.Scan(
			&resource.ResourceID,
			&resource.Name,
			&resource.Location,
			&resource.ResourceType,
			&resource.Owner,
			&resource.Status,
			&resource.SubscriptionID,
			&resource.LastSyncAt,
			&resource.CreatedAt,
			&resource.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// 获取资源标签
		resource.Tags, err = dao.getResourceTags(resource.ResourceID)
		if err != nil {
			return nil, err
		}

		resources = append(resources, &resource)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return resources, nil
}
