package main

import (
	"os"
	"time"

	"github.com/bloom42/rz-go"
	"github.com/bloom42/rz-go/log"

	sentry "github.com/getsentry/sentry-go"
)

// using this vars for build
var appName = "go-sigbro-mail-sender"
var appVersion = "1.0.1"

func main() {
	// set logger
	hostname, _ := os.Hostname()
	ll := os.Getenv("LOG_LEVEL")
	logLevel := rz.InfoLevel
	if ll == "DEBUG" {
		logLevel = rz.DebugLevel
	}

	log.SetLogger(
		log.With(
			rz.Fields(rz.String("ver", appVersion), rz.String("app", appName), rz.String("host", hostname)),
			rz.Level(logLevel),
		))

	// setup sentry
	sentryErr := sentry.Init(sentry.ClientOptions{
		Dsn:         sentryDSN,
		Environment: "",
		Release:     appName,
		Debug:       true,
	})
	failOnError(sentryErr, "Cannot initialize Sentry")
	defer sentry.Flush(2 * time.Second)

	startConsume()

	os.Exit(0)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Error(msg, rz.Error("rbmq error", err))
		os.Exit(1)
	}
}
