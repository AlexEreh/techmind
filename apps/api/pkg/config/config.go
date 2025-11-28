package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Postgres struct {
		Conn string `yaml:"conn" mapstructure:"conn"`
	} `yaml:"postgres" mapstructure:"postgres"`
	MinIO struct {
		Endpoint  string `yaml:"endpoint" mapstructure:"endpoint"`
		AccessKey string `yaml:"access_key" mapstructure:"access_key"`
		SecretKey string `yaml:"secret_key" mapstructure:"secret_key"`
		UseSSL    bool   `yaml:"use_ssl" mapstructure:"use_ssl"`
	} `yaml:"minio" mapstructure:"minio"`
	Elasticsearch struct {
		URL      string `yaml:"url" mapstructure:"url"`
		Username string `yaml:"username" mapstructure:"username"`
		Password string `yaml:"password" mapstructure:"password"`
	} `yaml:"elasticsearch" mapstructure:"elasticsearch"`
	Gotenberg struct {
		URL     string `yaml:"url" mapstructure:"url"`
		Enabled bool   `yaml:"enabled" mapstructure:"enabled"`
		Timeout int    `yaml:"timeout" mapstructure:"timeout"` // in seconds
	} `yaml:"gotenberg" mapstructure:"gotenberg"`
	HTTPPort int `yaml:"http_port" mapstructure:"http_port"`

	JWT struct {
		SecretKey            string `yaml:"secret_key" mapstructure:"secret_key"`
		AccessTokenLifetime  string `yaml:"access_token_lifetime" mapstructure:"access_token_lifetime"`
		RefreshTokenLifetime string `yaml:"refresh_token_lifetime" mapstructure:"refresh_token_lifetime"`
	} `yaml:"jwt" mapstructure:"jwt"`
}

func Load() (*Config, error) {
	var cfg Config

	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	var envPath string
	if strings.HasPrefix(wd, "/app") {
		wd = "/app"
		envPath = filepath.Join(wd, "prod.yml")
	} else {
		wd = filepath.Join(wd)
		envPath = filepath.Join(wd, "../../dev.yml")
	}

	fmt.Printf("Loading config from %s\n", envPath)

	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		return nil, err
	}

	viper.SetConfigFile(envPath)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}
