package logging

import (
	"github.com/buonotti/bus-stats-api/config/env"
	log "github.com/sirupsen/logrus"
)

func Setup() {
	if env.Env == env.Development {
		log.SetFormatter(&log.TextFormatter{
			ForceColors:  true,
			PadLevelText: true,
		})
	} else {
		log.SetFormatter(&log.JSONFormatter{})
	}
}
