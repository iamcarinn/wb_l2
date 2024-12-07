package main

import (
	"testing"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		hasError bool
	}{
		{
			name:     "Valid string with repetitions",
			input:    "a4bc2d5e",
			expected: "aaaabccddddde",
			hasError: false,
		},
		{
			name:     "Valid string without repetitions",
			input:    "abcd",
			expected: "abcd",
			hasError: false,
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
			hasError: false,
		},
		{
			name:     "String with only digits",
			input:    "45",
			expected: "",
			hasError: true,
		},
		{
			name:     "String with escape sequences",
			input:    `qwe\4\5`,
			expected: "qwe45",
			hasError: false,
		},
		{
			name:     "String with repeated escape sequences",
			input:    `qwe\45`,
			expected: "qwe44444",
			hasError: false,
		},
		{
			name:     "String with escaped backslash",
			input:    `qwe\\5`,
			expected: `qwe\\\\\`,
			hasError: false,
		},
		{
			name:     "Invalid string with non-alphanumeric characters",
			input:    "a2*b3",
			expected: "",
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := Unpack(tt.input)

			if tt.hasError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("did not expect error but got: %v", err)
				}
				if output != tt.expected {
					t.Errorf("expected %q but got %q", tt.expected, output)
				}
			}
		})
	}
}
