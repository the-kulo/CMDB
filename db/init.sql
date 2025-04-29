-- 创建数据库
CREATE DATABASE cmdb
    DEFAULT CHARACTER SET = 'utf8mb4'
    DEFAULT COLLATE = 'utf8mb4_unicode_ci';

-- 使用数据库
USE cmdb;

-- 创建资源表
CREATE TABLE resources (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    resource_id VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    location VARCHAR(255) NOT NULL,
    resource_type VARCHAR(50) NOT NULL,
    owner VARCHAR(255),
    status VARCHAR(50) NOT NULL,
    subscription_id VARCHAR(255) NOT NULL,
    last_sync_at DATETIME NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_resource_id (resource_id),
    INDEX idx_resource_type (resource_type),
    INDEX idx_subscription_id (subscription_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建资源标签表
CREATE TABLE resource_tags (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    resource_id VARCHAR(255) NOT NULL,
    tag_key VARCHAR(255) NOT NULL,
    tag_value VARCHAR(255) NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (resource_id) REFERENCES resources(resource_id) ON DELETE CASCADE,
    UNIQUE KEY uk_resource_tag (resource_id, tag_key)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建虚拟机表
CREATE TABLE vms (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    vm_id VARCHAR(255) NOT NULL UNIQUE,
    resource_id VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    location VARCHAR(255) NOT NULL,
    size VARCHAR(100) NOT NULL,
    status VARCHAR(50) NOT NULL,
    owner VARCHAR(255),
    subscription_id VARCHAR(255) NOT NULL,
    last_sync_at DATETIME NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (resource_id) REFERENCES resources(resource_id),
    INDEX idx_vm_id (vm_id),
    INDEX idx_resource_id (resource_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建虚拟机标签表
CREATE TABLE vm_tags (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    vm_id VARCHAR(255) NOT NULL,
    tag_key VARCHAR(255) NOT NULL,
    tag_value VARCHAR(255) NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (vm_id) REFERENCES vms(vm_id) ON DELETE CASCADE,
    UNIQUE KEY uk_vm_tag (vm_id, tag_key)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建数据库资源表
CREATE TABLE cmdb_databases (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    database_id VARCHAR(255) NOT NULL UNIQUE,
    resource_id VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    location VARCHAR(255) NOT NULL,
    server VARCHAR(255) NOT NULL,
    db_type VARCHAR(50) NOT NULL,
    version VARCHAR(50),
    status VARCHAR(50) NOT NULL,
    owner VARCHAR(255),
    subscription_id VARCHAR(255) NOT NULL,
    last_sync_at DATETIME NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (resource_id) REFERENCES resources(resource_id),
    INDEX idx_database_id (database_id),
    INDEX idx_resource_id (resource_id),
    INDEX idx_db_type (db_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建数据库标签表
CREATE TABLE cmdb_database_tags (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    database_id VARCHAR(255) NOT NULL,
    tag_key VARCHAR(255) NOT NULL,
    tag_value VARCHAR(255) NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (database_id) REFERENCES cmdb_databases(database_id) ON DELETE CASCADE,
    UNIQUE KEY uk_database_tag (database_id, tag_key)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;