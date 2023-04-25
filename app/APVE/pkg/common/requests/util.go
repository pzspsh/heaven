package requests

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
)

const (
	_tagName = "form"
)

//不处理嵌套的内容
func structToValues(i interface{}) url.Values {
	values := url.Values{}
	iVal := reflect.ValueOf(i)
	switch iVal.Kind() {
	case reflect.Map:
		for _, key := range iVal.MapKeys() {
			values.Add(key.String(), fmt.Sprint(iVal.MapIndex(key)))
		}
	default:
		reflectToValues(iVal, values)
	}
	return values
}

func mapToValues(params Value) url.Values {
	values := url.Values{}
	for k, v := range params {
		values.Set(k, v)
	}
	return values
}

func reflectToValues(iVal reflect.Value, values url.Values) {
	iVal = reflect.Indirect(iVal)
	typ := iVal.Type()
	//必须是一个结构体
	if typ.Kind() != reflect.Struct {
		panic(errors.New("only support struct or struct pointer"))
	}
	for i := 0; i < iVal.NumField(); i++ {
		var name string
		if typ.Field(i).Anonymous {
			//递归解析嵌套的匿名字段
			reflectToValues(iVal.Field(i), values)
			continue
		}
		if v, ok := typ.Field(i).Tag.Lookup(_tagName); ok {
			name = v
		} else {
			name = typ.Field(i).Name
		}
		values.Add(name, fmt.Sprint(iVal.Field(i)))
	}
}
