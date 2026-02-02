package fme

import (
	"errors"
)

// Sentinel errors for explicit semantics.
var (
	ErrSchemaCycle         = errors.New("argschema: dependency cycle")
	ErrSchemaContradiction = errors.New("argschema: schema contradiction between dependency and interference")
	ErrCombinationRequire  = errors.New("argschema: requirement not respected")
	ErrCombinationInterfer = errors.New("argschema: interfering arguments in combination")
)

// SchemaValidationError wraps a schema-level validation failure with:
//   - a well-known Kind (one of ErrSchemaCycle / ErrSchemaContradiction);
//   - a human-readable Detail string that can go to logs or user-facing errors.
type SchemaValidationError struct {
	Kind   error  // ErrSchemaCycle or ErrSchemaContradiction
	Detail string // human-readable explanation
}

func (e *SchemaValidationError) Error() string { return e.Detail }
func (e *SchemaValidationError) Unwrap() error { return e.Kind }


type CombinationVerificationError struct {
	Kind   error
	Detail string
	Path *PathInstance
}
func (e *CombinationVerificationError) Error() string { return e.Detail }
func (e *CombinationVerificationError) Unwrap() error { return e.Kind }
