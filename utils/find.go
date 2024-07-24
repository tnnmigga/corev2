package utils

// 列表如果长度大于0, 返回第一个元素, 否则返回默认值
func FirstOrDefault[T any](s []T, defaultVal T) T {
	if len(s) > 0 {
		return s[0]
	}
	return defaultVal
}

// 筛选列表中符合规则的元素
func Filter[T any](items []T, fn func(T) bool) (res []T) {
	for _, item := range items {
		if fn(item) {
			res = append(res, item)
		}
	}
	return res
}

// 查找元素在列表中的索引, 不存在返回-1
func Index[T comparable](slice []T, n T) int {
	for i, v := range slice {
		if v == n {
			return i
		}
	}
	return -1
}

// 判断元素是否在列表中
func Contain[T comparable](slice []T, n T) bool {
	return Index(slice, n) != -1
}
