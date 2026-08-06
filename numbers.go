package z

import (
	"fmt"
	"reflect"
)

type IntValidation struct {
	Validation
}

type FloatValidation struct {
	Validation
}

func Int(fns ...(func(int) error)) *IntValidation {
	return &IntValidation{
		Validation: *From(reflect.Int, fns...),
	}
}

func Float(fns ...(func(float64) error)) *FloatValidation {
	return &FloatValidation{
		Validation: *From(reflect.Float64, fns...),
	}
}

func (intValidation IntValidation) Min(minValue int) *IntValidation {
	intValidation.validationFns = append(intValidation.validationFns, func(val any) error {
		valI := val.(int)
		if (valI) <= minValue {
			return fmt.Errorf("value %d less than %d", valI, minValue)
		}
		return nil
	})

	return &intValidation
}

func (intValidation IntValidation) Max(maxValue int) *IntValidation {
	intValidation.validationFns = append(intValidation.validationFns, func(val any) error {
		valI := val.(int)
		if (valI) <= maxValue {
			return fmt.Errorf("value %d greater than %d", valI, maxValue)
		}
		return nil
	})
	return &intValidation
}

// Default Value
func (intValidation IntValidation) Default(defaultInt int) *IntValidation {
	intValidation.defaultValue = defaultInt
	intValidation.useDefault = true
	return &intValidation
}
