package main

import (
	"os"

	"github.com/erbalo/aws-lambda-connector/internal/client"
	"github.com/erbalo/aws-lambda-connector/internal/parser"
)

func main() {
	configuration, err := parser.Parse(os.Args)
	if err != nil {
		os.Stderr.WriteString("Error parsing arguments: " + err.Error() + "\n")
		os.Exit(1)
	}

	if configuration.ShowHelp {
		parser.ShowHelp()
		os.Exit(0)
	}

	lambdaClient := client.New(*configuration)

	res, err := lambdaClient.Invoke()
	if err != nil {
		os.Stderr.WriteString("Error invoking Lambda: " + err.Error() + "\n")
		os.Exit(2)
	}

	println(string(res))
}
