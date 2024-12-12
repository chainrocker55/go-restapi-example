package utils

import (
	"testing"
)

// TestIf tests the If function
func TestIf(t *testing.T) {
	trueValue := "trueValue"
	falseValue := "falseValue"
	tests := []struct {
		cond    bool
		vtrue   string
		vfalse  string
		want    string
	}{
		{true, trueValue, falseValue, trueValue},
		{false, trueValue, falseValue, falseValue},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := If(tt.cond, tt.vtrue, tt.vfalse)
			if got != tt.want {
				t.Errorf("If(%v, %v, %v) = %v; want %v", tt.cond, tt.vtrue, tt.vfalse, got, tt.want)
			}
		})
	}
}

// TestIfNil tests the IfNil function
func TestIfNil(t *testing.T) {
	defaultValue := "defaultValue"
	value := "someValue"

	tests := []struct {
		input *string
		def   string
		want  string
	}{
		{nil, defaultValue, defaultValue},
		{&value, defaultValue, value},
	}

	for _, tt := range tests {
		t.Run(tt.def, func(t *testing.T) {
			got := IfNil(tt.input, tt.def)
			if got != tt.want {
				t.Errorf("IfNil(%v, %v) = %v; want %v", tt.input, tt.def, got, tt.want)
			}
		})
	}
}

// TestDefault tests the Default function
func TestDefault(t *testing.T) {
	value := "someValue"
	tests := []struct {
		input *string
		def   string
		want  string
	}{
		{nil, "defaultValue", "defaultValue"},
		{&value, "defaultValue", value},
	}

	for _, tt := range tests {
		t.Run(tt.def, func(t *testing.T) {
			got := Default(tt.input, tt.def)
			if got != tt.want {
				t.Errorf("Default(%v, %v) = %v; want %v", tt.input, tt.def, got, tt.want)
			}
		})
	}
}
