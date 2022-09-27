package sort

import (
	"fmt"
	"log"
	"reflect"
)

const (
	ASC = iota
	DESC
)

type Sort struct {
}

func GetSortPackage() *Sort {
	return &Sort{}
}

func interfaceToSlice(slice interface{}) []interface{} {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		log.Println("InterfaceToSlice() given a non-slice type")
	}

	if s.IsNil() {
		return nil
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret
}

func reverseSlice(s interface{}) {
	n := reflect.ValueOf(s).Len()
	swap := reflect.Swapper(s)
	for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
		swap(i, j)
	}
}

func swap(x, y *interface{}) {
	temp := *x
	*x = *y
	*y = temp
}

func (s *Sort) SortDataByStructField(field string, method byte, data interface{}) interface{} {

	res := interfaceToSlice(data)
	n := len(res)

	for i := 0; i < n; i++ {
		for j := 0; j < n-i-1; j++ {

			r1 := reflect.ValueOf(res[j+1])
			f1 := reflect.Indirect(r1).FieldByName(field)

			r2 := reflect.ValueOf(res[j])
			f2 := reflect.Indirect(r2).FieldByName(field)

			switch f1.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				if f1.Int() < f2.Int() {
					swap(&res[j], &res[j+1])
				}
			case reflect.Float32, reflect.Float64:
				if f1.Float() < f2.Float() {
					swap(&res[j], &res[j+1])
				}
			case reflect.String:
				if f1.String() < f2.String() {
					swap(&res[j], &res[j+1])
				}
			default:
				if fmt.Sprint(f1) < fmt.Sprint(f2) {
					swap(&res[j], &res[j+1])
				}
			}
		}
	}

	if method == DESC {
		reverseSlice(res)
		return res
	}

	return res
}
