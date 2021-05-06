package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
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
	// this file must be present at configPath
	valuesConfigName = "values.json"
	configPermissions = 0644

	defaultUpdateInterval = time.Second
	minimalUpdateInterval = time.Millisecond * 10
)

var (
	// contains all configs
	configPath = filepath.Join(BaseConfigFolder, envSubFolder)
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

// Generate go struct from json config
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
	err = ioutil.WriteFile(to, output, configPermissions)

	return
}

// Read once read config files
func Read(env Env, valuesPath string) (Config, error) {
	var err error
	cfg := Config{}
	defaultConfig := map[string]interface{}{}
	envConfig := map[string]interface{}{}

	err = readConfigFile(valuesPath, &defaultConfig)
	if err != nil {
		return cfg, errors.Wrap(err, "default config")
	}

	err = readConfigFile(filepath.Join(configPath, env.String(), valuesConfigName), &envConfig)
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
func LiveRead(env Env, cfg *LiveConfig, errorCallback func(error)) {
	duration := getConfigInterval(cfg.GetCfg().ConfigReadInterval)
	for {
		time.Sleep(duration)

		newConfig, err := Read(env, GetDefaultValuesPath())
		if err != nil {
			errorCallback(err)
			continue
		}
		duration = getConfigInterval(newConfig.ConfigReadInterval)
		cfg.SetNew(newConfig)
	}
}

func getConfigInterval(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		d = defaultUpdateInterval
	}

	if d < minimalUpdateInterval {
		d = minimalUpdateInterval
	}

	return d
}

func GetDefaultValuesPath() string {
	return filepath.Join(configPath, valuesConfigName)
}
