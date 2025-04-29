package service

import (
	"CMDB/model"
	"CMDB/repository"
)

// QueryService 资源查询服务
type QueryService struct {
	resourceRepo *repository.ResourceRepository
	vmRepo       *repository.VMRepository
	databaseRepo *repository.DatabaseRepository
}

// NewQueryService 创建新的查询服务
func NewQueryService(
	resourceRepo *repository.ResourceRepository,
	vmRepo *repository.VMRepository,
	databaseRepo *repository.DatabaseRepository,
) *QueryService {
	return &QueryService{
		resourceRepo: resourceRepo,
		vmRepo:       vmRepo,
		databaseRepo: databaseRepo,
	}
}

// GetAllVMs 获取所有虚拟机
func (s *QueryService) GetAllVMs() ([]*model.VM, error) {
	return s.vmRepo.ListVMs()
}

// GetVMByID 根据ID获取虚拟机
func (s *QueryService) GetVMByID(vmID string) (*model.VM, error) {
	return s.vmRepo.GetVMByID(vmID)
}

// GetVMByResourceID 根据资源ID获取虚拟机
func (s *QueryService) GetVMByResourceID(resourceID string) (*model.VM, error) {
	return s.vmRepo.GetVMByResourceID(resourceID)
}

// GetAllDatabases 获取所有数据库
func (s *QueryService) GetAllDatabases() ([]*model.Database, error) {
	return s.databaseRepo.GetAllDatabases()
}

// GetAllResources 获取所有资源
func (s *QueryService) GetAllResources() ([]*model.Resource, error) {
	return s.resourceRepo.GetAllResources()
}