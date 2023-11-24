package parser

import (
	"flag"
	"fmt"
	"os"
	"time"
)

func Parse(args []string) (*Configuration, error) {
	config := &Configuration{}

	fs := flag.NewFlagSet("aws-lambda-connector", flag.ContinueOnError)
	fs.StringVar(&config.Address, "a", "localhost:8080", "address of the lambda function")

	var eventPath string
	fs.StringVar(&eventPath, "e", "", "path to the event JSON file")

	var data string
	fs.StringVar(&data, "d", "{}", "data in JSON format")
	fs.DurationVar(&config.Timeout, "t", 30*time.Second, "timeout handler execution")

	fs.BoolVar(&config.ShowHelp, "h", false, "show help")
	fs.BoolVar(&config.ShowHelp, "help", false, "show help")

	err := fs.Parse(args[1:])
	if err != nil {
		return nil, err
	}

	if eventPath != "" {
		fileContent, err := os.ReadFile(eventPath)
		if err != nil {
			return nil, fmt.Errorf("error reading event file: %w", err)
		}
		config.Payload = fileContent
	} else {
		config.Payload = []byte(data)
	}

	return config, nil
}

func ShowHelp() {
	fmt.Println(`Usage:
    aws-lambda-connector [flags]
        -e  path to the event JSON
        -d  data passed to the function, in JSON format, defaults to "{}"
        -a  the address of your local running function, defaults to localhost:8080
        -t  timeout for your handler execution, expressed as a duration, defaults to 30s
        -h, --help  show help`)
}
