package azure

import (
	"CMDB/model"
	"CMDB/repository"
	"fmt"
)

// AzureService 封装Azure资源同步服务
type AzureService struct {
	azureHelper  *AzureHelper
	vmRepo       *repository.VMRepository
	databaseRepo *repository.DatabaseRepository
	resourceRepo *repository.ResourceRepository
}

// NewAzureService 创建新的Azure服务
func NewAzureService(
	azureHelper *AzureHelper,
	vmRepo *repository.VMRepository,
	databaseRepo *repository.DatabaseRepository,
	resourceRepo *repository.ResourceRepository,
) *AzureService {
	return &AzureService{
		azureHelper:  azureHelper,
		vmRepo:       vmRepo,
		databaseRepo: databaseRepo,
		resourceRepo: resourceRepo,
	}
}

// SyncVirtualMachines 同步虚拟机资源
func (s *AzureService) SyncVirtualMachines() error {
	// 从Azure获取虚拟机资源
	azureVMs, err := s.azureHelper.GetVirtualMachines()
	if err != nil {
		return err
	}

	// 转换为模型
	var vms []*model.VM
	for _, azureVM := range azureVMs {
		vm := &model.VM{
			VMID:           azureVM.ID, // 确保 azureVM.ID 现在是完整的 ARM ID
			ResourceID:     azureVM.ID, // 确保 azureVM.ID 现在是完整的 ARM ID
			Name:           azureVM.Name,
			Location:       azureVM.Location,
			Type:           azureVM.Type,
			Status:         azureVM.Status,
			Owner:          azureVM.Owner,
			SubscriptionID: s.azureHelper.subscriptionID, // 这个可以保留，或者从 ARM ID 中解析出来
			Tags:           azureVM.Tags,
		}
		vms = append(vms, vm)
	}

	// 保存到数据库
	return s.vmRepo.BatchSaveVMs(vms)
}

// SyncDatabases 同步数据库资源
func (s *AzureService) SyncDatabases() error {
	// 从Azure获取SQL数据库资源
	sqlDatabases, err := s.azureHelper.GetSQLDatabases()
	if err != nil {
		return fmt.Errorf("获取SQL数据库资源失败: %v", err)
	}

	// 从Azure获取MySQL灵活服务器资源
	mysqlServers, err := s.azureHelper.GetMySQLFlexibleServers()
	if err != nil {
		return fmt.Errorf("获取MySQL灵活服务器资源失败: %v", err)
	}

	// 从Azure获取SQL服务器资源
	sqlServers, err := s.azureHelper.GetSQLServers()
	if err != nil {
		return fmt.Errorf("获取SQL服务器资源失败: %v", err)
	}

	// 合并所有数据库资源
	allDatabases := append(sqlDatabases, mysqlServers...)
	allDatabases = append(allDatabases, sqlServers...)

	// 转换为模型
	var databases []*model.Database
	for _, azureDB := range allDatabases {
		database := &model.Database{
			DatabaseID:     azureDB.ID,
			ResourceID:     azureDB.ID,
			Name:           azureDB.Name,
			Location:       azureDB.Location,
			Server:         azureDB.Server,
			DBType:         azureDB.DBType,
			Version:        azureDB.Version,
			Status:         azureDB.Status,
			Owner:          azureDB.Owner,
			SubscriptionID: s.azureHelper.subscriptionID,
			Tags:           azureDB.Tags,
		}
		databases = append(databases, database)
	}

	// 保存到数据库
	return s.databaseRepo.BatchSaveDatabases(databases)
}

// SyncResources 同步通用资源
func (s *AzureService) SyncResources() error {
	// 从Azure获取资源列表
	azureResources, err := s.azureHelper.GetResources()
	if err != nil {
		return fmt.Errorf("获取Azure资源列表失败: %v", err)
	}

	// 转换为模型
	var resources []*model.Resource
	for _, azureResource := range azureResources {
		resource := &model.Resource{
			ResourceID:     azureResource.ID,
			Name:           azureResource.Name,
			Location:       azureResource.Location,
			ResourceType:   azureResource.Type,
			Owner:          azureResource.Owner,
			SubscriptionID: s.azureHelper.subscriptionID,
			Tags:           azureResource.Tags,
		}
		resources = append(resources, resource)
	}

	// 保存到数据库
	return s.resourceRepo.BatchSaveResources(resources)
}

// SyncAllResources 同步所有资源
func (s *AzureService) SyncAllResources() error {
	// 同步通用资源
	if err := s.SyncResources(); err != nil {
		return err
	}

	// 同步虚拟机
	if err := s.SyncVirtualMachines(); err != nil {
		return err
	}

	// 同步数据库
	if err := s.SyncDatabases(); err != nil {
		return err
	}

	return nil
}

// 为了向后兼容，添加全局函数
func NewDefaultAzureService(
	vmRepo *repository.VMRepository,
	databaseRepo *repository.DatabaseRepository,
	resourceRepo *repository.ResourceRepository,
) *AzureService {
	azureHelper := NewAzureHelper()
	return NewAzureService(azureHelper, vmRepo, databaseRepo, resourceRepo)
}
