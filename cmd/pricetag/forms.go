package main

import (
	"bytes"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-playground/form/v4"
	"github.com/go-playground/validator/v10"
)

// HTML form validation errors. Keys are form input names
// and values are validation error messages.
type FormErrors map[string]string

// Implements error interface. This returns each key/value pair
// as its own line
func (formErrors FormErrors) Error() string {
	buff := bytes.NewBufferString("")

	for name, msg := range formErrors {
		buff.WriteString(name + ": " + msg)
		buff.WriteString("\n")
	}

	return strings.TrimSpace(buff.String())
}

const formErrorsSessionKey = "form-errors"

// Push form errors to session data
func (app *application) putFormErrors(r *http.Request, formErrors FormErrors) {
	app.sessionManager.Put(r.Context(), formErrorsSessionKey, formErrors)
}

// Pop form errors from session data
func (app *application) popFormErrors(r *http.Request) FormErrors {
	exists := app.sessionManager.Exists(r.Context(), formErrorsSessionKey)
	if exists {
		formErrors, ok := app.sessionManager.Pop(r.Context(), formErrorsSessionKey).(FormErrors)
		if ok {
			return formErrors
		}
	}

	return FormErrors{}
}

// Decode form values to struct and check for any form errors
func (app *application) parseForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = app.formDecoder.Decode(dst, r.Form)
	if err != nil {
		var invalidDecoderError *form.InvalidDecoderError
		switch {
		case errors.As(err, &invalidDecoderError):
			panic(err)
		default:
			return err
		}
	}

	err = app.validate.Struct(dst)
	if err != nil {
		var validationErrors validator.ValidationErrors
		switch {
		case errors.As(err, &validationErrors):
			formErrors := make(FormErrors)
			for _, fieldErr := range validationErrors {
				tag := fieldErr.Tag()
				param := fieldErr.Param()

				var msg string
				switch tag {
				case "email":
					msg = "invalid email"
				case "min":
					msg = "minimum length: " + param
				case "max":
					msg = "maximum length: " + param
				case "eqfield":
					// Field input value not equal to input[name="password"]
					if param == "Password" {
						msg = "passwords must be the same"
					}
				default:
					// This should be unexpected (akin to internal server error).
					// Return a generic error message and log the FieldError Tag
					// so that we can implement an error message.
					app.logger.Error(
						"unexpected form validation error",
						slog.String("tag", tag),
						slog.String("param", param))

					msg = ""
				}

				name := fieldErr.StructField()
				formErrors[name] = msg
			}
			return formErrors
		default:
			return err
		}
	}

	return nil
}
