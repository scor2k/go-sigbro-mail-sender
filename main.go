package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/bloom42/rz-go"
	"github.com/bloom42/rz-go/log"

	"github.com/streadway/amqp"

	sentry "github.com/getsentry/sentry-go"

	//go get -u github.com/aws/aws-sdk-go
	aws "github.com/aws/aws-sdk-go/aws"
	session "github.com/aws/aws-sdk-go/aws/session"
	ses "github.com/aws/aws-sdk-go/service/ses"
)

// using this vars for build
var appName = "go-sigbro-mail-sender"
var appVersion = "0.1.1"

// RabbitMQ connection
var rbmqHost = os.Getenv("MAILER_RABBITMQ_HOST")
var rbmqPort = os.Getenv("MAILER_RABBITMQ_PORT")
var rbmqUser = os.Getenv("MAILER_RABBITMQ_USER")
var rbmqPass = os.Getenv("MAILER_RABBITMQ_PASS")
var rbmqQueue = os.Getenv("MAILER_RABBITMQ_QUEUE")

// Sentry
var sentryDSN = os.Getenv("SENTRY_DSN")

// Some params for sending emails
const (
	sender    = "Sigbro  <noreply@nxter.org>"
	charSet   = "UTF-8"
	awsRegion = "eu-west-1"
)

// Template for every email via RabbitMQ
type emailJSON struct {
	Recipient string `json:"recipient"`
	Subject   string `json:"subject"`
	Body      string `json:"body"`
}

func main() {
	// set logger
	hostname, _ := os.Hostname()
	log.SetLogger(log.With(rz.Fields(rz.String("app", appName), rz.String("host", hostname))))

	// setup sentry
	sentryErr := sentry.Init(sentry.ClientOptions{
		Dsn:         sentryDSN,
		Environment: "",
		Release:     appName,
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
			// parse json from RabbitMQ
			email := emailJSON{}
			jsonErr := json.Unmarshal(d.Body, &email)

			if jsonErr != nil {
				log.Debug("Cannot parse message", rz.Bytes("Body", d.Body))
				sentry.CaptureMessage("Cannot parse email")
				continue
			}

			sendMailSES(email)
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

func sendMailSES(email emailJSON) {
	sess, err := session.NewSession(&aws.Config{Region: aws.String(awsRegion)})
	failOnError(err, "Cannot connect to Amazon SES")

	// Create an SES session.
	svc := ses.New(sess)

	// Assemble the email.
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(email.Recipient),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(charSet),
					Data:    aws.String(email.Body),
				},
				/*
					Text: &ses.Content{
							Charset: aws.String(charSet),
							Data:    aws.String(TextBody),
					},
				*/
			},
			Subject: &ses.Content{
				Charset: aws.String(charSet),
				Data:    aws.String(email.Subject),
			},
		},
		Source: aws.String(sender),
		// Uncomment to use a configuration set
		//ConfigurationSetName: aws.String(ConfigurationSet),
	}

	// Attempt to send the email.
	_, err = svc.SendEmail(input)

	if err != nil {
		sentry.CaptureMessage("Cannot send email")
		log.Error("Cannot send email", rz.Error("Amazon SES error", err))
		return
	}

	msg := fmt.Sprintf("Message to [%s] was sent.", email.Recipient)
	log.Info(msg)
}
