package utils

import (
	"fmt"
	"reflect"
)

func StructToMap(obj interface{}) (map[string]interface{}, error) {

	if reflect.ValueOf(obj).Kind() != reflect.Ptr || reflect.ValueOf(obj).IsNil() {
		return nil, fmt.Errorf("input must be a non-nil pointer")
	}

	val := reflect.ValueOf(obj).Elem()

	output := make(map[string]interface{})

	for i := 0; i < val.NumField(); i++ {

		valueField := val.Field(i)

		typeField := val.Type().Field(i)

		fieldName := typeField.Name

		output[fieldName] = valueField.Interface()
	}

	return output, nil
}
