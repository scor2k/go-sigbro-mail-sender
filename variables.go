package main

import "os"

// Sentry
var (
	sentryDSN = os.Getenv("SENTRY_DSN")
)

// Amazon SQS
var (
	sqsRetry     = 30 // 20 seconds for long pooling and 30 retries ~ 10 min
	sqsRegion    = os.Getenv("SIGBRO_SQS_REGION")
	sqsURL       = os.Getenv("SIGBRO_SQS_URL")
	sqsAccessKey = os.Getenv("SIGBRO_SQS_ACCESS_KEY")
	sqsSecretKey = os.Getenv("SIGBRO_SQS_SECRET_KEY")
)

// Some params for sending emails
const (
	sender    = "Sigbro  <noreply@nxter.org>"
	charSet   = "UTF-8"
	awsRegion = "eu-west-1"
)
