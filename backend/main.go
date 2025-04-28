package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"CMDB/azure"
)

func main() {
	// 加载环境变量
	godotenv.Load(".env.local")
	err := godotenv.Load()
	if err != nil {
		log.Println("error, cannot load .env file")
	}

	// 初始化Gin路由
	r := gin.Default()

	// 配置CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// 虚拟机资源API路由
	r.GET("/api/vms", func(c *gin.Context) {
		vms, err := azure.GetAzureVirtualMachines()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, vms)
	})

	// SQL数据库资源API路由
	r.GET("/api/sqldatabases", func(c *gin.Context) {
		databases, err := azure.GetAzureSQLDatabases()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, databases)
	})

	// MySQL灵活服务器资源API路由
	r.GET("/api/mysqlflexible", func(c *gin.Context) {
		servers, err := azure.GetAzureMySQLFlexibleServers()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, servers)
	})

	// SQL服务器资源API路由
	r.GET("/api/sqlservers", func(c *gin.Context) {
		servers, err := azure.GetAzureSQLServers()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, servers)
	})

	// 保留原有的数据库API路由（向后兼容）
	r.GET("/api/databases", func(c *gin.Context) {

		databases, err := azure.GetAzureDatabases()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, databases)
	})

	// 启动服务器
	log.Println("启动服务器在 :8080 端口...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("无法启动服务器: %v", err)
	}
}
