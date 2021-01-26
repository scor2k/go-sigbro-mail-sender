# go-sigbro-mail-sender

This app just connects to the RabbitMQ service and waiting for the new message to send it via Amazon SES. It uses Sentry for error tracking but you may skip this env. 

Every email should be Html and uses this pattern:
```
type emailJSON struct {
	Recipient string `json:"recipient"`
	Subject   string `json:"subject"`
	Body      string `json:"body"`
}
```

## Env

You should export these envs before start:

export MAILER_RABBITMQ_HOST=
export MAILER_RABBITMQ_USER=
export MAILER_RABBITMQ_PASS=
export MAILER_RABBITMQ_QUEUE=

export AWS_ACCESS_KEY_ID=
export AWS_SECRET_ACCESS_KEY=

export SENTRY_DSN=
