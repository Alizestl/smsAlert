package handler

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"smsAlert/database"
	"smsAlert/models"
	"strings"
)

type Receivers struct {
	Phone string
}

// Query 查询手机号
func Query(s *models.AlertStore, configPath *string) string {
	// 获取告警组，对告警组做去重
	var groups []string
	groupSet := make(map[string]struct{})
	for _, a := range s.Alerts {
		groupSet[a.Receiver] = struct{}{}
	}
	for group := range groupSet {
		groups = append(groups, group)
		logrus.Info(fmt.Sprintf("group to be informed: %s ", group))
	}

	var allPhones []string
	for _, group := range groups {
		phones := queryDB(group, configPath)
		if phones != "" {
			allPhones = append(allPhones, phones)
		}
	}
	uniquePhones := removeDuplicates(strings.Split(strings.Join(allPhones, ","), ","))
	finalPhones := strings.Join(uniquePhones, ",")
	logrus.Info(fmt.Sprintf("Numbers to be informed: %s\n\n", finalPhones))

	return finalPhones
}

func queryDB(group string, configPath *string) string {
	// 初始化数据库连接
	db, err := database.Init(*configPath)
	if err != nil {
		fmt.Println("database init error")
	}

	// 创建接收电话字符串的结构体切片
	var receiver []Receivers

	// 执行原生 SQL 查询
	//start := time.Now()
	result := db.Raw("SELECT phone FROM test.`phone` WHERE Name in (SELECT Name FROM test.`Group` WHERE Groupname = ?)", group).Scan(&receiver)
	//fmt.Printf("Query took: %s\n", time.Since(start))

	if result.Error != nil {
		fmt.Println("Error:", result.Error)
		return ""
	} else if result.RowsAffected == 0 {
		fmt.Println("No records found")
		return ""
	}

	var phones []string
	for _, r := range receiver {
		phones = append(phones, r.Phone)
	}

	return strings.Join(phones, ",")
}

// 工具函数，用于切片去重
func removeDuplicates(input []string) []string {
	// 使用 map 来记录已存在的元素
	seen := make(map[string]struct{})
	result := make([]string, 0, len(input)) // 预分配容量

	for _, item := range input {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}

	return result
}
