package main

import (
	"flag"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"log"
	"os"
)

func getUser(dfl string) string {
	extValue := os.Getenv("PG_USER")
	if extValue == "" {
		return dfl
	} else {
		return extValue
	}
}

func MakeIncidentKey(recordType string) string {
	myHostName, _ := os.Hostname()
	return fmt.Sprintf("standby_%s_%s", recordType, myHostName)
}

func main() {
	var configFilePath string

	flag.StringVar(&configFilePath, "config", "pgsyncmon.yml", "Load configuration from this (yaml-formatted) file")
	flag.Parse()

	log.Print("pgsyncmon, Andy Gallagher 2019. See https://github.com/fredex42/pgsyncmon for source code and details.")

	log.Printf("Loading config from %s", configFilePath)
	cfg, loadConfigErr := ConfigFromFile(configFilePath)

	if loadConfigErr != nil {
		log.Fatalf("Could not read config file %s: %s", configFilePath, loadConfigErr)
	}

	if cfg.PostgresUser == "" {
		cfg.PostgresUser = "postgres"
	}

	result, checkErr := TestRecoveryStatus(getUser(cfg.PostgresUser), false)

	if checkErr != nil {
		log.Fatal("Could not check current status")
	} else {
		spew.Dump(result)
	}

	myHostName, _ := os.Hostname()
	lag := result.Lag()

	if result.IsInRecovery == false {
		log.Print("Standby database problem - not in recovery mode")
		sendErr := SendIncident(cfg, fmt.Sprintf("Standby database lost recovery on %s", myHostName), MakeIncidentKey("lost"), 0)
		if sendErr != nil {
			os.Exit(1)
		}
	} else if lag.Upper > 0 || lag.Lower > 100 {
		log.Print("Standby database problem - recovery is lagging")
		sendErr := SendIncident(cfg, fmt.Sprintf("Standby database has a significant recovery lag on %s", myHostName), MakeIncidentKey("lag"), 0)
		if sendErr != nil {
			os.Exit(1)
		}
	}
	log.Print("Standby database status is OK")
	os.Exit(0)
}
