package form

import (
	"fmt"
	"reflect"
)

type customResponse struct {
	response map[string]any
}

func (c customResponse) String(key string) string {
	return c.response[key].(string)
}

func (c customResponse) Int(key string) int {
	if reflect.TypeOf(c.response[key]).Kind() == reflect.Float64 {
		return int(c.response[key].(float64))
	}
	return c.response[key].(int)
}

func (c customResponse) Float(key string) float64 {
	return c.response[key].(float64)
}

func (c customResponse) Bool(key string) bool {
	return c.response[key].(bool)
}

func (c customResponse) Bind(v any) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("bind expects a pointer to a struct, got %T", v)
	}

	val = val.Elem()
	typ := val.Type()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)

		tag := field.Tag.Get("form")
		if tag == "" {
			continue
		}

		mapValue, exists := c.response[tag]
		if !exists {
			continue
		}

		if !fieldValue.CanSet() {
			return fmt.Errorf("cannot set field %s", field.Name)
		}

		switch fieldValue.Kind() {
		case reflect.String:
			if strVal, ok := mapValue.(string); ok {
				fieldValue.SetString(strVal)
			} else {
				return fmt.Errorf("cannot assign %v to string field %s", mapValue, field.Name)
			}

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if floatVal, ok := mapValue.(float64); ok { // JSON numbers are float64 by default
				fieldValue.SetInt(int64(floatVal))
			} else if intVal, ok := mapValue.(int); ok {
				fieldValue.SetInt(int64(intVal))
			} else {
				return fmt.Errorf("cannot assign %v to int field %s", mapValue, field.Name)
			}

		case reflect.Float32, reflect.Float64:
			if floatVal, ok := mapValue.(float64); ok {
				fieldValue.SetFloat(floatVal)
			} else {
				return fmt.Errorf("cannot assign %v to float field %s", mapValue, field.Name)
			}

		case reflect.Bool:
			if boolVal, ok := mapValue.(bool); ok {
				fieldValue.SetBool(boolVal)
			} else {
				return fmt.Errorf("cannot assign %v to bool field %s", mapValue, field.Name)
			}

		default:
			return fmt.Errorf("unsupported field type %s for field %s", field.Type, field.Name)
		}
	}

	return nil
}

type CustomResponse interface {
	String(key string) string
	Int(key string) int
	Float(key string) float64
	Bool(key string) bool
	Bind(v any) error
}
