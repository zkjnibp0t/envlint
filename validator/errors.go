package validator

// ErrorKind categorises the type of validation failure.
type ErrorKind string

const (
	// ErrMissing indicates a required variable is absent from the .env file.
	ErrMissing ErrorKind = "missing"

	// ErrInvalidType indicates the variable value does not match the expected type.
	ErrInvalidType ErrorKind = "invalid_type"

	// ErrInvalidFormat indicates the variable value does not match a required format.
	ErrInvalidFormat ErrorKind = "invalid_format"
)

// ValidationError describes a single validation failure for a variable.
type ValidationError struct {
	// Variable is the name of the .env variable that failed validation.
	Variable string `json:"variable"`

	// Kind classifies the nature of the error.
	Kind ErrorKind `json:"kind"`

	// Message is a human-readable description of the error.
	Message string `json:"message"`
}

// Error implements the error interface.
func (e ValidationError) Error() string {
	return e.Variable + ": " + e.Message
}
