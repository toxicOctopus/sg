package main

import (
	"flag"
	"log"

	"github.com/toxicOctopus/sg/config"
)

type action int

const (
	generateConfig = action(iota)
)

var currentAction action

func main() {
	var err error

	switch currentAction {
	case generateConfig:
		err = config.Generate(config.GetDefaultValuesPath(), config.BaseConfigFolder + "/" + config.GeneratedConfigFile)
		if err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal("no action provided")
	}
}

func init() {
	var genConfig bool
	flag.BoolVar(&genConfig, "generate-config", false, "(re)generate config code")
	flag.Parse()

	if genConfig {
		currentAction = generateConfig
	}
}
