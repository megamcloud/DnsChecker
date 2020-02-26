package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

type app struct {
	Logger  *logrus.Entry
	Config  *config
	Metrics *metrics
}

func newApp() *app {
	var config = getConfig()
	var logger = newLogger(config)
	var metrics = newMetrics(config)

	return &app{
		Logger:  logger,
		Config:  config,
		Metrics: metrics,
	}
}

func main() {
	app := newApp()
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	ctx := context.Background()

	for _, nameServer := range app.Config.NameServers {
		go app.checkNameServer(ctx, nameServer)
	}

	go app.serveMetrics()

	<-c
	app.Logger.Info("Starting graceful shutdown")

	ctx, cancel := context.WithTimeout(ctx, 5000*time.Millisecond)
	defer cancel()

	app.Logger.Info("Stopping application")
}

func (app *app) checkNameServer(ctx context.Context, nameServer string) {
	nsLogger := app.Logger.WithField("NameServer", nameServer)
	nsLogger.Info("Start checking ", app.Config.HostNames, " on nameserver ", nameServer)

	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{}
			return d.DialContext(ctx, "udp", net.JoinHostPort(nameServer, "53"))
		},
	}

	for true {
		for _, hostName := range app.Config.HostNames {
			hostNameLogger := nsLogger.WithField("HostName", hostName)
			ips, err := resolver.LookupIPAddr(ctx, hostName)
			if err != nil {
				hostNameLogger.Error(err)
				app.Metrics.incDNSCounter(false, nameServer, hostName)
				continue
			}

			app.Metrics.incDNSCounter(true, nameServer, hostName)
			hostNameLogger.Debug(ips)
		}

		time.Sleep(10000 * time.Millisecond)
	}
}
