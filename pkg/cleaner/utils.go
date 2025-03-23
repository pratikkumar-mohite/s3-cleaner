package cleaner

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func getFromEnv(env string) string {
	variable := os.Getenv(env)
	if variable == "" {
		log.Errorf("Environment variable %s is not set", env)
	}
	return variable
}
