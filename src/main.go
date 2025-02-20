package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"smsAlert/handler"
	"syscall"

	"smsAlert/logger"
	"smsAlert/models"
)

//// 健康状态
//var (
//	isReady   atomic.Value // 就绪状态
//	isStarted atomic.Value // 启动状态
//)
//
//// 初始化状态
//func init() {
//	isReady.Store(false)
//	isStarted.Store(false)
//}

func main() {
	// 日志中打印当前目录，方便debug载入配置yaml出错
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory:", err)
	} else {
		fmt.Println("Current working directory:", dir)
	}

	// 命名行传递参数
	port := flag.String("addr", ":5000", "address to listen for webhook")
	capacity := flag.Int("cap", 64, "capacity of alertStore")
	configPath := flag.String("config.file", "/app/settings.yml", "Path to the configuration file")
	flag.Parse()

	// 创建存储报警的结构体
	// 变量类型为定义在module包中的AlertStore
	s := &models.AlertStore{
		Capacity: *capacity,
	}

	// 启动服务器，接收告警
	go Webhook(port, s, configPath)
	//select {}

	// 监听中断信号，阻塞主goroutine，收到终止信号后解除阻塞并打印日志
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	<-signalChan
	fmt.Println("Shutting down webhook...")
	logger.Log.Info(fmt.Sprintf("Shutting down webhook..."))
}

func Webhook(port *string, s *models.AlertStore, configPath *string) {
	// 初始化rolling-logger
	logger.InitLogger()

	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		// 同步处理 PostHandler，确保接收AM发送的告警json被存储
		s.PostHandler(w, r)

		// 异步处理 查询手机号、发送短信
		go func() {
			phone := handler.Query(s, configPath)
			if phone != "" {
				if err := handler.SendSMS(phone, s); err != nil {
					logger.Log.Errorf("Failed to send SMS: %v", err)
				} else {
					s.Lock()
					s.Alerts = nil // 清空 alertStore 存储的告警
					s.Unlock()
					logger.Log.Info("SendSMS successful, alertStore cleared")
				}
			}
		}()
	})

	//http.HandleFunc("/healthz", healthzHandler)
	//http.HandleFunc("/readyz", readyzHandler)
	//http.HandleFunc("/startup", startupHandler)

	fmt.Printf("Starting server on %s...\n\n", *port)
	if err := http.ListenAndServe(*port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

//// 健康检查 - /healthz
//func healthzHandler(w http.ResponseWriter, r *http.Request) {
//	// 简单地返回存活状态
//	w.WriteHeader(http.StatusOK)
//	w.Write([]byte("ok"))
//}

//// 就绪检查 - /readyz
//func readyzHandler(w http.ResponseWriter, r *http.Request) {
//	if !isReady.Load().(bool) {
//		w.WriteHeader(http.StatusServiceUnavailable)
//		json.NewEncoder(w).Encode(map[string]string{
//			"status": "not ready",
//		})
//		return
//	}
//
//	// 模拟检查外部依赖（如数据库连接）
//	if !checkDatabaseConnection() {
//		w.WriteHeader(http.StatusServiceUnavailable)
//		json.NewEncoder(w).Encode(map[string]string{
//			"status": "not ready",
//			"error":  "database connection failed",
//		})
//		return
//	}
//
//	w.WriteHeader(http.StatusOK)
//	json.NewEncoder(w).Encode(map[string]string{
//		"status": "ready",
//	})
//}

//// 启动检查 - /startup
//func startupHandler(w http.ResponseWriter, r *http.Request) {
//	if !isStarted.Load().(bool) {
//		w.WriteHeader(http.StatusServiceUnavailable)
//		json.NewEncoder(w).Encode(map[string]string{
//			"status": "starting",
//		})
//		return
//	}
//
//	w.WriteHeader(http.StatusOK)
//	json.NewEncoder(w).Encode(map[string]string{
//		"status": "started",
//	})
//}

//// 模拟数据库连接检查
//func checkDatabaseConnection() bool {
//	// 这里可以实现实际的数据库连接检查逻辑
//	// 比如 ping 数据库或检查连接池状态
//	return true // 假设数据库连接正常
//}
