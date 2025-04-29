package controller

import (
	"CMDB/azure"
	"CMDB/repository"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

// APIController 结构体
type APIController struct {
	vmRepo       *repository.VMRepository
	databaseRepo *repository.DatabaseRepository // 添加 DatabaseRepository
	azureService *azure.AzureService
}

// NewAPIController 创建新的API控制器
func NewAPIController(
	vmRepo *repository.VMRepository,
	databaseRepo *repository.DatabaseRepository, // 添加 DatabaseRepository
	azureService *azure.AzureService,
) *APIController {
	return &APIController{
		vmRepo:       vmRepo,
		databaseRepo: databaseRepo, // 初始化 DatabaseRepository
		azureService: azureService,
	}
}

// HandleGetAllVMs 处理获取所有虚拟机的请求
func (c *APIController) HandleGetAllVMs(w http.ResponseWriter, r *http.Request) {
	vms, err := c.vmRepo.ListVMs()
	if err != nil {
		http.Error(w, "Failed to get VMs", http.StatusInternalServerError)
		log.Printf("Error getting VMs: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vms)
}

// HandleGetVMByID 处理根据ID获取虚拟机的请求
func (c *APIController) HandleGetVMByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/vm/")
	if id == "" {
		http.Error(w, "VM ID is required", http.StatusBadRequest)
		return
	}
	vm, err := c.vmRepo.GetVMByID(id)
	if err != nil {
		http.Error(w, "VM not found", http.StatusNotFound)
		log.Printf("Error getting VM by ID %s: %v", id, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vm)
}

// HandleGetAllDatabases 处理获取所有数据库的请求 (新增)
func (c *APIController) HandleGetAllDatabases(w http.ResponseWriter, r *http.Request) {
	databases, err := c.databaseRepo.GetAllDatabases()
	if err != nil {
		http.Error(w, "Failed to get databases", http.StatusInternalServerError)
		log.Printf("Error getting databases: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(databases)
}

// HandleGetDatabaseByID 处理根据ID获取数据库的请求 (新增)
func (c *APIController) HandleGetDatabaseByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/database/")
	if id == "" {
		http.Error(w, "Database ID is required", http.StatusBadRequest)
		return
	}
	database, err := c.databaseRepo.GetDatabaseByResourceID(id)
	if err != nil {
		http.Error(w, "Database not found", http.StatusNotFound)
		log.Printf("Error getting database by ID %s: %v", id, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(database)
}

// HandleGetAllSQLDatabases 处理获取所有SQL数据库的请求
func (c *APIController) HandleGetAllSQLDatabases(w http.ResponseWriter, r *http.Request) {
	databases, err := c.databaseRepo.GetDatabasesByType("SQL Database")
	if err != nil {
		http.Error(w, "获取SQL数据库失败", http.StatusInternalServerError)
		log.Printf("获取SQL数据库错误: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(databases)
}

// HandleGetAllSQLServers 处理获取所有SQL服务器的请求
func (c *APIController) HandleGetAllSQLServers(w http.ResponseWriter, r *http.Request) {
	databases, err := c.databaseRepo.GetDatabasesByType("SQL Server")
	if err != nil {
		http.Error(w, "获取SQL服务器失败", http.StatusInternalServerError)
		log.Printf("获取SQL服务器错误: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(databases)
}

// HandleGetAllMySQLFlexibles 处理获取所有MySQL灵活服务器的请求
func (c *APIController) HandleGetAllMySQLFlexibles(w http.ResponseWriter, r *http.Request) {
	databases, err := c.databaseRepo.GetDatabasesByType("MySQL Flexible Server")
	if err != nil {
		http.Error(w, "获取MySQL灵活服务器失败", http.StatusInternalServerError)
		log.Printf("获取MySQL灵活服务器错误: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(databases)
}

// HandleSyncResources 处理同步资源的请求
func (c *APIController) HandleSyncResources(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 异步执行同步任务
	go func() {
		log.Println("Starting resource synchronization...")
		err := c.azureService.SyncAllResources()
		if err != nil {
			log.Printf("Error during resource synchronization: %v", err)
		} else {
			log.Println("Resource synchronization completed successfully.")
		}
	}()

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Resource synchronization started"))
}

// RegisterRoutes 注册API路由
func (c *APIController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/vms", c.HandleGetAllVMs)
	mux.HandleFunc("/api/vm/", c.HandleGetVMByID)
	mux.HandleFunc("/api/sqldatabase", c.HandleGetAllSQLDatabases)
	mux.HandleFunc("/api/sqlserver", c.HandleGetAllSQLServers)
	mux.HandleFunc("/api/mysqlflexible", c.HandleGetAllMySQLFlexibles)
}
