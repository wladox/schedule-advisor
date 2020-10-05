package logger

import (
	log "github.com/sirupsen/logrus"
	"github.com/wladox/schedule-advisor/pkg/config"
)

// New create a new logger instance
func New(env *config.Config) (*log.Logger, error) {
	level, err := log.ParseLevel(env.LogLevel)
	if err != nil {
		return nil, err
	}

	l := log.New()
	//l.Formatter = new(logger.JSONFormatter)
	l.Infof("logger level: %s", level.String())
	l.Level = level

	return l, nil
}
