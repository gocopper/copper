package qb

import (
	"fmt"
	"reflect"
	"strings"
)

type fieldExtractor func(reflect.StructField, reflect.Value) (string, any, bool)

func ValuePlaceholders(model any) string {
	columns := extractRawColumns(model)
	placeholders := make([]string, len(columns))
	for i := range columns {
		placeholders[i] = "?"
	}
	return strings.Join(placeholders, ", ")
}

func Columns(model any, alias ...string) string {
	a := ""
	if len(alias) > 0 {
		a = alias[0]
	}

	columns := extractRawColumns(model)
	if a == "" {
		return strings.Join(columns, ", ")
	}

	prefixedColumns := make([]string, len(columns))
	for i, col := range columns {
		prefixedColumns[i] = fmt.Sprintf("%s.%s", a, col)
	}
	return strings.Join(prefixedColumns, ", ")
}

func SetColumns(model any) string {
	return strings.Join(extractSetColumns(model), ", ")
}

func Values(model any) []any {
	return extractValues(model, columnExtractor)
}

func SetValues(model any) []any {
	return extractValues(model, setColumnExtractor)
}

func ValuesAndSetValues(model any) []any {
	values := extractValues(model, columnExtractor)
	setValues := extractValues(model, setColumnExtractor)
	return append(values, setValues...)
}

func extractRawColumns(model any) []string {
	var columns []string
	processModelFields(model, func(field reflect.StructField, value reflect.Value) bool {
		tag := field.Tag.Get("db")
		if tag == "" || tag == "-" {
			return false
		}

		parts := strings.Split(tag, ",")
		columns = append(columns, parts[0])
		return true
	})
	return columns
}

func extractSetColumns(model any) []string {
	var columns []string
	processModelFields(model, func(field reflect.StructField, value reflect.Value) bool {
		tag := field.Tag.Get("db")
		if tag == "" || tag == "-" {
			return false
		}

		parts := strings.Split(tag, ",")
		if len(parts) >= 2 && parts[1] == "readonly" {
			return false
		}

		columns = append(columns, fmt.Sprintf("%s = ?", parts[0]))
		return true
	})
	return columns
}

func extractValues(model any, extractor fieldExtractor) []any {
	var values []any
	processModelFields(model, func(field reflect.StructField, value reflect.Value) bool {
		_, val, ok := extractor(field, value)
		if ok {
			values = append(values, val)
		}
		return ok
	})
	return values
}

func columnExtractor(field reflect.StructField, value reflect.Value) (string, any, bool) {
	tag := field.Tag.Get("db")
	if tag == "" || tag == "-" {
		return "", nil, false
	}

	parts := strings.Split(tag, ",")
	return parts[0], value.Interface(), true
}

func setColumnExtractor(field reflect.StructField, value reflect.Value) (string, any, bool) {
	tag := field.Tag.Get("db")
	if tag == "" || tag == "-" {
		return "", nil, false
	}

	parts := strings.Split(tag, ",")
	if len(parts) >= 2 && parts[1] == "readonly" {
		return "", nil, false
	}

	return fmt.Sprintf("%s = ?", parts[0]), value.Interface(), true
}

func processModelFields(model any, processor func(reflect.StructField, reflect.Value) bool) {
	t := reflect.TypeOf(model)
	v := reflect.ValueOf(model)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		if !field.IsExported() {
			continue
		}

		if field.Anonymous && field.Type.Kind() == reflect.Struct {
			processModelFields(v.Field(i).Interface(), processor)
			continue
		}

		processor(field, v.Field(i))
	}
}
