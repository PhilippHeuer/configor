package configor

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"io/ioutil"
	"errors"
	"strings"
	yaml "gopkg.in/yaml.v2"
)

type Configor struct {
	*Config
}

type Config struct {
	Environment string
	ENVPrefix   string
	Debug       bool
	Verbose     bool
}

// New initialize a Configor
func New(config *Config) *Configor {
	if config == nil {
		config = &Config{}
	}

	if os.Getenv("CONFIGOR_DEBUG_MODE") != "" {
		config.Debug = true
	}

	if os.Getenv("CONFIGOR_VERBOSE_MODE") != "" {
		config.Verbose = true
	}

	return &Configor{Config: config}
}

// GetEnvironment get environment
func (configor *Configor) GetEnvironment() string {
	if configor.Environment == "" {
		if env := os.Getenv("CONFIGOR_ENV"); env != "" {
			return env
		}

		if isTest, _ := regexp.MatchString("/_test/", os.Args[0]); isTest {
			return "test"
		}

		return "development"
	}
	return configor.Environment
}

// Load will unmarshal configurations to struct from files that you provide
func (configor *Configor) Load(config interface{}, files ...string) error {
	defer func() {
		if configor.Config.Debug || configor.Config.Verbose {
			fmt.Printf("Configuration:\n  %#v\n", config)
		}
	}()

	for _, file := range configor.getConfigurationFiles(files...) {
		if configor.Config.Debug || configor.Config.Verbose {
			fmt.Printf("Loading configurations from file '%v'...\n", file)
		}
		if err := processFile(config, file); err != nil {
			return err
		}
	}

	prefix := configor.getENVPrefix(config)
	if prefix == "-" {
		return configor.processTags(config)
	}
	return configor.processTags(config, prefix)
}

// Save will save the configurations to the provided filename
func Save(config interface{}, filename string) error {
	var fileContent []byte
	var err error

	switch {
	  case strings.HasSuffix(filename, ".yaml") || strings.HasSuffix(filename, ".yml"):
		  fileContent, err = yaml.Marshal(&config)
	  case strings.HasSuffix(filename, ".json"):
		  fileContent, err = json.Marshal(&config)
	  default:
		  return errors.New("Unknown file type")
	}

	if err != nil {
		return nil
	}

	err = ioutil.WriteFile(filename, fileContent, 0600)
	return err
}

// ENV return environment
func ENV() string {
	return New(nil).GetEnvironment()
}

// Load will unmarshal configurations to struct from files that you provide
func Load(config interface{}, files ...string) error {
	return New(nil).Load(config, files...)
}
