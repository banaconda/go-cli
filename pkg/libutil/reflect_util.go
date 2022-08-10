package libutil

import (
	"fmt"
	"reflect"
)

func GetFieldNameList(t reflect.Type) []string {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	fieldNameSlice := make([]string, 0)
	for _, f := range reflect.VisibleFields(t) {
		if !f.IsExported() {
			continue
		}
		fieldNameSlice = append(fieldNameSlice, f.Name)
	}

	return fieldNameSlice
}

func PrintStructOne(elem any) {
	elems := make([]any, 0)
	elems = append(elems, elem)
	PrintStructAll(elems)
}

func PrintStructAll(elems interface{}) {
	slice, ok := CreateAnyTypeSlice(elems)
	if !ok || len(slice) == 0 {
		return
	}

	t := reflect.TypeOf(slice[0])
	fieldNameList := GetFieldNameList(t)
	fieldNum := len(fieldNameList)
	maxWidths := make([]int, fieldNum)
	lines := make([][]string, 0)

	line := make([]string, 0)
	for i, fieldName := range fieldNameList {
		maxWidths[i] = len(fieldName)
		line = append(line, fieldName)
	}
	lines = append(lines, line)

	for _, elem := range slice {
		e := reflect.ValueOf(elem)
		if e.Kind() == reflect.Ptr {
			e = e.Elem()
		}

		line := make([]string, 0)
		for i, name := range fieldNameList {
			value := e.FieldByName(name).Interface()
			valueString := fmt.Sprintf("%v", value)
			line = append(line, valueString)
			if maxWidths[i] < len(valueString) {
				maxWidths[i] = len(valueString)
			}
		}
		lines = append(lines, line)
	}

	for _, line := range lines {
		for i, value := range line {
			fmt.Printf("%*s", maxWidths[i]+1, value)
		}
		fmt.Println()
	}
}

// Interface {} to [] interface {}
func CreateAnyTypeSlice(slice interface{}) ([]interface{}, bool) {
	val, ok := isSlice(slice)

	if !ok {
		return nil, false
	}

	sliceLen := val.Len()

	out := make([]interface{}, sliceLen)

	for i := 0; i < sliceLen; i++ {
		out[i] = val.Index(i).Interface()
	}

	return out, true
}

// Determine whether it is slcie data
func isSlice(arg interface{}) (val reflect.Value, ok bool) {
	val = reflect.ValueOf(arg)

	if val.Kind() == reflect.Slice {
		ok = true
	}

	return
}
