package z

import (
	"fmt"
	"reflect"
	"strconv"
)

func validateRoot[T any](value *T, validations ...ValidationI) []ValidationErrorItem {
	context := validationContext{}
	objectValue := reflect.ValueOf(value)
	if objectValue.Type().Kind() == reflect.Pointer && objectValue.IsNil() {
		context.Errors = append(
			context.Errors,
			ValidationErrorItem{
				Error: "nil object found at root\n",
				Path:  []string{},
				Kind:  NilError,
			},
		)
	} else {
		validateAll(objectValue.Elem(), validations, []string{}, &context)
	}
	return context.Errors
}

func validateAll(objectValue reflect.Value, validations []ValidationI, path []string, context *validationContext) bool {
	status := true
	if len(validations) == 0 {
		return status
	}

	objectKind := objectValue.Type().Kind()

	switch objectKind {
	case reflect.Struct:
		return validateStruct(objectValue, validations, path, context)
	case reflect.Slice:
		return validateSlice(objectValue, validations, path, context)
	default:
		for index := range validations {
			subStatusOk := validateItemOrPointer(objectValue, objectKind, validations[index], path, context)
			if !subStatusOk {
				status = false
			}
		}
	}
	return status
}

func validateStruct(objectValue reflect.Value, validations []ValidationI, path []string, context *validationContext) bool {
	objectType := objectValue.Type()
	status := true

	for index := range validations {
		fieldName := validations[index].get().key

		// Root validation
		if fieldName == "" {
			return validateItemOrPointer(objectValue, objectValue.Type().Kind(), validations[index], path, context)
		}

		fieldType, found := objectType.FieldByName(fieldName)

		if !found {
			context.Errors = append(
				context.Errors,
				ValidationErrorItem{
					Error: fmt.Sprintf("validation field %v missing\n", validations[index].get().key),
					Path:  path,
					Kind:  FieldError,
				},
			)
			return false
		} else {
			fieldValue := objectValue.FieldByName(fieldName)
			nextPath := append(path, fieldType.Name)

			subStatusOk := validateItemOrPointer(fieldValue, fieldType.Type.Kind(), validations[index], nextPath, context)
			if !subStatusOk {
				status = false
			}
		}

	}
	return status
}

func validateSlice(objectValue reflect.Value, validations []ValidationI, path []string, context *validationContext) bool {
	status := true

	for index := range objectValue.Len() {
		itemValue := objectValue.Index(index)
		nextPath := append(path, strconv.Itoa(index))

		for validationIndex := range validations {
			subStatusOk := validateItemOrPointer(itemValue, itemValue.Type().Kind(), validations[validationIndex], nextPath, context)
			if !subStatusOk {
				status = false
			}
		}
	}
	return status
}

func validateItemOrPointer(objectValue reflect.Value, kind reflect.Kind, validation ValidationI, path []string, context *validationContext) bool {
	isPointer := objectValue.Type().Kind() == reflect.Pointer

	if isPointer {
		return validatePointer(objectValue, validation, path, context)
	}

	status := validateAll(objectValue, validation.get().children, path, context)

	if status {
		return validateItem(objectValue, kind, validation, path, context)
	}
	return status
}

func validatePointer(objectValue reflect.Value, validation ValidationI, path []string, context *validationContext) bool {
	if !objectValue.IsNil() {
		return validateItemOrPointer(objectValue.Elem(), objectValue.Elem().Kind(), validation, path, context)
	} else {
		if !validation.get().optional {
			context.Errors = append(
				context.Errors,
				ValidationErrorItem{
					Error: "expected non nil value\n",
					Path:  path,
					Kind:  NilError,
				},
			)
		}
		return false
	}
}

func validateItem(objectValue reflect.Value, kind reflect.Kind, validation ValidationI, path []string, context *validationContext) bool {
	value := objectValue.Interface()

	if kind != validation.get().kind {
		context.Errors = append(
			context.Errors,
			ValidationErrorItem{
				Error: fmt.Sprintf("expected kind %v, found %v\n", validation.get().kind, kind),
				Path:  path,
				Kind:  TypeError,
			},
		)
		return false
	}

	if validation.get().useDefault && objectValue.IsZero() {
		objectValue.Set(reflect.ValueOf(validation.get().defaultValue))
		return true
	}

	for index := range validation.get().validationFns {
		err := validation.get().validationFns[index](value)
		if err != nil {
			context.Errors = append(
				context.Errors,
				ValidationErrorItem{
					Error: fmt.Sprintln(err),
					Path:  path,
					Kind:  ValidationError,
				},
			)
			return false
		}
	}
	return true
}
