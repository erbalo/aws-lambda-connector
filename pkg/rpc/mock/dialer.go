package client

import (
	"errors"

	"github.com/aws/aws-lambda-go/lambda/messages"
	"github.com/erbalo/aws-lambda-connector/internal/rpc"
)

// MockRPCClient mocks an rpc.Client
type MockRPCClient struct {
	ResponsePayload []byte
	Err             error
}

func (m *MockRPCClient) Call(serviceMethod string, args interface{}, reply interface{}) error {
	if m.Err != nil {
		return m.Err
	}

	response := &messages.InvokeResponse{
		Payload: m.ResponsePayload,
	}

	*reply.(*messages.InvokeResponse) = *response
	return nil
}

// MockRPCDialer is a mock of the RPCDialer interface
type MockRPCDialer struct {
	MockClient *MockRPCClient
	ShouldFail bool
}

// Dial simulates a network dial operation
func (m *MockRPCDialer) Dial(network, address string) (rpc.Client, error) {
	if m.ShouldFail {
		return nil, errors.New("failed to dial")
	}

	return m.MockClient, nil
}
