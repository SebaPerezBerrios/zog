package z

import (
	"fmt"
	"reflect"
	"strconv"
)

func validate[T any](value *T, validations ...ValidationI) []ValidationErrorItem {
	context := validationContext{}
	objectValue := reflect.ValueOf(value)
	if objectValue.Type().Kind() == reflect.Pointer && objectValue.IsNil() {
		context.Errors = append(context.Errors,
			ValidationErrorItem{
				Error: "nil object found at root\n",
				Path:  []string{},
				Kind:  TypeError,
			},
		)
	} else {
		validateRoot(objectValue.Elem(), validations, []string{}, &context)
	}
	return context.Errors
}

func validateRoot(objectValue reflect.Value, validations []ValidationI, path []string, context *validationContext) {
	if len(validations) == 0 {
		return
	}

	objectKind := objectValue.Type().Kind()

	switch objectKind {
	case reflect.Struct:
		validateObject(objectValue, validations, path, context)
	case reflect.Slice:
		validateSlice(objectValue, validations, path, context)
	default:
		for index := range validations {
			validateItemOrPointer(objectValue, objectKind, validations[index], path, context)
		}
	}
}

func validateObject(objectValue reflect.Value, validations []ValidationI, path []string, context *validationContext) {
	objectType := objectValue.Type()

	for index := range validations {
		fieldName := validations[index].get().key
		fieldType, found := objectType.FieldByName(fieldName)

		if !found {
			context.Errors = append(context.Errors,
				ValidationErrorItem{
					Error: fmt.Sprintf("validation field %v missing\n", validations[index].get().key),
					Path:  path,
					Kind:  FieldError,
				},
			)
		} else {
			fieldValue := objectValue.FieldByName(fieldName)
			nextPath := append(path, fieldType.Name)

			validateItemOrPointer(fieldValue, fieldType.Type.Kind(), validations[index], nextPath, context)
		}

	}
}

func validateSlice(objectValue reflect.Value, validations []ValidationI, path []string, context *validationContext) {
	for index := range objectValue.Len() {
		itemValue := objectValue.Index(index)

		nextPath := append(path, strconv.Itoa(index))

		for validationIndex := range validations {
			validateItemOrPointer(itemValue, itemValue.Type().Kind(), validations[validationIndex], nextPath, context)
		}

	}
}

func validateItemOrPointer(objectValue reflect.Value, kind reflect.Kind, validation ValidationI, path []string, context *validationContext) {
	isPointer := objectValue.Type().Kind() == reflect.Pointer

	if isPointer {
		validatePointer(objectValue, validation, path, context)
		return
	}

	validateChildren := validateItem(objectValue, kind, validation, path, context)

	if validateChildren {
		validateRoot(objectValue, validation.get().children, path, context)
	}
}

func validatePointer(objectValue reflect.Value, validation ValidationI, path []string, context *validationContext) {
	if !objectValue.IsNil() {
		validateItemOrPointer(objectValue.Elem(), objectValue.Elem().Kind(), validation, path, context)
	} else {
		if !validation.get().optional {
			context.Errors = append(context.Errors,
				ValidationErrorItem{
					Error: "expected non nil value\n",
					Path:  path,
					Kind:  NilError,
				},
			)
		}
	}
}

func validateItem(objectValue reflect.Value, kind reflect.Kind, validation ValidationI, path []string, context *validationContext) bool {
	if len(validation.get().validationFns) == 0 {
		return true
	}

	value := objectValue.Interface()

	if kind != validation.get().kind {
		context.Errors = append(context.Errors,
			ValidationErrorItem{
				Error: fmt.Sprintf("expected kind %v, found %v\n", validation.get().kind, kind),
				Path:  path,
				Kind:  TypeError,
			},
		)
		return false
	}

	for index := range validation.get().validationFns {
		err := validation.get().validationFns[index](value)
		if err != nil {
			context.Errors = append(context.Errors,
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
