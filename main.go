package main

import (
	"fmt"
	"os"

	"github.com/erbalo/aws-lambda-connector/internal/client"
	"github.com/erbalo/aws-lambda-connector/internal/parser"
	"github.com/erbalo/aws-lambda-connector/internal/rpc"
)

func main() {
	conf, exit := handleArgs()
	if conf == nil {
		os.Exit(exit)
	}

	dialer := rpc.NewDialer()
	res, err := run(conf, dialer)
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(2)
	}

	println(res)
}

func handleArgs() (*parser.Configuration, int) {
	conf, err := parser.Parse(os.Args)
	if err != nil {
		os.Stderr.WriteString("Error parsing arguments: " + err.Error() + "\n")
		return nil, 1
	}

	if conf.ShowHelp {
		parser.ShowHelp()
		return nil, 0
	}

	return conf, 0
}

func run(conf *parser.Configuration, dialer rpc.Dialer) (string, error) {
	lambdaClient := client.NewLambda(*conf, dialer)

	res, err := lambdaClient.Invoke()
	if err != nil {
		return "", fmt.Errorf("error invoking Lambda: %w", err)
	}

	return string(res), nil
}
