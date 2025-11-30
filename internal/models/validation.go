// Package models contains the core data structures and validation logic for the application.
package models

import (
	"time"
	"unicode"

	valid "github.com/asaskevich/govalidator"
)

// init initializes custom validation rules for the application using the govalidator package.
//
// Custom Validation Rules:
//   - "born": Validates that a given date is in the past.
//   - "passwordcheck": Validates that a password contains at least one lowercase letter, one uppercase letter, and one digit.
//
// Behavior:
//   - The "born" rule ensures that the provided date is before the current time.
//   - The "passwordcheck" rule iterates through the password string to check for the presence of required character types.
func init() {
	valid.CustomTypeTagMap.Set("born", func(date interface{}, context interface{}) bool {
		d, ok := date.(Date)
		if !ok {
			return false
		}

		return d.Time.Before(time.Now())
	})

	valid.TagMap["passwordcheck"] = valid.Validator(func(password string) bool {
		hasLowLetters := false
		hasUpperLetters := false
		hasDigits := false

		for _, c := range password {
			switch {
			case unicode.IsNumber(c):
				hasDigits = true
			case unicode.IsUpper(c):
				hasUpperLetters = true
			case unicode.IsLower(c):
				hasLowLetters = true
			}
		}

		return hasLowLetters && hasUpperLetters && hasDigits
	})
}
