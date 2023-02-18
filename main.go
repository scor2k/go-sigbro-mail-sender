package main

import (
	"github.com/bloom42/rz-go"
	"github.com/bloom42/rz-go/log"
	"os"
)

// using this vars for build
var appName = "go-sigbro-mail-sender"
var appVersion = "1.3.0"

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

	startConsume()
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Error(msg, rz.Error("rbmq error", err))
		os.Exit(1)
	}
}
