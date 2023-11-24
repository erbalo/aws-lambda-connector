package rpc

import "net/rpc"

type Dialer interface {
	Dial(network, address string) (Client, error)
}

type Client interface {
	Call(serviceMethod string, args interface{}, reply interface{}) error
}

func NewDialer() Dialer {
	return &DefaultDialer{}
}

type DefaultDialer struct{}

func (d *DefaultDialer) Dial(network, address string) (Client, error) {
	return rpc.Dial(network, address)
}
