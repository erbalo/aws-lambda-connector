package main

import (
	"fmt"
	"os"

	"github.com/erbalo/aws-lambda-connector/internal/client"
	"github.com/erbalo/aws-lambda-connector/internal/parser"
	"github.com/erbalo/aws-lambda-connector/internal/rpc"
)

func main() {
	conf, err := parser.Parse(os.Args)
	if err != nil {
		os.Stderr.WriteString("Error parsing arguments: " + err.Error() + "\n")
		os.Exit(1)
	}

	if conf.ShowHelp {
		parser.ShowHelp()
		os.Exit(0)
	}

	dialer := rpc.NewDialer()
	res, err := run(conf, dialer)
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(2)
	}

	println(res)
}

func run(conf *parser.Configuration, dialer rpc.Dialer) (string, error) {
	lambdaClient := client.NewLambda(*conf, dialer)

	res, err := lambdaClient.Invoke()
	if err != nil {
		return "", fmt.Errorf("error invoking Lambda: %w", err)
	}

	return string(res), nil
}
