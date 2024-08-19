package cmd

import (
	"fmt"
	"os"
)

func getFromEnv(env string) string {
	variable := os.Getenv(env)
	if variable == "" {
		panic(fmt.Sprintf("Environment variable %s is not set", env))
	}
	return variable
}
