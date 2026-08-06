package z

import "reflect"

// StructSchema is the type for key value pairs used to validate structs
type StructSchema = map[string]ValidationI

// ValidationI is the interface which all sub validators should implement
type ValidationI interface {
	get() *Validation
}

type Validation struct {
	key           string
	kind          reflect.Kind
	optional      bool
	defaultValue  any
	useDefault    bool
	validationFns [](func(any) error)
	children      []ValidationI
}

func (validation *Validation) get() *Validation {
	return validation
}

type ValidationErrorKind int

const (
	TypeError ValidationErrorKind = iota
	ValidationError
	FieldError
	NilError
)

// ValidationErrorItem gives information about validation errors, it implements json marshall (json tags)
type ValidationErrorItem struct {
	Error string              `json:"error"`
	Path  []string            `json:"path"`
	Kind  ValidationErrorKind `json:"kind"`
}

type validationContext struct {
	Errors []ValidationErrorItem `json:"errors"`
}
