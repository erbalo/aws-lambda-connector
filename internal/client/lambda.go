package client

import (
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/lambda/messages"
	"github.com/erbalo/aws-lambda-connector/internal/parser"
	"github.com/erbalo/aws-lambda-connector/internal/rpc"
)

type LambdaClient struct {
	configuration parser.Configuration
	dialer        rpc.Dialer
}

func NewLambda(configuration parser.Configuration, dialer rpc.Dialer) *LambdaClient {
	return &LambdaClient{
		configuration,
		dialer,
	}
}

func (client LambdaClient) Invoke() ([]byte, error) {
	deadline := time.Now().Add(client.configuration.Timeout)
	request := messages.InvokeRequest{Payload: client.configuration.Payload, Deadline: messages.InvokeRequest_Timestamp{
		Seconds: deadline.Unix(),
		Nanos:   int64(deadline.Nanosecond()),
	}}

	invoker, err := client.dialer.Dial("tcp", client.configuration.Address)
	if err != nil {
		return nil, err
	}

	var response messages.InvokeResponse
	err = invoker.Call("Function.Invoke", request, &response)
	if err != nil {
		return nil, err
	}

	if response.Error != nil {
		return nil, fmt.Errorf("lambda returned error:\n%s", response.Error.Message)
	}

	return response.Payload, nil
}
