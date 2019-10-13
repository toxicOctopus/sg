package main

import (
	"log"
	"os"
	"time"

	"github.com/jessevdk/go-flags"
)

var (
	BuildTime time.Time

	opts struct {
		Env     string `long:"env" description:"type of environment"`
	}
)

func main() {
	log.Printf("pid: %d\n", os.Getpid())
	log.Printf("build: %s", BuildTime.Format(time.RFC822))
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatalf("Arguments error: %s", err)
	}
}

func init() {
	BuildTime = time.Now()
}
