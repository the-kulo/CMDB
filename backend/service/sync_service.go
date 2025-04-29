package service

import (
	"CMDB/azure"
	"CMDB/repository"
	"time"
)

// SyncService 资源同步服务
type SyncService struct {
	azureService     *azure.AzureService
	resourceRepo     *repository.ResourceRepository
	vmRepo           *repository.VMRepository
	databaseRepo     *repository.DatabaseRepository
}

// NewSyncService 创建新的同步服务
func NewSyncService(
	azureService *azure.AzureService,
	resourceRepo *repository.ResourceRepository,
	vmRepo *repository.VMRepository,
	databaseRepo *repository.DatabaseRepository,
) *SyncService {
	return &SyncService{
		azureService: azureService,
		resourceRepo: resourceRepo,
		vmRepo:       vmRepo,
		databaseRepo: databaseRepo,
	}
}

// SyncAllResources 同步所有资源
func (s *SyncService) SyncAllResources() error {
	// 使用Azure服务同步所有资源
	return s.azureService.SyncAllResources()
}

// SyncVirtualMachines 同步虚拟机资源
func (s *SyncService) SyncVirtualMachines() error {
	// 使用Azure服务同步虚拟机资源
	return s.azureService.SyncVirtualMachines()
}

// GetLastSyncTime 获取最后同步时间
func (s *SyncService) GetLastSyncTime() (time.Time, error) {
	// 这里需要实现获取最后同步时间的逻辑
	// 可能需要在数据库中添加一个表来记录同步状态
	return time.Now(), nil
}