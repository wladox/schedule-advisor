package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	FlinkMetricsAPIUrl   string `envconfig:"FLINK_METRICS_API" default:"http://localhost:8081"`
	PredictionServiceUrl string `envconfig:"PREDICTION_API" default:"http://localhost:5000"`
	Interval             int    `envconfig:"POLL_INTERVAL" default:"5000"`
	Port                 int    `envconfig:"PORT" default:"8000"`
	LogLevel             string `envconfig:"LOG_LEVEL"`
}

// LoadFromENV loads config values from env
func (sc *Config) LoadFromENV() error {
	err := envconfig.Process("", sc)
	if err != nil {
		return err
	}

	return nil
}
