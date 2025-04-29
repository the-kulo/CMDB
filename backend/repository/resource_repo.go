// repository/resource_repo.go
package repository

import (
	"CMDB/dao"
	"CMDB/model"
	"log"
)

// ResourceRepository 资源仓库
type ResourceRepository struct {
	resourceDAO *dao.ResourceDAO
}

// NewResourceRepository 创建资源仓库
func NewResourceRepository(resourceDAO *dao.ResourceDAO) *ResourceRepository {
	return &ResourceRepository{resourceDAO: resourceDAO}
}

// SaveResource 保存资源及其标签
func (repo *ResourceRepository) SaveResource(resource *model.Resource) error {
	// 开始事务
	tx, err := repo.resourceDAO.BeginTx()
	if err != nil {
		return err
	}

	// 保存资源基本信息
	err = repo.resourceDAO.UpsertResourceTx(tx, resource)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 保存资源标签
	err = repo.resourceDAO.UpsertResourceTagsTx(tx, resource.ResourceID, resource.Tags)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 提交事务
	return tx.Commit()
}

// BatchSaveResources 批量保存资源
func (repo *ResourceRepository) BatchSaveResources(resources []*model.Resource) error {
	var lastErr error
	successCount := 0

	for _, resource := range resources {
		if err := repo.SaveResource(resource); err != nil {
			log.Printf("保存资源 %s 失败: %v", resource.ResourceID, err)
			lastErr = err
		} else {
			successCount++
		}
	}

	log.Printf("批量保存资源完成: 总数 %d, 成功 %d, 失败 %d", len(resources), successCount, len(resources)-successCount)

	return lastErr
}

// GetResourceByID 根据ID获取资源
func (repo *ResourceRepository) GetResourceByID(resourceID string) (*model.Resource, error) {
	return repo.resourceDAO.GetResourceByID(resourceID)
}

// GetResourcesByType 根据类型获取资源列表
func (repo *ResourceRepository) GetResourcesByType(resourceType string) ([]*model.Resource, error) {
	return repo.resourceDAO.GetResourcesByType(resourceType)
}

// GetAllResources 获取所有资源
func (repo *ResourceRepository) GetAllResources() ([]*model.Resource, error) {
	return repo.resourceDAO.GetAllResources()
}
