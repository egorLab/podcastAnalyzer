package misc

import (
	"reflect"
	"strconv"
)

// SQL syntax helper structs

type TableMapper struct {
	ColumnNames string
	Entity      interface{}
}

type Request struct {
	Field interface{}
	Value interface{}
}

func GetFieldsAndWildcards(entry interface{}) (fields []interface{}, wildcardStr string) {
	// create interface to pass and wildcards
	wildcard := "("
	s := reflect.ValueOf(entry)
	fieldsToFill := make([]interface{}, s.NumField())
	for i := 0; i < s.NumField(); i++ {
		fieldsToFill[i] = s.Field(i).Interface()
		wildcard += "$" + strconv.Itoa(i+1)
		if i != s.NumField()-1 {
			wildcard += ", " // TODO better way to construct statement?
		}
	}
	wildcard += ")"

	return fieldsToFill, wildcard
}
