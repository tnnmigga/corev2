package tables

import (
	"reflect"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/tnnmigga/corev2/conc"
	"github.com/tnnmigga/corev2/log"
	"github.com/tnnmigga/corev2/system"
	"github.com/tnnmigga/corev2/utils"
	"golang.org/x/exp/constraints"
)

var (
	loader  ILoader
	store   = atomic.Value{}
	handles = map[string]func(Raw) (any, error){}
	checks  = []func(Tables) error{}
	tables  = []string{}
)

type ILoader interface {
	Load() (map[string]any, error)
}

func Init(from ILoader) {
	loader = from
	update()
	conc.Go(func() {
		ticker := time.NewTicker(time.Minute)
		for {
			select {
			case <-ticker.C:
				utils.ExecAndRecover(update)
			case <-system.RootCtx().Done():
				return
			}
		}
	})
}

func Register[T TableItem](h func(item Raw) (*T, error)) {
	var t T
	table := t.Table()
	handles[table] = func(m Raw) (any, error) {
		return h(m)
	}
	tables = append(tables, table)
}

func Check(h func(newConfs Tables) error) {
	checks = append(checks, h)
}

func Default[T TableItem](m Raw) (*T, error) {
	item := new(T)
	err := mapstructure.WeakDecode(m, item)
	return item, err
}

func Get[T TableItem, N constraints.Integer](id N) *T {
	var t T
	table := t.Table()
	m := read()
	confs, ok := m[table]
	if !ok {
		log.Errorf("table not found %s", table)
		return nil
	}
	item, ok := confs[int(id)]
	if !ok {
		// log.Errorf("table item not found %s %d", table, id)
		return nil
	}
	return item.(*T)
}

func Range[T TableItem](h func(*T) (stop bool)) {
	var t T
	table := t.Table()
	m := read()
	for _, v := range m[table] {
		item := v.(*T)
		h(item)
	}
}

func All[T TableItem]() []*T {
	var t T
	table := t.Table()
	m := read()
	items := make([]*T, 0, len(m[table]))
	for _, v := range m[table] {
		items = append(items, v.(*T))
	}
	return items
}

func Filter[T TableItem](h func(*T) bool) []*T {
	var t T
	table := t.Table()
	m := read()
	var items []*T
	for _, v := range m[table] {
		conf := v.(*T)
		if h(conf) {
			items = append(items, conf)
		}
	}
	return items
}

func Find[T TableItem](is func(*T) bool) *T {
	var conf *T
	Range(func(item *T) (stop bool) {
		if is(item) {
			conf = item
			return true
		}
		return false
	})
	return conf
}

type TableItem interface {
	Table() string
	GetID() int
}

type Tables map[string]map[int]any

func (c Tables) Get(name string, ID int) TableItem {
	return c[name][ID].(TableItem)
}

func (c Tables) Range(name string, h func(item TableItem) bool) {
	table := c[name]
	for _, item := range table {
		if h(item.(TableItem)) {
			return
		}
	}
}

func (c Tables) Filter(name string, h func(item TableItem) bool) []TableItem {
	var items []TableItem
	c.Range(name, func(item TableItem) bool {
		if h(item) {
			items = append(items, item)
		}
		return false
	})
	return items
}

func (c Tables) All(name string) []TableItem {
	var items []TableItem
	table := c[name]
	for _, item := range table {
		items = append(items, item.(TableItem))
	}
	return items
}

func (c Tables) Save(item TableItem) {
	c[item.Table()][item.GetID()] = item
}

func save(m Tables) {
	store.Store(m)
}

func read() Tables {
	return store.Load().(Tables)
}

func update() {
	confs, err := loader.Load()
	if err != nil {
		log.Panic(err)
	}
	newTables := Tables{}
	for _, name := range tables {
		conf := confs[name]
		h := handles[name]
		sub := map[int]any{}
		switch reflect.TypeOf(conf).Kind() {
		case reflect.Slice:
			l := conf.([]any)
			for i, v := range l {
				if v == nil {
					continue
				}
				item, err := h(v.(Raw))
				if err != nil {
					log.Panicf("handle table %s %v error", name, v)
				}
				if item == nil { // 无效的数据直接跳过
					continue
				}
				sub[i+1] = item
			}
		case reflect.Map:
			m := conf.(map[string]any)
			for k, v := range m {
				item, err := h(v.(Raw))
				if err != nil {
					log.Panicf("handle table %s %v error", name, v)
				}
				if v := reflect.ValueOf(item); v.Kind() == reflect.Ptr && v.IsNil() { // 无效的数据直接跳过
					continue
				}
				id, err := strconv.Atoi(k)
				if err != nil {
					log.Panicf("handler table parse id error %v", id)
				}
				sub[id] = item
			}
		}
		newTables[name] = sub
	}
	for _, check := range checks {
		if err := check(newTables); err != nil {
			log.Panic(err)
		}
	}
	save(newTables)
}
