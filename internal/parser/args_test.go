package parser

import (
	"bytes"
	"io"
	"os"
	"reflect"
	"testing"
	"time"
)

type argsTestCase struct {
	name           string
	args           []string
	expectedConfig *Configuration
	expectErr      bool
}

func createTempFileWithContent(t *testing.T, content string) string {
	tmpfile, err := os.CreateTemp("", "example")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	return tmpfile.Name()
}

func TestParse(t *testing.T) {
	// Create a temporary file for testing event file reading
	eventFilePath := createTempFileWithContent(t, `{"event":"data"}`)
	defer os.Remove(eventFilePath) // Clean up the file afterwards

	tests := []argsTestCase{
		{
			name: "Default values",
			args: []string{"cmd"},
			expectedConfig: &Configuration{
				Address: "localhost:8080",
				Timeout: 30 * time.Second,
				Payload: []byte("{}"),
			},
			expectErr: false,
		},
		{
			name: "Custom address",
			args: []string{"cmd", "-a", "127.0.0.1:9000"},
			expectedConfig: &Configuration{
				Address: "127.0.0.1:9000",
				Timeout: 30 * time.Second,
				Payload: []byte("{}"),
			},
			expectErr: false,
		},
		{
			name: "Custom data",
			args: []string{"cmd", "-d", "{\"key\":\"value\"}"},
			expectedConfig: &Configuration{
				Address: "localhost:8080",
				Timeout: 30 * time.Second,
				Payload: []byte("{\"key\":\"value\"}"),
			},
			expectErr: false,
		},
		{
			name: "Custom timeout",
			args: []string{"cmd", "-t", "1m"},
			expectedConfig: &Configuration{
				Address: "localhost:8080",
				Timeout: 1 * time.Minute,
				Payload: []byte("{}"),
			},
			expectErr: false,
		},
		{
			name: "Show help",
			args: []string{"cmd", "-h"},
			expectedConfig: &Configuration{
				Address:  "localhost:8080",
				Timeout:  30 * time.Second,
				Payload:  []byte("{}"),
				ShowHelp: true,
			}, // When showing help, default configuration is expected
			expectErr: false,
		},
		{
			name: "Event file",
			args: []string{"cmd", "-e", eventFilePath},
			expectedConfig: &Configuration{
				Address: "localhost:8080",
				Timeout: 30 * time.Second,
				Payload: []byte(`{"event":"data"}`),
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := Parse(tt.args)
			if (err != nil) != tt.expectErr {
				t.Errorf("Parse() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if !reflect.DeepEqual(config, tt.expectedConfig) {
				t.Errorf("Parse() = %v, want %v", config, tt.expectedConfig)
			}
		})
	}
}

func TestShowHelp(t *testing.T) {
	expectedOutput := `Usage:
    aws-lambda-connector [flags]
        -e  path to the event JSON
        -d  data passed to the function, in JSON format, defaults to "{}"
        -a  the address of your local running function, defaults to localhost:8080
        -t  timeout for your handler execution, expressed as a duration, defaults to 30s
        -h, --help  show help
`

	// Temporarily redirect stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Call ShowHelp
	ShowHelp()

	// Close the write-end of the pipe
	w.Close()

	// Read the captured output
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Restore stdout
	os.Stdout = oldStdout

	// Assert
	if output != expectedOutput {
		t.Errorf("ShowHelp() output = %v, want %v", output, expectedOutput)
	}
}
