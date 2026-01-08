package z

import (
	"fmt"
	"strings"
)

func (validationErrorKind ValidationErrorKind) String() string {
	switch validationErrorKind {
	case TypeError:
		return "type cast"
	case ValidationError:
		return "validation"
	case FieldError:
		return "missing field"
	case NilError:
		return "nil field"
	default:
		return fmt.Sprintf("ValidationErrorKind(%d)", validationErrorKind)
	}
}

func (validationErrorItem ValidationErrorItem) String() string {
	if len(validationErrorItem.Path) == 0 {
		return fmt.Sprintf(
			"kind: %s, path: <ROOT>, error: %s",
			validationErrorItem.Kind,
			validationErrorItem.Error)
	}
	return fmt.Sprintf(
		"kind: %s, path: %s, error: %s",
		validationErrorItem.Kind,
		strings.Join(validationErrorItem.Path, "."),
		validationErrorItem.Error)
}
