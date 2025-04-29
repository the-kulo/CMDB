// dao/database_dao.go
package dao

import (
	"database/sql"
	"time"

	"CMDB/model"
)

// DatabaseDAO 数据库资源数据访问对象
type DatabaseDAO struct {
	db *sql.DB
}

// NewDatabaseDAO 创建新的DatabaseDAO实例
func NewDatabaseDAO(db *sql.DB) *DatabaseDAO {
	return &DatabaseDAO{db: db}
}

// UpsertDatabase 插入或更新数据库信息
func (dao *DatabaseDAO) UpsertDatabase(database *model.Database) error {
	query := `
        INSERT INTO cmdb_databases (database_id, resource_id, name, location, server, db_type, version, status, owner, subscription_id, last_sync_at)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
        ON DUPLICATE KEY UPDATE
            resource_id = VALUES(resource_id),
            name = VALUES(name),
            location = VALUES(location),
            server = VALUES(server),
            db_type = VALUES(db_type),
            version = VALUES(version),
            status = VALUES(status),
            owner = VALUES(owner),
            subscription_id = VALUES(subscription_id),
            last_sync_at = VALUES(last_sync_at)
    `

	now := time.Now()
	_, err := dao.db.Exec(
		query,
		database.DatabaseID,
		database.ResourceID,
		database.Name,
		database.Location,
		database.Server,
		database.DBType,
		database.Version,
		database.Status,
		database.Owner,
		database.SubscriptionID,
		now,
	)

	return err
}

// UpsertDatabaseTags 更新数据库标签
// UpsertDatabaseTags 更新数据库标签
func (dao *DatabaseDAO) UpsertDatabaseTags(databaseID string, tags map[string]string) error {
    // 先删除该数据库的所有标签
    _, err := dao.db.Exec("DELETE FROM cmdb_database_tags WHERE database_id = ?", databaseID)
    if err != nil {
        return err
    }

    // 如果没有标签，直接返回
    if len(tags) == 0 {
        return nil
    }

    // 批量插入新标签
    stmt, err := dao.db.Prepare("INSERT INTO cmdb_database_tags (database_id, tag_key, tag_value) VALUES (?, ?, ?)")
    if err != nil {
        return err
    }
    defer stmt.Close()

    for key, value := range tags {
        _, err = stmt.Exec(databaseID, key, value)
        if err != nil {
            return err
        }
    }

    return nil
}

// GetDatabaseByID 根据ID获取数据库信息
func (dao *DatabaseDAO) GetDatabaseByID(databaseID string) (*model.Database, error) {
	query := `
        SELECT id, database_id, resource_id, name, location, server, db_type, version, status, owner, subscription_id, last_sync_at, created_at, updated_at
        FROM cmdb_databases
        WHERE database_id = ?
    `
	
	database := &model.Database{}
	err := dao.db.QueryRow(query, databaseID).Scan(
		&database.ID,
		&database.DatabaseID,
		&database.ResourceID,
		&database.Name,
		&database.Location,
		&database.Server,
		&database.DBType,
		&database.Version,
		&database.Status,
		&database.Owner,
		&database.SubscriptionID,
		&database.LastSyncAt,
		&database.CreatedAt,
		&database.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	
	// 获取标签
	database.Tags = make(map[string]string)
	rows, err := dao.db.Query("SELECT tag_key, tag_value FROM cmdb_database_tags WHERE database_id = ?", databaseID)
	if err != nil {
		return database, err
	}
	defer rows.Close()
	
	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			return database, err
		}
		database.Tags[key] = value
	}
	
	return database, nil
}

// ListDatabases 列出所有数据库
func (dao *DatabaseDAO) ListDatabases() ([]*model.Database, error) {
	query := `
        SELECT id, database_id, resource_id, name, location, server, db_type, version, status, owner, subscription_id, last_sync_at, created_at, updated_at
        FROM cmdb_databases
        ORDER BY name
    `
	
	rows, err := dao.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var databases []*model.Database
	for rows.Next() {
		database := &model.Database{Tags: make(map[string]string)}
		err := rows.Scan(
			&database.ID,
			&database.DatabaseID,
			&database.ResourceID,
			&database.Name,
			&database.Location,
			&database.Server,
			&database.DBType,
			&database.Version,
			&database.Status,
			&database.Owner,
			&database.SubscriptionID,
			&database.LastSyncAt,
			&database.CreatedAt,
			&database.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		databases = append(databases, database)
	}
	
	// 获取所有数据库的标签
	for _, database := range databases {
		tagRows, err := dao.db.Query("SELECT tag_key, tag_value FROM cmdb_database_tags WHERE database_id = ?", database.DatabaseID)
		if err != nil {
			return databases, err
		}
		
		for tagRows.Next() {
			var key, value string
			if err := tagRows.Scan(&key, &value); err != nil {
				tagRows.Close()
				return databases, err
			}
			database.Tags[key] = value
		}
		tagRows.Close()
	}
	
	return databases, nil
}

// BeginTx 开始事务
func (dao *DatabaseDAO) BeginTx() (*sql.Tx, error) {
	return dao.db.Begin()
}

// UpsertDatabaseTx 在事务中插入或更新数据库信息
func (dao *DatabaseDAO) UpsertDatabaseTx(tx *sql.Tx, database *model.Database) error {
	query := `
        INSERT INTO cmdb_databases (database_id, resource_id, name, location, server, db_type, version, status, owner, subscription_id, last_sync_at)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
        ON DUPLICATE KEY UPDATE
            resource_id = VALUES(resource_id),
            name = VALUES(name),
            location = VALUES(location),
            server = VALUES(server),
            db_type = VALUES(db_type),
            version = VALUES(version),
            status = VALUES(status),
            owner = VALUES(owner),
            subscription_id = VALUES(subscription_id),
            last_sync_at = VALUES(last_sync_at)
    `

	now := time.Now()
	_, err := tx.Exec(
		query,
		database.DatabaseID,
		database.ResourceID,
		database.Name,
		database.Location,
		database.Server,
		database.DBType,
		database.Version,
		database.Status,
		database.Owner,
		database.SubscriptionID,
		now,
	)

	return err
}

// UpsertDatabaseTagsTx 在事务中更新数据库标签
func (dao *DatabaseDAO) UpsertDatabaseTagsTx(tx *sql.Tx, databaseID string, tags map[string]string) error {
	// 先删除该数据库的所有标签
	_, err := tx.Exec("DELETE FROM cmdb_database_tags WHERE database_id = ?", databaseID)
	if err != nil {
		return err
	}

	// 如果没有标签，直接返回
	if len(tags) == 0 {
		return nil
	}

	// 批量插入新标签
	stmt, err := tx.Prepare("INSERT INTO cmdb_database_tags (database_id, tag_key, tag_value) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for key, value := range tags {
		_, err = stmt.Exec(databaseID, key, value)
		if err != nil {
			return err
		}
	}

	return nil
}