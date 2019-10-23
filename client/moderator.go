package client

type ValidationContext struct {}

type ValidationStatus struct {}


// Moderator is a Moderator interface.
type Moderator interface {
	*Member

	// GetAccount returns the underlying Account
	PostValidate(context *ValidationContext) ValidationStatus


}