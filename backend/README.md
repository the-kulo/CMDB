# CMDB

CMDB/
├── backend/
│   ├── main.go                 # 主程序入口
│   ├── config/                 # 配置文件
│   │   └── config.go
│   ├── model/                  # 数据模型
│   │   ├── resource.go
│   │   ├── vm.go
│   │   └── database.go
│   ├── dao/                    # 数据访问层
│   │   ├── resource_dao.go
│   │   ├── vm_dao.go
│   │   └── database_dao.go
│   ├── repository/             # 仓库层
│   │   ├── resource_repo.go
│   │   ├── vm_repo.go
│   │   └── database_repo.go
│   ├── service/                # 业务逻辑层
│   │   ├── sync_service.go
│   │   └── query_service.go
│   ├── controller/             # 控制器层
│   │   └── api_controller.go
│   ├── scheduler/              # 定时任务
│   │   └── cron_scheduler.go
│   └── azure/                  # Azure API 封装
│       └── azure.go
└── README.md
