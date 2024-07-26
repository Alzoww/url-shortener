package config

import (
	"github.com/segmentio/kafka-go"
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

type Config struct {
	Env         string      `yaml:"env"`
	StoragePath string      `yaml:"storage_path" env-required:"true"`
	HttpServer  HttpServer  `yaml:"http_server"`
	Kafka       KafkaConfig `yaml:"kafka"`
}

type HttpServer struct {
	Address     string        `yaml:"address"`
	Timeout     time.Duration `yaml:"timeout"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
}

type KafkaConfig struct {
	Producers *kafka.WriterConfig `yaml:"producers"`
}

func LoadConfig(filename string, cfg interface{}) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return err
	}

	return nil
}
