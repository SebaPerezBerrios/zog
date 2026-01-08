package z

import (
	"reflect"
)

type StructValidation struct {
	Validation
}

// Struct builds a struc validation schema, accepts a map (StructSchema) with key validation pairs.
func Struct(schemaMap StructSchema) *StructValidation {
	validations := make([]ValidationI, 0, len(schemaMap))

	for key, validation := range schemaMap {
		validationCopy := *validation.get()
		validationCopy.key = key
		validations = append(validations, &validationCopy)
	}

	return &StructValidation{
		Validation: Validation{
			optional: false,
			kind:     reflect.Struct,
			children: validations,
		},
	}
}

// SubSchema allows to compose sub struct schemas
func (structValidation StructValidation) SubSchema(subSchemas ...ValidationI) *StructValidation {
	structValidation.children = append(structValidation.children, subSchemas...)

	return &structValidation
}
