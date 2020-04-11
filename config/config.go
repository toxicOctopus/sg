package config

import (
	"github.com/ChimeraCoder/gojson"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
)

// Generate from json config to go struct
func Generate(from, to string) (err error) {
	if len(from) == 0 {
		return errors.New("no source specified")
	}

	inputFile, err := os.Open(from)
	if err != nil {
		return errors.Wrap(err, "failed to read from source")
	}
	defer func() {
		closeErr := inputFile.Close()
		if err == nil {
			err = closeErr
		}
	}()

	output, err := gojson.Generate(inputFile, gojson.ParseJson, "Config", "config", []string{"json"}, false, true)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(to, output, 0644)

	return
}

// LiveRead config read blocking call
func LiveRead() {

}
