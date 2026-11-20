package dto

// VoteRequest is the validated query of POST /docs/{slug}/vote.
type VoteRequest struct {
	Vote string `validate:"required,oneof=yes no"`
}

// CommandRequest is the validated form of the project command builder on the home page.
type CommandRequest struct {
	Name     string `validate:"required,max=40"`
	Database string `validate:"required,oneof=sqlite postgres mysql"`
}
