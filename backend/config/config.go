// config/config.go
package config

import (
	"fmt"
	"os"
	"github.com/joho/godotenv"
)

// Config 应用配置结构
type Config struct {
	DatabaseDSN   string
	ServerPort    string
	ServerAddress string
	AzureConfig   AzureConfig
}

// AzureConfig Azure配置
type AzureConfig struct {
	ClientID       string
	TenantID       string
	ClientSecret   string
	SubscriptionID string
}

// LoadConfig 从环境变量加载配置
func LoadConfig() (*Config, error) {
	// 尝试加载.env文件
	_ = godotenv.Load()
	
	// 数据库配置
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	
	// 使用默认值（如果环境变量未设置）
	if dbUser == "" {
		dbUser = "root"
	}
	if dbPass == "" {
		dbPass = ""
	}
	if dbHost == "" {
		dbHost = "localhost"
	}
	if dbPort == "" {
		dbPort = "3306"
	}
	if dbName == "" {
		dbName = "cmdb"
	}
	
	// 构建DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", 
		dbUser, dbPass, dbHost, dbPort, dbName)
	
	// 服务器配置
	serverPort := getEnvOrDefault("SERVER_PORT", "8080")
	serverAddress := fmt.Sprintf(":%s", serverPort)
	
	// Azure配置
	azureConfig := AzureConfig{
		ClientID:       os.Getenv("CLIENT_ID"),
		TenantID:       os.Getenv("TENANT_ID"),
		ClientSecret:   os.Getenv("CLIENT_SECRET"),
		SubscriptionID: os.Getenv("SUBSCRIPTION_ID"),
	}
	
	return &Config{
		DatabaseDSN:   dsn,
		ServerPort:    serverPort,
		ServerAddress: serverAddress,
		AzureConfig:   azureConfig,
	}, nil
}

// getEnvOrDefault 获取环境变量，如果不存在则返回默认值
func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}