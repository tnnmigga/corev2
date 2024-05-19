package utils

import (
	"math/rand"

	"github.com/tnnmigga/corev2/algorithm"
	"github.com/tnnmigga/corev2/zlog"
	"golang.org/x/exp/constraints"
)

// RandomInterval 生成一个随机数, 范围是[low, high]
func RandomInterval[T constraints.Integer](low, high T) T {
	a, b := int64(low), int64(high)
	return T(rand.Int63n(b-a+1) + a)
}

// RandomIntervalN 生成指定数量的随机数, 范围是[low, high], 保证不重复
func RandomIntervalN[T constraints.Signed](low, high T, num T) algorithm.Set[T] {
	maxNum := high - low + 1
	if maxNum < num {
		zlog.Errorf("max random num not enough")
		return nil
	}
	var set algorithm.Set[T]
	if float64(num)/float64(maxNum) < 0.75 {
		set = make(algorithm.Set[T], num)
		for len(set) < int(num) {
			v := RandomInterval(low, high)
			if set.Find(v) {
				continue
			}
			set.Insert(v)
		}
		return set
	}
	set = make(algorithm.Set[T], maxNum)
	for i := low; i <= high; i++ {
		set.Insert(i)
	}
	count := high - low - num + 1
	for key := range set {
		delete(set, key)
		count--
		if count <= 0 {
			break
		}
	}
	return set
}

// 生成生成指定长度的只包含字母和数字的随机字符串
func GenerateToken(charLen int) string {
	charset := []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	token := make([]byte, charLen)
	for i := 0; i < charLen; i++ {
		token[i] = charset[rand.Intn(len(charset))]
	}
	return string(token)
}
