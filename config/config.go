package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/ChimeraCoder/gojson"
	"github.com/peterbourgon/mergemap"
	"github.com/pkg/errors"
)

const (
	BaseConfigFolder = "config"
	GeneratedConfigFile = "generated_config.go"

	envSubFolder = "env"
	// contains all configs
	configPath = BaseConfigFolder + "/" + envSubFolder
	// this file must be present at configPath
	valuesConfigName = "values.json"

	defaultUpdateInterval = time.Second
)

type LiveConfig struct {
	cfg Config
	mu sync.RWMutex
}

func (c *LiveConfig) SetNew(new Config) {
	c.mu.Lock()
	c.cfg = new
	c.mu.Unlock()
}

func (c *LiveConfig) GetCfg() Config {
	c.mu.RLock()
	cp := c.cfg
	c.mu.RUnlock()

	return cp
}

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

// Read once read config files
func Read(env Env) (Config, error) {
	var err error
	cfg := Config{}
	defaultConfig := map[string]interface{}{}
	envConfig := map[string]interface{}{}

	err = readConfigFile(GetDefaultValuesPath(), &defaultConfig)
	if err != nil {
		return cfg, errors.Wrap(err, "default config")
	}
	err = readConfigFile(configPath + "/" + env.String() + "/" + valuesConfigName, &envConfig)
	if err == nil {
		defaultConfig = mergemap.Merge(defaultConfig, envConfig)
	}
	cfgJson, err := json.Marshal(defaultConfig)
	if err != nil {
		return cfg, errors.Wrap(err, "after config merge")
	}
	err = json.Unmarshal(cfgJson, &cfg)
	if err != nil {
		return cfg, errors.Wrap(err, "after merged config unmarshal")
	}

	return cfg, nil
}

func readConfigFile(path string, out *map[string]interface{}) error {
	cfg, err := ioutil.ReadFile(path)
	if err != nil {
		return errors.Wrap(err, "reading config file")
	}
	err = json.Unmarshal(cfg, &out)
	if err != nil {
		return errors.Wrap(err, "config file is not valid json")
	}

	return nil
}

// LiveRead config read blocking call
func LiveRead(env Env, cfg *LiveConfig, d time.Duration, errorCallback func(error)) {
	for {
		time.Sleep(d)

		newConfig, err := Read(env)
		if err != nil {
			errorCallback(err)
			continue
		}
		d = StringToUpdateInterval(newConfig.ConfigReadInterval)
		cfg.SetNew(newConfig)
	}
}

// StringToUpdateInterval convert with default value
func StringToUpdateInterval(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		d = defaultUpdateInterval
	}

	return d
}

func GetDefaultValuesPath() string {
	return configPath + "/" + valuesConfigName
}
