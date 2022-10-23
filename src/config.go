package main

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type ConfigAction struct {
	Run string `yaml:"run"`
	Name string `yaml:"name"`
	If string `yaml:"if"`
	Env map[string]string `yaml:"env"`
}

type Config struct {
	Secret string `yaml:"secret"`
	Port string `yaml:"port"`
	Actions []ConfigAction `yaml:"actions"`
	Env map[string]string `yaml:"env"`
}

var configCache *Config

func readConfig() Config {
	if configCache != nil {
		return *configCache
	}

	var config, err = os.ReadFile("config.yml")

	if err != nil {
		log.Fatal("Cannot read config.yml. Use \"deployment init\" to create template")
	}

	var yamlConfig Config
	err = yaml.Unmarshal(config, &yamlConfig)

	if err != nil {
		log.Fatal("Unexpected syntax in config.yml. Use \"deployment init\" to create template")
	}

	configCache = &yamlConfig
	return yamlConfig
}