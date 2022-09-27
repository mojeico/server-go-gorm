package validation

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func ValidateJSON(input []byte, model interface{}) error {

	var data map[string]interface{}

	err := json.Unmarshal(input, &data)
	if err != nil {
		return err
	}

	keys := reflect.ValueOf(data).MapKeys()

	var tags []string
	for i := 1; i < reflect.TypeOf(model).NumField(); i++ {
		tag := strings.Split(string(reflect.TypeOf(model).Field(i).Tag), "\"")
		if tag == nil {
			return errors.New("can't get tag")
		}
		tags = append(tags, tag[1])
	}

	for _, keyValue := range keys {
		if !contains(keyValue.String(), tags) {
			return fmt.Errorf("doesn't exist json tag '%s'", keyValue.String())
		}
	}

	return nil
}

func contains(str string, arrStr []string) bool {
	for _, v := range arrStr {
		if str == v {
			return true
		}
	}
	return false
}
