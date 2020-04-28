package util

import (
	"strings"
)

// 根据Label生成对应的labelSelector
// k8s-app=nginx,app=nginx01 多标签查询
func GenerateLabelSelector(selector map[string]string) string {
	var labelSelector string
	for k, v := range selector {
		labelSelector += k + "=" + v + ","
	}
	return strings.TrimRight(labelSelector, ",")
}
