package main

import "strings"

// Messages are user facing strings. Do not expose app specific data.
type Message string

const (
	MessageInternalServerError Message = "something went wrong..."
	MessageInvalidCredentials  Message = "invalid credentials"
)

func (m Message) ToString() string {
	return string(m)
}

func (m Message) Capitalize() string {
	return strings.ToUpper(m.ToString())
}
