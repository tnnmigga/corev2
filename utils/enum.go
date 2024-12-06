package utils

import "reflect"

func NewEnum[T any]() T {
	var enum T
	valueOf := reflect.ValueOf(&enum)
	typeOf := reflect.TypeOf(&enum)
	valueOf = valueOf.Elem()
	typeOf = typeOf.Elem()
	for i := 0; i < valueOf.NumField(); i++ {
		field := valueOf.Field(i)
		fieldType := typeOf.Field(i)
		if !field.CanSet() {
			continue
		}
		switch k := field.Kind(); k {
		default:
			panic("enum field type must be string or int")
		case reflect.String:
			field.SetString(fieldType.Name)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			field.SetInt(int64(i))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			field.SetUint(uint64(i))
		}
	}
	return enum
}
