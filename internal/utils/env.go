package utils

import "os"

// GetEnv 获取环境变量，若不存在则返回默认值
func GetEnv(key string, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}
