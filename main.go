package main

import (
	_ "embed"
	"flag"
	"log"
	"os"

	"github.com/aprx/around/internal/config"
	"github.com/aprx/around/internal/generator"
)

//go:embed fixtures/fixtures.sh
var fixtures string

const (
	defaultConfigPath string = "$HOME/.config/around"
	defaultOutputPath string = "$HOME/.aroundrc"
)

var clOptions struct {
	configPath string
	output     string
}

func init() {
	flag.StringVar(&clOptions.configPath, "config", os.ExpandEnv(defaultConfigPath), "Path to configuration")
	flag.StringVar(&clOptions.output, "output", os.ExpandEnv(defaultOutputPath), "Output path")
}

func main() {
	flag.Parse()
	configurations := []config.Configuration{}
	configFiles := config.CollectConfigurations(clOptions.configPath)
	for _, name := range configFiles {
		conf, err := config.ParseFile(name, config.ParseOptions{})
		if err != nil {
			log.Fatal(err)
		}
		configurations = append(configurations, conf...)
	}
	if err := config.Validate(configurations); err != nil {
		log.Fatal(err)
	}
	output, err := os.OpenFile(clOptions.output, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o640)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = output.Close() }()

	_, _ = output.WriteString(fixtures)
	if err := generator.Generate(configurations, output); err != nil {
		log.Fatal(err)
	}
}
