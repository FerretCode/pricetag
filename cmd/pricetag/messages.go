package main

// These are user facing messages (typically for errors).
// Do not expose app specific data.
type Message string

const (
	MessageInternalServerError  Message = "something went wrong..."
	MessageInvalidCredentials   Message = "invalid credentials"
	MessageAlreadyAuthenticated Message = "already authenticated"
	MessageNotFound             Message = "how did I get here?"
)
