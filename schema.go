/*
Package z (zog) allows data validation in the form of composable schemas inspired by zod https://zod.dev/, using reflect https://pkg.go.dev/reflect.

	schema := z.Struct(
		z.StructSchema{
			"ID": z.Int().Min(0),
			"Title": z.String(func(title string) error {
				if strings.HasPrefix(title, " ") {
					return errors.New("title should not begin with empty spaces")
				}
				return nil
			}).Min(10).Max(200),
			"Authors": z.Slice[string]().Min(1),
			"Appendix": z.Struct(z.StructSchema{
				"Body": z.String().Min(1),
				"Annex": z.Struct(z.StructSchema{
					"AnnexBody": z.String().Max(100),
				}).Optional(),
			}),
		},
	)

	errors := z.Parse(&book, schema)

Using Gin https://gin-gonic.com/

		func (service *Service) testEndpoint(c *gin.Context) error {
			request := dtos.testEndpoint{}

			if err := c.ShouldBind(&request); err != nil {
				return err
			}

			errors := z.Parse(&request, testEndpointDto)
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

// Parse takes a struct pointer and a schema, returns the list of validation errors found, an empty list means successful validation
func Parse[T any](value *T, schema ValidationI) []ValidationErrorItem {
	if schema == nil {
		return []ValidationErrorItem{{
			Error: "nil schema reference",
			Path:  []string{},
			Kind:  NilError,
		}}
	}

	if schema.get().kind == reflect.Struct {
		return validate(value, schema.get().children...)
	}
	return validate(value, schema)
}
