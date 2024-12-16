package main

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Messages are user facing strings. Do not expose app specific data.
type Message string

const (
	MessageInternalServerError Message = "something went wrong..."
	MessageInvalidCredentials  Message = "invalid credentials"
	MessageNotFound            Message = "how did I get here?"
)

func (m Message) ToString() string {
	return string(m)
}

func (m Message) Capitalize() string {
	s := m.ToString()
	if len(s) == 0 {
		return s
	}

	// TODO: accept languag tag as argument and detect user's
	// language via request headers (default to English)
	caser := cases.Title(language.English)
	cap := caser.String(s[:1])
	return cap + s[1:]
}
