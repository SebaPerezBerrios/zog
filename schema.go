/*
Package z (zog) allows data validation in the form of composable schemas inspired by zod https://zod.dev/, using reflect https://pkg.go.dev/reflect.

CustomerSchema := z.Struct[Customer](

	z.StructSchema{
		"ID":    z.Int().Min(0),
		"Email": z.String(),
		"Phone": z.String().Optional(),
		"Address": z.Struct(z.StructSchema{
			"ZipCode":  z.String().Min(6),
			"City":     z.String().Min(1),
			"District": z.String().Min(1),
			"Address":  z.String().Min(1),
		}, func(address Address) error {
			// further validate address
			return nil
		}),
		"Account": z.Struct[Account](z.StructSchema{
			"ProviderID": z.String().Min(1).Max(3),
			"AccountID":  z.String().Min(1),
			"Tokens": z.Slice[Token]().Min(1).ForEach(z.Struct[Token](z.StructSchema{
				"ID": z.String().Min(1),
			})),
		}).Optional(),
	},

)

	errors := z.Validate(&customer, CustomerSchema)

Using Gin https://gin-gonic.com/

		func (service *Service) testEndpoint(c *gin.Context) error {
			request := dtos.testEndpoint{}

			if err := c.ShouldBind(&request); err != nil {
				return err
			}

			errors := z.Validate(&request, testEndpointDto)
	}
*/
package z

import (
	"reflect"
)

func wrapFns[T any](fns [](func(T) error)) [](func(any) error) {
	result := make([](func(any) error), len(fns))
	for index := range fns {
		result[index] = func(x any) error {
			return fns[index](x.(T))
		}
	}
	return result
}

/*
From allows for user defined validations, using reflect for type validation, keep in mind that a mismatch between reflect and actual type conversion will result in runtime panic

	type CustomStringValidation struct {
		z.Validation
	}

	func CustomString() *CustomStringValidation {
		fns := [](func(string) error){func(s string) error {
			if s == "ilegal_value" {
				return errors.New("invalid string")
			}
			return nil
		}}

		return &CustomStringValidation{
			Validation: *z.From(reflect.String, fns...),
		}
	}
*/
func From[T any](reflectKind reflect.Kind, fns ...(func(T) error)) *Validation {
	return &Validation{
		optional:      false,
		kind:          reflectKind,
		validationFns: wrapFns(fns),
	}
}

// Optional allows for nil pointers to be a valid schema, this allows for optional values to be validated if they are present or skipped if they are a nil pointer
func (validation Validation) Optional() *Validation {
	validation.optional = true
	return &validation
}

// Required fails the validation on nil pointers (default), see Optional
func (validation Validation) Required() *Validation {
	validation.optional = false
	return &validation
}

// Validate takes a struct pointer and a schema, returns the list of validation errors found, an empty list means successful validation
func Validate[T any](value *T, schema ValidationI) []ValidationErrorItem {
	if schema == nil {
		return []ValidationErrorItem{{
			Error: "nil schema reference",
			Path:  []string{},
			Kind:  NilError,
		}}
	}

	return validateRoot(value, schema)
}
