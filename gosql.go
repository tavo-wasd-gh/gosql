package gosql

import (
	"fmt"
	"database/sql"
	"reflect"
)

func ScanRows(rows *sql.Rows, dest interface{}) error {
	v := reflect.ValueOf(dest)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("dest must be a pointer to a struct")
	}

	structValue := v.Elem()
	structType := structValue.Type()

	keyFieldMap := map[string]int{}
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		tag := field.Tag.Get("db")
		if tag == "key" {
			keyFieldMap[field.Name] = i
		}
	}

	columns, err := rows.Columns()
	if err != nil {
		return fmt.Errorf("failed to get columns: %v", err)
	}

	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return fmt.Errorf("failed to scan row: %v", err)
		}

		for i, colName := range columns {
			if fieldIdx, ok := keyFieldMap[colName]; ok {
				field := structValue.Field(fieldIdx)
				if field.CanSet() {
					field.Set(reflect.ValueOf(values[i]))
				}
			}
		}
	}

	return nil
}

func ScanRow(row *sql.Row, dest interface{}) error {
	v := reflect.ValueOf(dest)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("dest must be a pointer to a struct")
	}

	structValue := v.Elem()
	structType := structValue.Type()

	keyFieldMap := map[string]int{}
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		tag := field.Tag.Get("db")
		if tag != "" {
			keyFieldMap[tag] = i
		}
	}

	columns := make([]string, 0, len(keyFieldMap))
	for key := range keyFieldMap {
		columns = append(columns, key)
	}

	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	if err := row.Scan(valuePtrs...); err != nil {
		return fmt.Errorf("failed to scan row: %v", err)
	}

	for i, colName := range columns {
		if fieldIdx, ok := keyFieldMap[colName]; ok {
			field := structValue.Field(fieldIdx)
			if field.CanSet() {
				val := reflect.ValueOf(values[i])
				if val.Kind() == reflect.Ptr && val.IsNil() {
					continue
				}
				field.Set(val)
			}
		}
	}

	return nil
}
