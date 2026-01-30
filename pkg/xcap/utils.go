package xcap

import "strings"

// SanitizeFilename 将字符串转换为安全的文件名
// 移除或替换不安全的文件名字符
func SanitizeFilename(name string) string {
	replacer := strings.NewReplacer(
		"/", "_", "\\", "_", ":", "_", "*", "_",
		"?", "_", "\"", "_", "<", "_", ">", "_",
		"|", "_", "\n", "_", "\r", "_",
	)
	result := replacer.Replace(name)
	if len(result) > 50 {
		result = result[:50]
	}
	result = strings.TrimSpace(result)
	if result == "" {
		result = "unnamed"
	}
	return result
}
