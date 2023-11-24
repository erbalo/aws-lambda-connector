package main

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/erbalo/aws-lambda-connector/internal/parser"
	rpcMock "github.com/erbalo/aws-lambda-connector/pkg/rpc/mock"
	"github.com/stretchr/testify/assert"
)

type runTestCase struct {
	name           string
	args           []string
	mockRPCClient  *rpcMock.MockRPCClient
	shouldFailDial bool
	expectedResult string
	expectedError  bool
}

type argsTestCase struct {
	name          string
	args          []string
	expectedError bool
	expectedExit  int
}

func TestRun(t *testing.T) {
	testCases := []runTestCase{
		{
			name:           "Successful invocation",
			args:           []string{"cmd", "-a", "localhost:8080", "-t", "5s"},
			mockRPCClient:  &rpcMock.MockRPCClient{ResponsePayload: []byte("success")},
			shouldFailDial: false,
			expectedResult: "success",
			expectedError:  false,
		},
		{
			name:           "Failure to dial",
			args:           []string{"cmd", "--invalid-arg"},
			mockRPCClient:  &rpcMock.MockRPCClient{},
			shouldFailDial: true,
			expectedResult: "",
			expectedError:  true,
		},
		{
			name:           "Failure in lambda invocation",
			args:           []string{"cmd", "-a", "localhost:8080", "-t", "5s"},
			mockRPCClient:  &rpcMock.MockRPCClient{Err: errors.New("lambda error")},
			shouldFailDial: false,
			expectedResult: "",
			expectedError:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockDialer := &rpcMock.MockRPCDialer{MockClient: tc.mockRPCClient, ShouldFail: tc.shouldFailDial}

			os.Args = tc.args

			// Invoke the function that should write to stdout
			result, err := runWrapper(mockDialer)

			if tc.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, result)
			}
		})
	}
}

func runWrapper(dialer *rpcMock.MockRPCDialer) (string, error) {
	configuration, err := parser.Parse(os.Args)
	if err != nil {
		return "", fmt.Errorf("invalid arguments")
	}

	res, err := run(configuration, dialer)
	if err != nil {
		return "", fmt.Errorf("error invoking Lambda: %w", err)
	}

	return string(res), nil
}

func TestHandleArgs(t *testing.T) {
	testCases := []argsTestCase{
		{
			name:          "Valid args",
			args:          []string{"cmd", "-a", "localhost:8080", "-t", "5s"},
			expectedError: false,
			expectedExit:  0,
		},
		{
			name:          "Invalid args",
			args:          []string{"cmd", "--invalid-arg"},
			expectedError: true,
			expectedExit:  1,
		},
		{
			name:          "Show help",
			args:          []string{"cmd", "--help"},
			expectedError: false,
			expectedExit:  0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			os.Args = tc.args
			_, exitCode := handleArgs()

			if tc.expectedError {
				assert.NotEqual(t, 0, exitCode)
			} else {
				assert.Equal(t, tc.expectedExit, exitCode)
			}
		})
	}
}
