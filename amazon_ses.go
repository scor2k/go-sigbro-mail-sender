package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/bloom42/rz-go"
	"github.com/bloom42/rz-go/log"
	tg "github.com/scor2k/go-telegram-sender"
)

func sendMailSES(email emailJSON) error {
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
			},
			Subject: &ses.Content{
				Charset: aws.String(charSet),
				Data:    aws.String(email.Subject),
			},
		},
		Source: aws.String(sender),
	}

	// Attempt to send the email.
	_, err = svc.SendEmail(input)

	if err != nil {
		errTelegram := tg.SendMessage("[go-sigbro-mail-sender] Cannot send email via Amazon SES")
		if errTelegram != nil {
			log.Error("Cannot send Alert to the telegram", rz.Error("Error", errTelegram))
		}

		log.Error("Cannot send email", rz.Error("Amazon SES error", err))
		return err
	}

	msg := fmt.Sprintf("Message to [%s] was sent.", email.Recipient)
	log.Info(msg)
	return nil
}
