package handler

import (
	"bytes"
	"fmt"
	"text/template"

	"smsAlert/models"
)

// GenerateMessage 提取告警结构体中所需要的字段，解析至template，生成短信内容
func GenerateMessage(a models.AlertItem) string {
	tmpl := `
{{ if eq .Status "firing" }}
【告警触发】
{{ else if eq .Status "resolved" }}
【告警恢复】
{{ end }}
域名：{{ .Labels.Instance }} 
实例名：{{ .Labels.InstanceName }} 
{{ if .Labels.Mountpoint }}
挂载点：{{ .Labels.Mountpoint }} 
{{ end }}
描述：{{ .Annotations.Summary }} 
严重性：{{ .Labels.Severity }} 
触发时间：{{ .StartsAt }} 
{{ if eq .Status "resolved" }}
恢复时间：{{ .EndsAt }} 
{{ end }}
`

	// 解析模板
	t, err := template.New("models").Parse(tmpl)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return ""
	}

	// 执行模板并生成结果
	var result bytes.Buffer
	err = t.Execute(&result, a)
	if err != nil {
		fmt.Println("Error executing template:", err)
		return ""
	}

	return result.String()
}
