package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/sebest/logrusly"
	log "github.com/sirupsen/logrus"
)

var applicationName = getEnvOrDefault("APP_NAME", "DnsChecker-Local")
var logglyToken = getEnvOrDefault("APP_LOGGLY_TOKEN", "")
var nameServers = getEnvArrayOrDefault("APP_NAMESERVERS", []string{"8.8.8.8"})
var hostNames = getEnvArrayOrDefault("APP_HOSTNAMES", []string{"google.com"})

var logger = log.WithField("Application", applicationName)
var loggly = logrusly.NewLogglyHook(logglyToken, applicationName, log.InfoLevel, applicationName)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	if logglyToken != "" {
		log.AddHook(loggly)
	} else {
		log.Info("No Loggly token found, only logging to console")
	}
}

func main() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	ctx := context.Background()

	for _, nameServer := range nameServers {
		go checkNameServer(ctx, nameServer, hostNames)
	}

	<-c
	logger.Info("Starting graceful shutdown")

	loggly.Flush()
	ctx, cancel := context.WithTimeout(ctx, 5000*time.Millisecond)
	defer cancel()

	logger.Info("Stopping application")
}

func checkNameServer(ctx context.Context, nameServer string, hostNames []string) {
	nsLogger := logger.WithField("NameServer", nameServer)
	nsLogger.Info("Start checking ", hostNames, " on nameserver ", nameServer)

	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{}
			return d.DialContext(ctx, "udp", net.JoinHostPort(nameServer, "53"))
		},
	}

	for true {
		for _, hostName := range hostNames {
			hostNameLogger := nsLogger.WithField("HostName", hostName)
			ips, err := resolver.LookupIPAddr(ctx, hostName)
			if err != nil {
				hostNameLogger.Error(err)
				continue
			}

			hostNameLogger.Debug(ips)
		}

		time.Sleep(10000 * time.Millisecond)
	}
}

func getEnvOrDefault(envVar, fallback string) string {
	if value, ok := os.LookupEnv(envVar); ok {
		return value
	}

	return fallback
}

func getEnvArrayOrDefault(envVar string, fallback []string) []string {
	if value := getEnvOrDefault(envVar, ""); value != "" {
		return strings.Split(value, ",")
	}

	return fallback
}
