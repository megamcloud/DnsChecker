package main

import (
	"fmt"

	"github.com/caarlos0/env/v6"
)

type config struct {
	Debug           bool     `env:"APP_DEBUG" envDefault:"false"`
	ListenAddr      string   `env:"APP_LISTEN_ADDR" envDefault:":8080"`
	ApplicationName string   `env:"APP_NAME" envDefault:"DnsChecker-Local"`
	LogglyToken     string   `env:"APP_LOGGLY_TOKEN" envDefault:""`
	NameServers     []string `env:"APP_NAMESERVERS" envSeparator:"," envDefault:"8.8.8.8"`
	HostNames       []string `env:"APP_HOSTNAMES" envSeparator:"," envDefault:"google.com"`
}

func getConfig() *config {
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	if cfg.Debug {
		fmt.Printf("%+v\n", cfg)
	}

	return &cfg
}
