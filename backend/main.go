package main

import (
	"CMDB/azure"
	"CMDB/config"
	"CMDB/controller"
	"CMDB/dao"
	"CMDB/repository"
	"CMDB/scheduler"
	"CMDB/service"
	"database/sql"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// CORS中间件，用于处理跨域请求
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 设置CORS头
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// 处理预检请求
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// 继续处理请求
		next.ServeHTTP(w, r)
	})
}

func main() {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 连接数据库
	// 连接数据库
	db, err := sql.Open("mysql", cfg.DatabaseDSN)
	if err != nil {
	    log.Fatalf("连接数据库失败: %v", err)
	}
	defer db.Close()
	
	// 配置连接池
	db.SetMaxOpenConns(25)  // 最大连接数
	db.SetMaxIdleConns(10)  // 最大空闲连接数
	db.SetConnMaxLifetime(5 * time.Minute)  // 连接最大生命周期
	
	// 验证数据库连接
	if err := db.Ping(); err != nil {
	    log.Fatalf("数据库连接测试失败: %v", err)
	}
	// 初始化DAO
	vmDAO := dao.NewVMDAO(db)
	databaseDAO := dao.NewDatabaseDAO(db)
	resourceDAO := dao.NewResourceDAO(db)

	// 初始化Repository
	vmRepo := repository.NewVMRepository(vmDAO)
	databaseRepo := repository.NewDatabaseRepository(databaseDAO)
	resourceRepo := repository.NewResourceRepository(resourceDAO)

	// 初始化Azure Helper
	azureHelper := azure.NewAzureHelper()
	if err := azureHelper.Initialize(); err != nil {
		log.Fatalf("初始化Azure Helper失败: %v", err)
	}

	// 初始化Azure Service
	azureService := azure.NewAzureService(azureHelper, vmRepo, databaseRepo, resourceRepo)

	// 初始化Service
	syncService := service.NewSyncService(azureService, resourceRepo, vmRepo, databaseRepo)
	// 删除未使用的queryService变量

	// 初始化Controller
	apiController := controller.NewAPIController(vmRepo, databaseRepo, azureService)

	// 注册路由
	mux := http.NewServeMux()
	apiController.RegisterRoutes(mux)

	// 初始化定时任务
	cronScheduler := scheduler.NewCronScheduler(syncService, 6*time.Hour)
	cronScheduler.Start()
	defer cronScheduler.Stop()

	// 使用CORS中间件包装HTTP处理器
	corsHandler := corsMiddleware(mux)

	// 启动HTTP服务器
	log.Printf("HTTP服务器启动在 %s", cfg.ServerAddress)
	if err := http.ListenAndServe(cfg.ServerAddress, corsHandler); err != nil {
		log.Fatalf("HTTP服务器启动失败: %v", err)
	}
}
