package utils

import (
	"strings"
	"testing"
)

func TestValidatePathParam_Valid(t *testing.T) {
	cases := []string{"hello", "World123", "a", "Z"}
	for _, c := range cases {
		if err := ValidatePathParam(c); err != nil {
			t.Errorf("expected valid for %q, got error: %v", c, err)
		}
	}
}

func TestValidatePathParam_TooLong(t *testing.T) {
	long := strings.Repeat("a", 101)
	if err := ValidatePathParam(long); err == nil {
		t.Error("expected error for 101-char path param, got nil")
	}
}

func TestValidatePathParam_Empty(t *testing.T) {
	if err := ValidatePathParam(""); err == nil {
		t.Error("expected error for empty path param, got nil")
	}
}

func TestValidatePathParam_InvalidChars(t *testing.T) {
	cases := []string{"hello world", "hello/world", "hello@world", "foo-bar", "foo_bar"}
	for _, c := range cases {
		if err := ValidatePathParam(c); err == nil {
			t.Errorf("expected error for %q (invalid chars), got nil", c)
		}
	}
}

func TestValidateQueryParam_Empty(t *testing.T) {
	if err := ValidateQueryParam(""); err != nil {
		t.Errorf("expected nil for empty query param (optional), got: %v", err)
	}
}

func TestValidateQueryParam_Valid(t *testing.T) {
	cases := []string{"hello", "foo123", "World"}
	for _, c := range cases {
		if err := ValidateQueryParam(c); err != nil {
			t.Errorf("expected valid for %q, got error: %v", c, err)
		}
	}
}

func TestValidateQueryParam_TooLong(t *testing.T) {
	long := strings.Repeat("a", 201)
	if err := ValidateQueryParam(long); err == nil {
		t.Error("expected error for 201-char query param, got nil")
	}
}
