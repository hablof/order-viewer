package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	// to use another config file set env var CFGFILENAME
	defaultFile = "default_config.yaml"
)

type Config struct {
	Debug     bool `yaml:"debug"`
	Consumers int  `yaml:"consumers"`

	PostgresHost     string `yaml:"postgres_host"`
	PostgresPort     string `yaml:"postgres_port"`
	PostgresDB       string `yaml:"postgres_database"`
	PostgresUser     string `yaml:"postgres_user"`
	PostgresPassword string

	HTTPPort    string        `yaml:"http_port"`
	HTTPTimeout time.Duration `yaml:"http_timeout"`

	NatsHost              string `yaml:"nats_host"`
	NatsPort              string `yaml:"nats_port"`
	NatsCluster           string `yaml:"nats_cluster"`
	NatsClientName        string `yaml:"nats_client_name"`
	NatsClientID          string `yaml:"nats_client_id"`
	NatsDurableQueueGroup string `yaml:"nats_durable_queue_group"`
	NatsSubject           string `yaml:"nats_subject"`
}

func NewConfig() (Config, error) {

	psqlPassword, ok := os.LookupEnv("PSQLPASS")
	if !ok {
		return Config{}, errors.New("env var psqlpass is not set")
	}

	file := defaultFile
	overwritedFileName, ok := os.LookupEnv("CFGFILENAME")
	if ok {
		file = overwritedFileName
	}

	f, err := os.Open(file)
	if err != nil {
		return Config{}, err
	}

	cfg := Config{}
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return Config{}, err
	}

	cfg.PostgresPassword = psqlPassword

	//todo: overwrite yaml values with env vars?

	return cfg, nil
}

func (c Config) GetPgURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.PostgresUser,
		c.PostgresPassword,
		c.PostgresHost,
		c.PostgresPort,
		c.PostgresDB,
	)
}

func (c Config) GetNatsURL() string {
	return fmt.Sprintf("nats://%s:%s", c.NatsHost, c.NatsPort)
}
