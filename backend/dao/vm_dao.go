// dao/vm_dao.go
package dao

import (
	"database/sql"
	"time"

	"CMDB/model"
)

// VMDAO 虚拟机数据访问对象
type VMDAO struct {
	db *sql.DB
}

// NewVMDAO 创建新的VMDAO实例
func NewVMDAO(db *sql.DB) *VMDAO {
	return &VMDAO{db: db}
}

// UpsertVM 插入或更新虚拟机信息
func (dao *VMDAO) UpsertVM(vm *model.VM) error {
	query := `
        INSERT INTO vms (vm_id, resource_id, name, location, type, status, owner, subscription_id, last_sync_at)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
        ON DUPLICATE KEY UPDATE
            resource_id = VALUES(resource_id),
            name = VALUES(name),
            location = VALUES(location),
            type = VALUES(type),
            status = VALUES(status),
            owner = VALUES(owner),
            subscription_id = VALUES(subscription_id),
            last_sync_at = VALUES(last_sync_at)
    `

	now := time.Now()
	_, err := dao.db.Exec(
		query,
		vm.VMID,
		vm.ResourceID,
		vm.Name,
		vm.Location,
		vm.Type,
		vm.Status,
		vm.Owner,
		vm.SubscriptionID,
		now,
	)

	return err
}

// UpsertVMTags 更新虚拟机标签
func (dao *VMDAO) UpsertVMTags(vmID string, tags map[string]string) error {
	// 先删除该虚拟机的所有标签
	_, err := dao.db.Exec("DELETE FROM vm_tags WHERE vm_id = ?", vmID)
	if err != nil {
		return err
	}

	// 如果没有标签，直接返回
	if len(tags) == 0 {
		return nil
	}

	// 批量插入新标签
	stmt, err := dao.db.Prepare("INSERT INTO vm_tags (vm_id, tag_key, tag_value) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for key, value := range tags {
		_, err = stmt.Exec(vmID, key, value)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetVMByID 根据ID获取虚拟机信息
func (dao *VMDAO) GetVMByID(vmID string) (*model.VM, error) {
	query := `
        SELECT id, vm_id, resource_id, name, location, type, status, owner, subscription_id, last_sync_at, created_at, updated_at
        FROM vms
        WHERE vm_id = ?
    `

	vm := &model.VM{}
	err := dao.db.QueryRow(query, vmID).Scan(
		&vm.ID,
		&vm.VMID,
		&vm.ResourceID,
		&vm.Name,
		&vm.Location,
		&vm.Type,
		&vm.Status,
		&vm.Owner,
		&vm.SubscriptionID,
		&vm.LastSyncAt,
		&vm.CreatedAt,
		&vm.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// 获取标签
	vm.Tags = make(map[string]string)
	rows, err := dao.db.Query("SELECT tag_key, tag_value FROM vm_tags WHERE vm_id = ?", vmID)
	if err != nil {
		return vm, err
	}
	defer rows.Close()

	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			return vm, err
		}
		vm.Tags[key] = value
	}

	return vm, nil
}

// ListVMs 列出所有虚拟机
func (dao *VMDAO) ListVMs() ([]*model.VM, error) {
	query := `
        SELECT id, vm_id, resource_id, name, location, type, status, owner, subscription_id, last_sync_at, created_at, updated_at
        FROM vms
        ORDER BY name
    `

	rows, err := dao.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vms []*model.VM
	for rows.Next() {
		vm := &model.VM{Tags: make(map[string]string)}
		err := rows.Scan(
			&vm.ID,
			&vm.VMID,
			&vm.ResourceID,
			&vm.Name,
			&vm.Location,
			&vm.Type,
			&vm.Status,
			&vm.Owner,
			&vm.SubscriptionID,
			&vm.LastSyncAt,
			&vm.CreatedAt,
			&vm.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		vms = append(vms, vm)
	}

	// 获取所有VM的标签
	for _, vm := range vms {
		tagRows, err := dao.db.Query("SELECT tag_key, tag_value FROM vm_tags WHERE vm_id = ?", vm.VMID)
		if err != nil {
			return vms, err
		}

		for tagRows.Next() {
			var key, value string
			if err := tagRows.Scan(&key, &value); err != nil {
				tagRows.Close()
				return vms, err
			}
			vm.Tags[key] = value
		}
		tagRows.Close()
	}

	return vms, nil
}

// BeginTx 开始事务
func (dao *VMDAO) BeginTx() (*sql.Tx, error) {
	return dao.db.Begin()
}

// UpsertVMTx 在事务中插入或更新虚拟机
func (dao *VMDAO) UpsertVMTx(tx *sql.Tx, vm *model.VM) error {
	query := `
        INSERT INTO vms (vm_id, resource_id, name, location, type, status, owner, subscription_id, last_sync_at)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
        ON DUPLICATE KEY UPDATE
            resource_id = VALUES(resource_id),
            name = VALUES(name),
            location = VALUES(location),
            type = VALUES(type),
            status = VALUES(status),
            owner = VALUES(owner),
            subscription_id = VALUES(subscription_id),
            last_sync_at = VALUES(last_sync_at)
    `

	now := time.Now()
	_, err := tx.Exec(
		query,
		vm.VMID,
		vm.ResourceID,
		vm.Name,
		vm.Location,
		vm.Type,
		vm.Status,
		vm.Owner,
		vm.SubscriptionID,
		now,
	)

	return err
}

// UpsertVMTagsTx 在事务中更新虚拟机标签
func (dao *VMDAO) UpsertVMTagsTx(tx *sql.Tx, vmID string, tags map[string]string) error {
	// 先删除该虚拟机的所有标签
	_, err := tx.Exec("DELETE FROM vm_tags WHERE vm_id = ?", vmID)
	if err != nil {
		return err
	}

	// 如果没有标签，直接返回
	if len(tags) == 0 {
		return nil
	}

	// 批量插入新标签
	stmt, err := tx.Prepare("INSERT INTO vm_tags (vm_id, tag_key, tag_value) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for key, value := range tags {
		_, err = stmt.Exec(vmID, key, value)
		if err != nil {
			return err
		}
	}

	return nil
}
