package client

import (
	"errors"
	"testing"
	"time"

	"github.com/erbalo/aws-lambda-connector/internal/parser"
	rpcMock "github.com/erbalo/aws-lambda-connector/pkg/rpc/mock"
	"github.com/stretchr/testify/assert"
)

type testCase struct {
	name           string
	mockRPCClient  *rpcMock.MockRPCClient
	config         parser.Configuration
	expectedResult []byte
	expectedError  bool
}

func TestLambdaClient_Invoke(t *testing.T) {
	testCases := []testCase{
		{
			name:          "SuccessWithPayload",
			mockRPCClient: &rpcMock.MockRPCClient{ResponsePayload: []byte("payload")},
			config: parser.Configuration{
				Address: "test-address",
				Payload: []byte("test-payload"),
				Timeout: 1 * time.Minute,
			},
			expectedResult: []byte("payload"),
			expectedError:  false,
		},
		{
			name:          "FailureToDial",
			mockRPCClient: &rpcMock.MockRPCClient{Err: errors.New("failed to dial")},
			config: parser.Configuration{
				Address: "test-address",
				Payload: []byte("test-payload"),
				Timeout: 1 * time.Minute,
			},
			expectedResult: nil,
			expectedError:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockDialer := &rpcMock.MockRPCDialer{MockClient: tc.mockRPCClient, ShouldFail: tc.mockRPCClient.Err != nil}
			lambdaClient := NewLambda(tc.config, mockDialer)
			response, err := lambdaClient.Invoke()

			if tc.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, response)
			}
		})
	}
}
