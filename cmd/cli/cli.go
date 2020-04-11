package main

import (
	"flag"
	"log"
	"sg/config"
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
		err = config.Generate("config/env/values.json", "config/generated/config.go")
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
	if genConfig {
		currentAction = generateConfig
	}

	flag.Parse()
}
