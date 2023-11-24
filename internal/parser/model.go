package parser

import "time"

type Configuration struct {
	Address  string
	Payload  []byte
	Timeout  time.Duration
	ShowHelp bool
}
