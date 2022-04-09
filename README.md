# go-sigbro-mail-sender

This app just connects to the RabbitMQ service and waiting for the new message to send it via Amazon SES. It uses Sentry for error tracking but you may skip this env. 

Every email should be Html and uses this pattern:
```
type emailJSON struct {
	Recipient string `json:"recipient"`
	Subject   string `json:"subject"`
	Body      string `json:"body"`
	Sender    string `json:"sender"`
}
```

## Envs

You should/may export these envs before start:

```
export LOG_LEVEL=DEBUG

# Amazon SES
export AWS_ACCESS_KEY_ID=
export AWS_SECRET_ACCESS_KEY=

# Sentry
export SENTRY_DSN=

# Amazon SQS
export SIGBRO_SQS_REGION=eu-west-1
export SIGBRO_SQS_URL=
export SIGBRO_SQS_ACCESS_KEY=
export SIGBRO_SQS_SECRET_KEY=

# Telegram
export TELEGRAM_TOKEN=
export TELEGRAM_CHAT=

```

## Changelog
*1.2.0 version*
 - feat: remove sentry

*1.1.0 version*
 - feat: ability to change sender

*1.0.1 version*
 - switched to json body instead of plain/text + params

*1.0.0 version*
 - uses Amazon SQS as a queue
 - send alerts about SQS/SES issues via telegram

*0.1.1 version*
 - uses RabbitMQ as a queue
 - send mail via Amazon SES
