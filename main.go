package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/bloom42/rz-go"
	"github.com/bloom42/rz-go/log"

	"github.com/streadway/amqp"

	sentry "github.com/getsentry/sentry-go"
)

var rbmqHost = os.Getenv("MAILER_RABBITMQ_HOST")
var rbmqPort = os.Getenv("MAILER_RABBITMQ_PORT")
var rbmqUser = os.Getenv("MAILER_RABBITMQ_USER")
var rbmqPass = os.Getenv("MAILER_RABBITMQ_PASS")
var rbmqQueue = os.Getenv("MAILER_RABBITMQ_QUEUE")

var sentryDSN = os.Getenv("SENTRY_DSN")

func main() {
	// set logger
	hostname, _ := os.Hostname()
	log.SetLogger(log.With(rz.Fields(rz.String("app", "go-sigbro-mail-sender"), rz.String("host", hostname))))

	// setup sentry
	sentryErr := sentry.Init(sentry.ClientOptions{
		Dsn:         sentryDSN,
		Environment: "",
		Release:     "go-sigbro-mail-sender",
		Debug:       true,
	})
	failOnError(sentryErr, "Cannot initialyze Sentry")
	defer sentry.Flush(2 * time.Second)

	//sentry.CaptureMessage("go-sigbro-mail-sender running...")

	// set rabbitMQ port
	port, _ := strconv.Atoi(rbmqPort)
	if port == 0 {
		port = 5672
	}

	// connect to RabbitMQ
	rbmqDSN := fmt.Sprintf("amqp://%s:%s@%s:%d/", rbmqUser, rbmqUser, rbmqHost, port)
	conn, err := amqp.Dial(rbmqDSN)

	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	queue, err := ch.QueueDeclare(
		rbmqQueue,
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		queue.Name,
		"go-sigbro-mail-sender",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register a consumer")

	doForever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Debug("Message", rz.Bytes("Body", d.Body))
		}
	}()

	log.Debug("Waiting for messages")

	<-doForever

}

func failOnError(err error, msg string) {
	if err != nil {
		log.Error(msg, rz.Error("rbmq error", err))
		os.Exit(1)
	}
}
