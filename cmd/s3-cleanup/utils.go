package cmd

import (
	"fmt"
	"os"
)

func getFromEnv(env string) string {
	home := os.Getenv(env)
	if home == "" {
		panic(fmt.Sprintf("Environment variable %s is not set", env))
	}
	return home
}