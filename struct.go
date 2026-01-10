package z

import (
	"reflect"
)

type StructValidation[T any] struct {
	Validation
}

// Struct builds a struc validation schema, accepts a map (StructSchema) with key validation pairs.
func Struct[T any](schemaMap StructSchema, fns ...(func(T) error)) *StructValidation[T] {
	validations := make([]ValidationI, 0, len(schemaMap))

	for key, validation := range schemaMap {
		validationCopy := *validation.get()
		validationCopy.key = key
		validations = append(validations, &validationCopy)
	}

	return &StructValidation[T]{
		Validation{
			optional:      false,
			kind:          reflect.Struct,
			validationFns: wrapFns(fns),
			children:      validations,
		},
	}
}
