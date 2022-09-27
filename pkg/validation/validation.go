package validation

import (
	"fmt"
	"reflect"
)

type ValidMethods interface {
	IsEmpty(value interface{}) error
	IsEmail(value string) error
	IsLength(value interface{}) error
	IsString(value interface{}) error
	IsNumber(value interface{}) error
}

type Validation struct {
	ValidMethods
}

func (valid *Validation) IsEmail(value string) bool {
	return RegexpEmail.MatchString(value)
}

func (valid *Validation) IsEmpty(value interface{}) bool {
	return fmt.Sprintf("%s", value) != ""
}

func (valid *Validation) IsLength(value interface{}, min int, max int) bool {
	var stringVal = fmt.Sprintf("%s", value)
	if len(stringVal) < min && len(stringVal) > max {
		return false
	}
	return true
}

func (valid *Validation) IsString(value interface{}) bool {
	return (reflect.TypeOf(value)).String() == "string"
}

func (valid *Validation) IsNumber(value interface{}) bool {
	return (reflect.TypeOf(value)).String() == "int"
}

func GetValidation() *Validation {
	return new(Validation)
}
