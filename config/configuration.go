package config

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"time"
	
	"github.com/heetch/confita"
	"github.com/heetch/confita/backend"
	"github.com/heetch/confita/backend/env"
	"github.com/heetch/confita/backend/flags"
)

// Config is the exporter CLI configuration.
type Config struct {
	ServerPort      string        `config:"server_port"`
	WgPeerFile      string        `config:"wg_peer_file"`
	Interval        time.Duration `config:"interval"`
}

func getDefaultConfig() *Config {
	return &Config{
		ServerPort:      "9586",
		WgPeerFile:      "",
		Interval:        10 * time.Second,
	}
}

// Load method loads the configuration by using both flag or environment variables.
func Load() *Config {
	loaders := []backend.Backend{
		env.NewBackend(),
		flags.NewBackend(),
	}

	loader := confita.NewLoader(loaders...)

	cfg := getDefaultConfig()
	err := loader.Load(context.Background(), cfg)
	if err != nil {
		panic(err)
	}

	cfg.show()

	return cfg
}

func (c Config) show() {
	val := reflect.ValueOf(&c).Elem()
	log.Println("---------------------------------------")
	log.Println("- Wireguard exporter configuration -")
	log.Println("---------------------------------------")
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)

		log.Println(fmt.Sprintf("%s : %v", typeField.Name, valueField.Interface()))
	}
	log.Println("---------------------------------------")
}
