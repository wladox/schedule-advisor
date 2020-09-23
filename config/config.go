package config

import "github.com/kelseyhightower/envconfig"


type Config struct {
	FlinkMetricsAPIUrl          string  `envconfig:"FLINK_METRICS_API"`
}

// LoadFromENV loads config values from env
func (sc *Config) LoadFromENV() error {
	err := envconfig.Process("", sc)
	if err != nil {
		return err
	}

	return nil
}