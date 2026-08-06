package z

import (
	"fmt"
	"reflect"
	"regexp"
)

type StringValidation struct {
	Validation
}

// String builds a string validation schema, accepts a list of validation functions as arguments
func String(fns ...(func(string) error)) *StringValidation {
	return &StringValidation{
		Validation: *From(reflect.String, fns...),
	}
}

// Min string length
func (stringValidation StringValidation) Min(stringLen int) *StringValidation {
	stringValidation.validationFns = append(stringValidation.validationFns, func(val any) error {
		valS := val.(string)
		if len(valS) < stringLen {
			return fmt.Errorf("string length %d less than %d", len(valS), stringLen)
		}
		return nil
	})
	return &stringValidation
}

// Max string length
func (stringValidation StringValidation) Max(stringLen int) *StringValidation {
	stringValidation.validationFns = append(stringValidation.validationFns, func(val any) error {
		valS := val.(string)
		if len(valS) > stringLen {
			return fmt.Errorf("string length %d greater than %d", len(valS), stringLen)
		}
		return nil
	})
	return &stringValidation
}

// Regex validation
func (stringValidation StringValidation) Regex(regex *regexp.Regexp) *StringValidation {
	stringValidation.validationFns = append(stringValidation.validationFns, func(val any) error {
		valS := val.(string)
		if regex.MatchString(valS) {
			return fmt.Errorf("regex match failed for string %s", valS)
		}
		return nil
	})
	return &stringValidation
}

// Default Value
func (stringValidation StringValidation) Default(defaultString string) *StringValidation {
	stringValidation.defaultValue = defaultString
	stringValidation.useDefault = true
	return &stringValidation
}
