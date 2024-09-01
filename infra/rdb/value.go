package rdb

import (
	"errors"

	"github.com/gomodule/redigo/redis"
)

var (
	ErrUintOverflow = errors.New("redigo: overflow converting uint64 to uint32")
)

type Result struct {
	value interface{}
	err   error
}

func (r Result) Value() (interface{}, error) {
	return r.value, r.err
}

func (r Result) Error() error {
	return r.err
}

func (r Result) Int16() (int16, error) {
	n, err := redis.Int(r.value, r.err)
	return int16(n), err
}

func (r Result) Int32() (int32, error) {
	n, err := redis.Int(r.value, r.err)
	return int32(n), err
}

func (r Result) Int64() (int64, error) {
	return redis.Int64(r.value, r.err)
}

func (r Result) Int() (int, error) {
	return redis.Int(r.value, r.err)
}

func (r Result) Uint32() (uint32, error) {
	n, err := redis.Uint64(r.value, r.err)
	if n > 0xFFFFFFFF {
		return 0, ErrUintOverflow
	}
	return uint32(n), err
}

func (r Result) Uint64() (uint64, error) {
	return redis.Uint64(r.value, r.err)
}

func (r Result) Float64() (float64, error) {
	return redis.Float64(r.value, r.err)
}

func (r Result) String() (string, error) {
	return redis.String(r.value, r.err)
}

func (r Result) Bytes() ([]byte, error) {
	return redis.Bytes(r.value, r.err)
}

func (r Result) Bool() (bool, error) {
	return redis.Bool(r.value, r.err)
}

func (r Result) Float64s() ([]float64, error) {
	return redis.Float64s(r.value, r.err)
}

func (r Result) Strings() ([]string, error) {
	return redis.Strings(r.value, r.err)
}

func (r Result) ByteSlices() ([][]byte, error) {
	return redis.ByteSlices(r.value, r.err)
}

func (r Result) Int16s() ([]int16, error) {
	m, err := redis.Int64s(r.value, r.err)
	if err != nil {
		return nil, err
	}
	n := make([]int16, len(m))
	for i, v := range m {
		n[i] = int16(v)
	}
	return n, nil
}

func (r Result) Int32s() ([]int32, error) {
	m, err := redis.Int64s(r.value, r.err)
	if err != nil {
		return nil, err
	}
	n := make([]int32, len(m))
	for i, v := range m {
		n[i] = int32(v)
	}
	return n, nil
}

func (r Result) Int64s() ([]int64, error) {
	return redis.Int64s(r.value, r.err)
}

func (r Result) Ints() ([]int, error) {
	return redis.Ints(r.value, r.err)
}

func (r Result) StringMap() (map[string]string, error) {
	return redis.StringMap(r.value, r.err)
}

func (r Result) Int32Map() (map[string]int32, error) {
	m, err := redis.IntMap(r.value, r.err)
	if err != nil {
		return nil, err
	}
	n := make(map[string]int32, len(m))
	for k, v := range m {
		n[k] = int32(v)
	}
	return n, nil
}

func (r Result) IntMap() (map[string]int, error) {
	return redis.IntMap(r.value, r.err)
}

func (r Result) Int64Map() (map[string]int64, error) {
	return redis.Int64Map(r.value, r.err)
}

func (r Result) Positions() ([]*[2]float64, error) {
	return redis.Positions(r.value, r.err)
}

func (r Result) Values() ([]interface{}, error) {
	return redis.Values(r.value, r.err)
}

// obj为一个指针对象
func (r Result) ScanStruct(obj interface{}) error {
	v, err := r.Values()
	if err != nil {
		return err
	}
	return redis.ScanStruct(v, obj)
}
