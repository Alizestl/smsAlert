package models

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"smsAlert/logger"
	"sync"
)

type (
	AlertItem struct {
		Status string `json:"status"`
		Labels struct {
			Alertname    string `json:"alertname"`
			Cluster      string `json:"cluster"`
			Deployment   string `json:"deployment"`
			Namespace    string `json:"namespace"`
			Instance     string `json:"instance"`
			InstanceName string `json:"instance_name"`
			Mountpoint   string `json:"mountpoint"`
			Severity     string `json:"Severity"`
			Team         string `json:"team"`
		} `json:"labels"`
		Annotations struct {
			Summary string `json:"summary"`
		} `json:"annotations"`
		StartsAt     string `json:"startsAt"`
		EndsAt       string `json:"endsAt"`
		GeneratorURL string `json:"generatorURL"`
		Fingerprint  string `json:"fingerprint"`
	}

	Alert struct {
		Receiver    string      `json:"receiver"`
		Status      string      `json:"status"`
		Alerts      []AlertItem `json:"alerts"`
		GroupLabels struct {
			Alertname string `json:"alertname"`
			Cluster   string `json:"cluster"`
		} `json:"groupLabels"`
		CommonLabels struct {
			Alertname  string `json:"alertname"`
			Cluster    string `json:"cluster"`
			Deployment string `json:"deployment"`
			Namespace  string `json:"namespace"`
			Team       string `json:"team"`
		} `json:"commonLabels"`
		CommonAnnotations struct {
			Summary string `json:"summary"`
		} `json:"commonAnnotations"`
		ExternalURL     string `json:"externalURL"`
		Version         string `json:"version"`
		GroupKey        string `json:"groupKey"`
		TruncatedAlerts int    `json:"truncatedAlerts"`
	}

	AlertStore struct {
		sync.Mutex
		Capacity int
		Alerts   []*Alert
	}
)

// PostHandler 接收AM发送的告警json，解析到Alert结构体，存储到alertStore.alerts中
func (s *AlertStore) PostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		logger.Log.Warnf("Unsupported HTTP method: %s", r.Method)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576) // 限制请求体为 1MB
	defer r.Body.Close()

	// 读取、打印POST的内容
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Fail:", http.StatusInternalServerError)
		log.Printf("Error reading requests body %v\n", err)
		return
	}
	logger.Log.Infof("Received raw JSON:\n %s\n", string(body)) // 记录raw JSON到日志

	log.Printf("Receive JSON:\n %s\n", body) // 在webhook中打印原始POST数据

	// 将接收到的json数据解析为Alert结构体
	var alert Alert
	err = json.Unmarshal(body, &alert)
	if err != nil {
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		logger.Log.Errorf("Error parsing JSON: %v", err)
		return
	}

	if alert.Receiver == "" || alert.Status == "" || len(alert.Alerts) == 0 {
		http.Error(w, "Missing required alert fields", http.StatusBadRequest)
		logger.Log.Warn("Received alert with missing fields")
		return
	}

	s.Lock()
	defer s.Unlock()

	// 将转化后的Alert结构体添加到alertStore.Alert切片中
	s.Alerts = append(s.Alerts, &alert)
	
	if len(s.Alerts) > s.Capacity {
		s.Alerts = s.Alerts[1:]
	}
	logger.Log.Infof("Alert stored: Receiver=%s, Status=%s", alert.Receiver, alert.Status)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Alert received successfully"))
}
