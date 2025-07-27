package config

import (
	"fmt"
	"log"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

type Config struct {
	WsServer struct {
		Port             int    `yaml:"port"`
		WSEndpoint       string `yaml:"endpoint"`
		WriterBufferSize int    `yaml:"writer_buffer_size"`
	} `yaml:"ws_server"`

	NATS struct {
		URL string `yaml:"url"`
	} `yaml:"nats"`
}

var (
	cfg  *Config
	once sync.Once
)

const PATH string = "config/config.yaml"

func loadConfig() *Config {
	f, err := os.Open(PATH)
	if err != nil {
		log.Fatalf("cannot open config file: %v", err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	var c Config
	if err := decoder.Decode(&c); err != nil {
		log.Fatalf("cannot decode config: %v", err)
	}
	return &c
}

func Init() {
	once.Do(func() {
		cfg = loadConfig()
	})
}

func Get() (*Config, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config not initialized. Call config.Init() first")
	}
	return cfg, nil
}
