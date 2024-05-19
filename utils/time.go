package utils

import "time"

// 获取纳秒级时间戳
func NowNs() time.Duration {
	return time.Duration(time.Now().UnixNano())
}

// 获取秒级时间戳
func NowSec() float64 {
	return float64(time.Now().UnixNano()) / float64(time.Second)
}
