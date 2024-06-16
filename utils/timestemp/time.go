package timestemp

import (
	"time"

	"golang.org/x/exp/constraints"
)

// 获取纳秒级时间戳
func NowNs() time.Duration {
	return time.Duration(time.Now().UnixNano())
}

// 获取秒级时间戳
func NowSec[N constraints.Float | constraints.Integer]() N {
	return N(float64(time.Now().UnixNano()) / float64(time.Second))
}
