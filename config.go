package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Config struct {
	PagerdutyKey string `yaml:"pagerDutyKey"`
	PostgresUser string `yaml:"postgresUser"`
	AllowedDrift int    `yaml:"allowedDrift"`
	ServiceId    string `yaml:"pagerDutyService"`
	Urgency      string `yaml:"pagerDutyUrgency"`
	AlertTitle   string `yaml:"pagerDutyAlertTitle"`
}

func ConfigFromFile(fileName string) (*Config, error) {
	content, openErr := ioutil.ReadFile(fileName)
	if openErr != nil {
		log.Printf("Could not open config file %s: %s", fileName, openErr)
		return nil, openErr
	}

	var rtn Config

	readErr := yaml.Unmarshal(content, &rtn)
	if readErr != nil {
		log.Printf("Could not read data from file: %s", readErr)
		return nil, readErr
	}

	return &rtn, nil
}
