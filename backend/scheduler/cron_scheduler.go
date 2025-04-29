package scheduler

import (
	"CMDB/service"
	"log"
	"time"
)

// CronScheduler 定时任务调度器
type CronScheduler struct {
	syncService *service.SyncService
	interval    time.Duration
	stopChan    chan struct{}
}

// NewCronScheduler 创建新的定时任务调度器
func NewCronScheduler(syncService *service.SyncService, interval time.Duration) *CronScheduler {
	return &CronScheduler{
		syncService: syncService,
		interval:    interval,
		stopChan:    make(chan struct{}),
	}
}

// Start 启动定时任务
func (s *CronScheduler) Start() {
	go func() {
		ticker := time.NewTicker(s.interval)
		defer ticker.Stop()

		// 启动时立即执行一次同步
		if err := s.syncService.SyncAllResources(); err != nil {
			log.Printf("资源同步失败: %v", err)
		} else {
			log.Println("资源同步成功")
		}

		for {
			select {
			case <-ticker.C:
				// 定时执行同步
				if err := s.syncService.SyncAllResources(); err != nil {
					log.Printf("资源同步失败: %v", err)
				} else {
					log.Println("资源同步成功")
				}
			case <-s.stopChan:
				return
			}
		}
	}()
}

// Stop 停止定时任务
func (s *CronScheduler) Stop() {
	close(s.stopChan)
}