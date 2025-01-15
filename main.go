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

func main() {
	// 打印当前目录，为载入settings.yml出错时debug提供便利
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory:", err)
	} else {
		fmt.Println("Current working directory:", dir)
	}

	// AM生成的告警为json格式，将其解析成Alert结构体
	// 创建全局变量alertStore存储Alert结构体
	port := flag.String("addr", ":5000", "address to listen for webhook")
	capacity := flag.Int("cap", 64, "capacity of alertStore")
	flag.Parse()
	s := &models.AlertStore{
		Capacity: *capacity,
	}

	// 启动服务器，接收告警
	go Webhook(port, s)
	//select {}

	// 监听中断信号，阻塞主goroutine，收到终止信号后解除阻塞并打印日志
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	<-signalChan
	fmt.Println("Shutting down webhook...")
	logger.Log.Info(fmt.Sprintf("Shutting down webhook..."))
}

func Webhook(port *string, s *models.AlertStore) {
	// 初始化rolling-logger
	logger.InitLogger()

	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		// 接收AM发送的告警json
		s.PostHandler(w, r)
		phone := handler.Query(s)
		if phone != "" {
			if err := handler.SendSMS(phone, s); err != nil {
				logger.Log.Errorf("Failed to send SMS: %v", err)
			} else {
				s.Lock()
				s.Alerts = nil // 清空alertStore存储的告警，否则下次调用SendSMS会重复发送短信
				s.Unlock()
				logger.Log.Info("SendSMS successful, alertStore cleared")
			}
		}
	})

	fmt.Printf("Starting server on %s...\n\n", *port)
	if err := http.ListenAndServe(*port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
