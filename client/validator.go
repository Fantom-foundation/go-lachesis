package client

type ValidationContext struct {}

type ValidationStatus struct {}


// Validator is a Validator interface.
type Validator interface {
	*Member

	// GetAccount returns the underlying Account
	Validate(context *ValidationContext) ValidationStatus

}