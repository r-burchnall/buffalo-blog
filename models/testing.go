package models

import (
	"github.com/gobuffalo/validate/v3"
	"testing"
)

func assertNoValidationErrors(t *testing.T, got *validate.Errors) {
	t.Helper()

	if got.HasAny() {
		t.Errorf("Was expecting no validation errors but got %q", got)
	}
}

func assertError(t *testing.T, got error, want error) {
	t.Helper()

	if got != want {
		t.Errorf("was expecting an error %q but got %q", got, want)
	}
}
