package main

import (
	"github.com/davecgh/go-spew/spew"
	"log"
	"os"
)

func getUser() string {
	extValue := os.Getenv("PG_USER")
	if extValue == "" {
		return "postgres"
	} else {
		return extValue
	}
}

func main() {
	log.Print("pgsyncmon, Andy Gallagher 2019. See https://github.com/fredex42/pgsyncmon for source code and details.")

	result, checkErr := TestRecoveryStatus(getUser(), false)

	if checkErr != nil {
		log.Fatal("Could not check current status")
	} else {
		spew.Dump(result)
	}
	os.Exit(0)
}
