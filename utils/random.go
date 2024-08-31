package utils

import (
	"math/rand"

	"github.com/tnnmigga/corev2/log"
	"golang.org/x/exp/constraints"
)

func RandomInterval[T constraints.Integer](low, high T) T {
	a, b := int64(low), int64(high)
	return T(rand.Int63n(b-a+1) + a)
}

// RandomIntervalN 生成指定数量的随机数, 范围是[low, high], 保证不重复
func RandomIntervalN[T constraints.Integer](low, high T, num T) map[T]struct{} {
	maxNum := high - low + 1
	if maxNum < num {
		log.Errorf("max random num not enough")
		return nil
	}
	var m map[T]struct{}
	if float64(num)/float64(maxNum) < 0.75 {
		m = make(map[T]struct{}, num)
		for len(m) < int(num) {
			v := RandomInterval(low, high)
			if _, ok := m[v]; ok {
				continue
			}
			m[v] = struct{}{}
		}
		return m
	}
	m = make(map[T]struct{}, maxNum)
	for i := low; i <= high; i++ {
		m[i] = struct{}{}
	}
	count := high - low - num + 1
	for key := range m {
		delete(m, key)
		count--
		if count <= 0 {
			break
		}
	}
	return m
}
