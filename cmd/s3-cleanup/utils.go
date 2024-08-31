package main

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func getFromEnv(env string) string {
	variable := os.Getenv(env)
	if variable == "" {
		log.Fatalf("Environment variable %s is not set", env)
	}
	return variable
}
