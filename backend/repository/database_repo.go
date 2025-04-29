package repository

import (
	"CMDB/dao"
	"CMDB/model"
)

// DatabaseRepository 数据库资源仓库
type DatabaseRepository struct {
	databaseDAO *dao.DatabaseDAO
}

// NewDatabaseRepository 创建数据库资源仓库
func NewDatabaseRepository(databaseDAO *dao.DatabaseDAO) *DatabaseRepository {
	return &DatabaseRepository{databaseDAO: databaseDAO}
}

// SaveDatabaseResource 保存数据库资源
func (repo *DatabaseRepository) SaveDatabaseResource(database *model.Database) error {
	return repo.databaseDAO.UpsertDatabase(database)
}

// BatchSaveDatabaseResources 批量保存数据库资源
func (repo *DatabaseRepository) BatchSaveDatabaseResources(databases []*model.Database) error {
	for _, database := range databases {
		if err := repo.SaveDatabaseResource(database); err != nil {
			return err
		}
	}
	return nil
}

// GetDatabaseByResourceID 根据资源ID获取数据库资源
func (repo *DatabaseRepository) GetDatabaseByResourceID(resourceID string) (*model.Database, error) {
	return repo.databaseDAO.GetDatabaseByID(resourceID)
}

// GetDatabasesByType 根据数据库类型获取数据库资源
func (repo *DatabaseRepository) GetDatabasesByType(dbType string) ([]*model.Database, error) {
	// 获取所有数据库资源
	allDatabases, err := repo.databaseDAO.ListDatabases()
	if err != nil {
		return nil, err
	}

	// 筛选指定类型的数据库资源
	var result []*model.Database
	for _, db := range allDatabases {
		if db.DBType == dbType {
			result = append(result, db)
		}
	}

	return result, nil
}

// GetAllDatabases 获取所有数据库资源
func (repo *DatabaseRepository) GetAllDatabases() ([]*model.Database, error) {
	return repo.databaseDAO.ListDatabases()
}

// BatchSaveDatabases 批量保存数据库资源
func (repo *DatabaseRepository) BatchSaveDatabases(databases []*model.Database) error {
	return repo.BatchSaveDatabaseResources(databases)
}
