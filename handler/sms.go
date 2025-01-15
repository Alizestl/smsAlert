package handler

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"smsAlert/models"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// getSign 生成API所需要的 HMAC-SHA256 签名
func getSign(data, secret string) string {
	secretBytes, _ := hex.DecodeString(secret)
	h := hmac.New(sha256.New, secretBytes)
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

// httpPost 发送给API POST请求
func httpPost(url string, data []byte, headers map[string]string) (string, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// SendSMS 调用短信API
func SendSMS(phone string, s *models.AlertStore) error {
	key := "xxx"
	secret := "xxx"

	timestamp := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)

	for _, a := range s.Alerts {
		var message string
		for _, v := range a.Alerts {
			message = GenerateMessage(v)

			// 将换行符\n，替换为\\n，不然短信API接收POST包后会报json解析错误
			message = strings.ReplaceAll(message, "\n", "\\n")
			logrus.Info(fmt.Sprintf(message))

			data := fmt.Sprintf(`{"phone": "%s", "msg": "%s"}`, phone, message)
			dataBytes := []byte(data)

			dataToSig := key + "/api/sms/sendMsg" + timestamp + data
			sig := getSign(dataToSig, secret)

			url := "https://url/api/sms/sendMsg"
			headers := map[string]string{
				"Accept":          "application/json",
				"Content-Type":    "application/json;charset=utf-8",
				"X-Api-Key":       key,
				"X-Api-Timestamp": timestamp,
				"X-Api-Signature": sig,
			}

			start := time.Now()
			resp, err := httpPost(url, dataBytes, headers)
			logrus.Info("短信API响应包: ", resp)
			logrus.Info(fmt.Sprintf("发包花费时间: %s\n", time.Since(start)))
			logrus.Info(fmt.Sprintf("当前时间: %s\n", time.Now().Format("2006-01-02 15:04:05")))

			if err != nil {
				return err
			}
		}
	}
	return nil
}
