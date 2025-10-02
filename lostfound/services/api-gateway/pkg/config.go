package pkg

import (
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type (
	logging struct {
		Level   int    `mapstructure:"level,omitempty"`
		Path    string `mapstructure:"path,omitempty"`
		MaxSize int    `mapstructure:"max_size_mb,omitempty"`
	}
	configuration struct {
		Version      string  `mapstructure:"version,omitempty"`
		Logging      logging `mapstructure:"logging"`
		ServerPort   string  `mapstructure:"server_port"`
		IsProduction bool    `mapstructure:"is_production,omitempty"`
	}
)

var (
	AppConfig configuration
)

func ReadConfiguration() error {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("../config")
	viper.AddConfigPath("../../config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Error().Msgf("failed to read configuration file with error: %+v", err)
		return fmt.Errorf("config load failed: %w", err)
	}

	if err := viper.Unmarshal(&AppConfig); err != nil {
		return fmt.Errorf("config unmarshal failed: %w", err)
	}

	if AppConfig.ServerPort == "" {
		return errors.New("server_port is required")
	}

	return nil

}
