package handler

import (
	"errors"
	"gocourse/pkg/utils"
	"reflect"
	"strings"
)

func CheckBlankFields(value interface{}) error {
	val := reflect.ValueOf(value)
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if field.Kind() == reflect.String && field.String() == "" {
			// http.Error(w, "All fields are required", http.StatusBadRequest)
			return utils.ErrorHandler(errors.New("all fields is required"), "All fields is required")
		}
	}
	return nil
}

func GetFieldNames(model interface{}) []string {
	val := reflect.TypeOf(model)
	fields := []string{}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		// fields = append(fields, field.Tag.Get("json")) //get JSON tag
		fieldToAdd := strings.TrimSuffix(field.Tag.Get("json"), ",omitempty")
		fields = append(fields, fieldToAdd) //get JSON tag
	}
	return fields
}
