package suggester_test

import (
	"testing"

	"github.com/user/envlint/schema"
	"github.com/user/envlint/suggester"
	"github.com/user/envlint/validator"
)

func makeSchema(vars map[string]schema.VarDef) *schema.Schema {
	return &schema.Schema{Vars: vars}
}

func TestSuggest_MissingWithTypo(t *testing.T) {
	s := makeSchema(map[string]schema.VarDef{
		"DATABASE_URL": {Required: true, Type: "url"},
	})
	errs := []validator.ValidationError{
		{Variable: "DATABASE_URL", Kind: validator.ErrMissing},
	}
	envKeys := []string{"DATABASE_UR", "PORT"}

	result := suggester.Suggest(errs, s, envKeys)

	if len(result) == 0 {
		t.Fatal("expected at least one suggestion")
	}
	if result[0].Variable != "DATABASE_URL" {
		t.Errorf("expected variable DATABASE_URL, got %s", result[0].Variable)
	}
	if result[0].Message == "" {
		t.Error("expected non-empty suggestion message")
	}
}

func TestSuggest_MissingNoTypo(t *testing.T) {
	s := makeSchema(map[string]schema.VarDef{
		"SECRET_KEY": {Required: true, Type: "string"},
	})
	errs := []validator.ValidationError{
		{Variable: "SECRET_KEY", Kind: validator.ErrMissing},
	}
	envKeys := []string{"PORT", "HOST"}

	result := suggester.Suggest(errs, s, envKeys)

	if len(result) == 0 {
		t.Fatal("expected a suggestion")
	}
	expected := "add SECRET_KEY to your .env file"
	if result[0].Message != expected {
		t.Errorf("expected %q, got %q", expected, result[0].Message)
	}
}

func TestSuggest_InvalidType(t *testing.T) {
	s := makeSchema(map[string]schema.VarDef{
		"PORT": {Required: true, Type: "int"},
	})
	errs := []validator.ValidationError{
		{Variable: "PORT", Kind: validator.ErrInvalidType, Message: "not an int"},
	}
	envKeys := []string{"PORT"}

	result := suggester.Suggest(errs, s, envKeys)

	if len(result) == 0 {
		t.Fatal("expected a suggestion")
	}
	if result[0].Variable != "PORT" {
		t.Errorf("expected PORT, got %s", result[0].Variable)
	}
}

func TestSuggest_NoErrors(t *testing.T) {
	s := makeSchema(map[string]schema.VarDef{})
	result := suggester.Suggest(nil, s, []string{})
	if len(result) != 0 {
		t.Errorf("expected no suggestions, got %d", len(result))
	}
}

// TestSuggest_MultipleErrors verifies that suggestions are returned for each
// validation error when multiple variables fail validation.
func TestSuggest_MultipleErrors(t *testing.T) {
	s := makeSchema(map[string]schema.VarDef{
		"DATABASE_URL": {Required: true, Type: "url"},
		"PORT":         {Required: true, Type: "int"},
	})
	errs := []validator.ValidationError{
		{Variable: "DATABASE_URL", Kind: validator.ErrMissing},
		{Variable: "PORT", Kind: validator.ErrInvalidType, Message: "not an int"},
	}
	envKeys := []string{"PORT"}

	result := suggester.Suggest(errs, s, envKeys)

	if len(result) != 2 {
		t.Fatalf("expected 2 suggestions, got %d", len(result))
	}
}
