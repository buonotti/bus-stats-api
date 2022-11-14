package logging

import (
	"github.com/buonotti/bus-stats-api/config"
	log "github.com/sirupsen/logrus"
)

func Setup() {
	if config.Env == config.Development {
		log.SetFormatter(&log.TextFormatter{
			ForceColors:  true,
			PadLevelText: true,
		})
	} else {
		log.SetFormatter(&log.JSONFormatter{})
	}
}
