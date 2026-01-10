# Zog, a zod inspired schema validation tool

Zog allows user defined schemas for runtime validation, it's also extensible to allow for custom validations and error messages.

## Usage

### Install

```bash
	go get "github.com/SebaPerezBerrios/zog"
```

### Import

```go
import (
	z "github.com/SebaPerezBerrios/zog"
)

```

### Define the base type

```go
type Customer struct {
	ID      int
	Email   string
	Phone   *string
	Address Address
	Account *Account
}

type Address struct {
	ZipCode  string
	City     string
	District string
	Address  string
}

type Account struct {
	ProviderID string
	AccountID  string
	Tokens     []Token
}

type Token struct {
	Data string
}

```

### Define the schema

```go
var CustomerSchema = z.Struct[Customer](
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
				"Data": z.String().Min(1),
			})),
		}).Optional(),
	},
)


```

### Validate

```go
	errors := z.Validate(&customer, CustomerSchema)

```

### Extensible API

You can define custom validation logic by creating new types that embed the core validation struct

```go
type EmailValidation struct {
	z.Validation
}
```

Then Create builder functions

```go
func EmailValidator() *EmailValidation {
	fns := [](func(string) error){func(email string) error {
		// your validation logic here
		return nil
	}}

	return &EmailValidation{
		Validation: *z.From(reflect.String, fns...),
	}
}

```

### Validation functions

Validations support custom functions that get computed only and after every sub validation has succeeded.

```go
z.Struct(z.StructSchema{
			"ZipCode":  z.String().Min(6),
			"City":     z.String().Min(1),
			"District": z.String().Min(1),
			"Address":  z.String().Min(1),
		}, addressValidation)
```

Given this schema. `addressValidation` will only recieve a valid address as defined by the struct schema.
