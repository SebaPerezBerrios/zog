package z

import (
	"fmt"
	"reflect"
)

type SliceValidation[T any] struct {
	Validation
}

func Slice[T any](fns ...(func([]T) error)) *SliceValidation[T] {
	return &SliceValidation[T]{
		Validation: *From(reflect.Slice, fns...),
	}
}

func (sliceValidation SliceValidation[T]) ForEach(subSchemas ...ValidationI) *SliceValidation[T] {
	sliceValidation.children = append(sliceValidation.children, subSchemas...)

	return &sliceValidation
}

func (sliceValidation SliceValidation[T]) Min(minLen int) *SliceValidation[T] {
	sliceValidation.validationFns = append(sliceValidation.validationFns, func(val any) error {
		valS := val.([]T)
		if len(valS) <= minLen {
			return fmt.Errorf("slice length %d less than %d", len(valS), minLen)
		}
		return nil
	})
	return &sliceValidation
}

func (sliceValidation SliceValidation[T]) Max(maxLen int) *SliceValidation[T] {
	sliceValidation.validationFns = append(sliceValidation.validationFns, func(val any) error {
		valS := val.([]T)
		if len(valS) >= maxLen {
			return fmt.Errorf("slice length %d greater than %d", len(valS), maxLen)
		}
		return nil
	})
	return &sliceValidation
}
