package main

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func runCommand(input string, args []string) (string, error) {
	// Create a command to execute the Go program with provided arguments
	cmd := exec.Command("go", append([]string{"run", "task.go"}, args...)...)

	// Set the input for the command
	cmd.Stdin = strings.NewReader(input)

	// Capture the output of the command
	var out bytes.Buffer
	cmd.Stdout = &out

	// Run the command and return the output
	err := cmd.Run()
	return out.String(), err
}

func TestCutUtility(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		args     []string
		expected string
	}{
		{
			name:     "Select fields 1 and 3",
			input:    "apple\tbanana\tcherry\n",
			args:     []string{"-f", "1,3"},
			expected: "apple\tcherry\n",
		},
		{
			name:     "Custom delimiter",
			input:    "a;b;c\nd;e;f\n",
			args:     []string{"-f", "2", "-d", ";"},
			expected: "b\ne\n",
		},
		{
			name:     "Skip lines without delimiter",
			input:    "no-delimiter\napple\tbanana\tcherry\n",
			args:     []string{"-f", "1", "-s"},
			expected: "apple\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := runCommand(tt.input, tt.args)
			if err != nil {
				t.Fatalf("Error running command: %v", err)
			}

			if output != tt.expected {
				t.Errorf("Expected output:\n%q\nGot:\n%q", tt.expected, output)
			}
		})
	}
}
