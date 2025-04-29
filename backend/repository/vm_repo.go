// repository/vm_repo.go
package repository

import (
	"CMDB/dao"
	"CMDB/model"
)

// VMRepository 虚拟机仓库
type VMRepository struct {
	vmDAO *dao.VMDAO
}

// NewVMRepository 创建虚拟机仓库
func NewVMRepository(vmDAO *dao.VMDAO) *VMRepository {
	return &VMRepository{vmDAO: vmDAO}
}

// SaveVirtualMachine 保存虚拟机
func (repo *VMRepository) SaveVirtualMachine(vm *model.VM) error {
	return repo.vmDAO.UpsertVM(vm)
}

// BatchSaveVirtualMachines 批量保存虚拟机
func (repo *VMRepository) BatchSaveVirtualMachines(vms []*model.VM) error {
	for _, vm := range vms {
		if err := repo.SaveVirtualMachine(vm); err != nil {
			return err
		}
	}
	return nil
}

// GetVMByResourceID 根据资源ID获取虚拟机
func (repo *VMRepository) GetVMByResourceID(resourceID string) (*model.VM, error) {
	// 由于DAO中没有直接通过ResourceID获取VM的方法，我们需要获取所有VM然后筛选
	vms, err := repo.ListVMs()
	if err != nil {
		return nil, err
	}
	
	for _, vm := range vms {
		if vm.ResourceID == resourceID {
			return vm, nil
		}
	}
	
	return nil, nil
}

// GetAllVMs 获取所有虚拟机
func (repo *VMRepository) GetAllVMs() ([]*model.VM, error) {
	return repo.vmDAO.ListVMs()
}

// SaveVM 保存虚拟机及其标签
func (repo *VMRepository) SaveVM(vm *model.VM) error {
	// 开始事务
	tx, err := repo.vmDAO.BeginTx()
	if err != nil {
		return err
	}

	// 保存虚拟机基本信息
	err = repo.vmDAO.UpsertVMTx(tx, vm)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 保存虚拟机标签
	err = repo.vmDAO.UpsertVMTagsTx(tx, vm.VMID, vm.Tags)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 提交事务
	return tx.Commit()
}

// BatchSaveVMs 批量保存虚拟机
func (repo *VMRepository) BatchSaveVMs(vms []*model.VM) error {
	for _, vm := range vms {
		if err := repo.SaveVM(vm); err != nil {
			return err
		}
	}
	return nil
}

// GetVMByID 根据ID获取虚拟机
func (repo *VMRepository) GetVMByID(vmID string) (*model.VM, error) {
	return repo.vmDAO.GetVMByID(vmID)
}

// ListVMs 列出所有虚拟机
func (repo *VMRepository) ListVMs() ([]*model.VM, error) {
	return repo.vmDAO.ListVMs()
}
