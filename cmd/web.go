package main

import (
	"fmt"
	"os"
	"time"

	"github.com/jessevdk/go-flags"
	"sg/utils"
)

type arguments struct {
	utils.Arguments
	Port string `short:"p" long:"port" description:"http port" default:"8080" optional:"y"`
}

var (
	args      arguments
	startTime time.Time
)

func main() {
	fmt.Println(startTime)
}

func init() {
	_, err := flags.Parse(&args)
	if nil != err {
		os.Exit(1)
	}
	startTime = time.Now()
}
