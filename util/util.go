package util

import (
	"hotbox-adm-backend/pkg/constant"

	"golang.org/x/exp/rand"
)

// GenerateRandomString生成随机字符串
func GenerateRandomString(length int) string {
	const characters = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = characters[rand.Intn(len(characters))]
	}
	return string(result)
}

// 转换活动类型为字符串
func StringActivityType[T int | int64 | int32](activityType T) string {
	switch activityType {
	case constant.ACTIVITY_TYPE_COMBINATION:
		return "合成"
	case constant.ACTIVITY_TYPE_UPGRADE:
		return "升级"
	case constant.ACTIVITY_TYPE_REPLACE:
		return "置换"
	case constant.ACTIVITY_TYPE_DECOMPOSE:
		return "分解"
	}
	return ""
}
