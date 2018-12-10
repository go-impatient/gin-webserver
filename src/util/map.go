package util

import (
	"fmt"
	"reflect"
	"strconv"
)

type JSONObject map[string]interface{}

func (c JSONObject) Set(k interface{}, v interface{}) JSONObject {
	c[toString(k)] = v
	return c
}

func (c JSONObject) Bool(k interface{}, def bool) bool {
	if v, ok := c[toString(k)].(bool); ok {
		return v
	}
	return def
}

func (c JSONObject) Int(k interface{}, def int) int {
	if v, ok := c[toString(k)].(int); ok {
		return v
	}
	return def
}

func (c JSONObject) String(k interface{}, def string) string {
	if v, ok := c[toString(k)].(string); ok {
		return v
	}
	return def
}

func (c JSONObject) Object(k interface{}) JSONObject {
	key := toString(k)
	if v, ok := c[key].(JSONObject); ok {
		return v
	}
	v := make(JSONObject)
	c[key] = v
	return v
}

func toString(v interface{}) string {
	r := reflect.ValueOf(v)
	switch r.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(r.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(r.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(r.Float(), 'f', -1, 64)
	case reflect.String:
		return r.String()
	default:
		return fmt.Sprintf("%#v", v)
	}
}
