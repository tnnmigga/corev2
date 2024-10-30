package tables

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/tnnmigga/corev2/log"
)

type Raw map[string]any

func (raw Raw) TranslateJSON(keys ...string) error {
	for _, key := range keys {
		value := raw[key]
		str, ok := value.(string)
		if len(str) == 0 {
			raw[key] = nil
			continue
		}
		if !ok {
			return fmt.Errorf("TranslateJSON must string %v", raw)
		}
		var result any
		err := json.Unmarshal([]byte(str), &result)
		if err != nil {
			return err
		}
		raw[key] = result
	}
	return nil
}

func (raw Raw) TransformExpression(keys ...string) error {
	for _, key := range keys {
		txt := raw[key].(string)
		e, err := parseExpression(txt)
		if err != nil {
			return err
		}
		raw[key] = e
	}
	return nil
}

func parseExpression(s string) (*Expression, error) {
	if len(s) == 0 {
		return nil, nil
	}
	subs := strings.Split(s, " or ")
	if len(subs) > 1 {
		expr := &Expression{Op: "or", Subs: make([]*Expression, 0, len(subs))}
		for _, sub := range subs {
			subexpr, err := parseExpression(strings.TrimSpace(sub))
			if err != nil {
				return nil, err
			}
			expr.Subs = append(expr.Subs, subexpr)
		}
		return expr, nil
	}
	subs = strings.Split(s, " and ")
	if len(subs) > 1 {
		expr := &Expression{Op: "and", Subs: make([]*Expression, 0, len(subs))}
		for _, sub := range subs {
			subexpr, err := parseExpression(strings.TrimSpace(sub))
			if err != nil {
				return nil, err
			}
			expr.Subs = append(expr.Subs, subexpr)
		}
		return expr, nil
	}
	replace := regexp.MustCompile(`\s+`)
	s = replace.ReplaceAllString(s, "") // 删除空字符
	match := regexp.MustCompile(`([a-zA-Z0-9_]+)([><=]+)(\d+)`)
	items := match.FindStringSubmatch(s)
	if len(items) != 4 {
		return nil, fmt.Errorf("parseExpression error %v", s)
	}
	value, err := strconv.Atoi(items[3])
	if err != nil {
		return nil, fmt.Errorf("parseExpression error %v", s)
	}
	return &Expression{
		Op:    items[2],
		Key:   items[1],
		Value: int64(value),
	}, nil
}

type Expression struct {
	Op    string
	Key   string
	Value int64
	Subs  []*Expression
}

func (e *Expression) Check(param map[string]int64) bool {
	if e == nil {
		return true
	}
	if param == nil {
		return false
	}
	switch e.Op {
	case ">=":
		return param[e.Key] >= e.Value
	case "<=":
		return param[e.Key] <= e.Value
	case "=", "==":
		return param[e.Key] == e.Value
	case "<":
		return param[e.Key] < e.Value
	case ">":
		return param[e.Key] > e.Value
	case "and":
		for _, sub := range e.Subs {
			if !sub.Check(param) {
				return false
			}
		}
		return true
	case "or":
		for _, sub := range e.Subs {
			if sub.Check(param) {
				return true
			}
		}
		return false
	}
	log.Errorf("Check Expression invalid Op %s", e.Op)
	return false
}

func (raw Raw) TransformTime(keys ...string) error {
	for _, key := range keys {
		txt := raw[key].(string)
		if len(txt) == 0 {
			raw[key] = 0
			continue
		}
		t, err := time.Parse(time.DateTime, txt)
		if err != nil {
			t, err = time.Parse("2006/01/02 15:04:05", txt)
		}
		if err != nil {
			return err
		}
		raw[key] = t.UTC().Unix()
	}
	return nil
}
