package database

import "fmt"

func GetCondition(queries map[string]string) string {
	var condition = ""
	for key, val := range queries {
		if val != "" && key != "id" {
			condition += " AND " + key + " = " + fmt.Sprintf("'%s'", val)
		}
	}
	return condition
}
