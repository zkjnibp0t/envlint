package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/envlint/envparser"
	"github.com/yourorg/envlint/reporter"
	"github.com/yourorg/envlint/schema"
	"github.com/yourorg/envlint/validator"
)

func main() {
	envFile := flag.String("env", ".env", "Path to the .env file")
	schemaFile := flag.String("schema", ".env.schema.yaml", "Path to the schema YAML file")
	format := flag.String("format", "text", "Output format: text or json")
	flag.Parse()

	sch, err := schema.Load(*schemaFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading schema: %v\n", err)
		os.Exit(2)
	}

	env, err := envparser.Parse(*envFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing env file: %v\n", err)
		os.Exit(2)
	}

	validationErrors := validator.Validate(sch, env)

	r := reporter.New(os.Stdout, *format)
	if err := r.Write(validationErrors); err != nil {
		fmt.Fprintf(os.Stderr, "error writing report: %v\n", err)
		os.Exit(2)
	}

	if len(validationErrors) > 0 {
		os.Exit(1)
	}
}
