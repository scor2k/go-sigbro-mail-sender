package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/bloom42/rz-go"
	"github.com/bloom42/rz-go/log"
	tg "github.com/scor2k/go-telegram-sender"
	"os"
	"time"
)

func startConsume() {
	log.Info("Start consuming")

	sess, errSession := session.NewSession(&aws.Config{
		Region:      aws.String(sqsRegion),
		Credentials: credentials.NewStaticCredentials(sqsAccessKey, sqsSecretKey, ""),
	})
	if errSession != nil {
		log.Error("Cannot connect to the Amazon", rz.Error("Error", errSession))

		errTelegram := tg.SendMessage("[go-sigbro-mail-sender] Cannot connect to Amazon SQS")
		if errTelegram != nil {
			log.Error("Cannot send Alert to the telegram", rz.Error("Error", errTelegram))
		}

		return
	}

	svc := sqs.New(sess)

	retryCounter := 0
consume:
	retryCounter += 1
	if retryCounter > sqsRetry {
		log.Info("buy-buy, see you next time.")
		os.Exit(0)
	}

	msgResult, errSqs := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
		AttributeNames: []*string{
			aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
		},
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		QueueUrl:            &sqsURL,
		MaxNumberOfMessages: aws.Int64(1),
		VisibilityTimeout:   aws.Int64(600), // do not try to send the same email twice for 10 minutes
	})

	if errSqs != nil {
		log.Error("Cannot get Msg from the Queue", rz.Error("Error", errSqs))
		errTelegram := tg.SendMessage("[go-sigbro-mail-sender] Cannot get message from Amazon SQS")
		if errTelegram != nil {
			log.Error("Cannot send Alert to the telegram", rz.Error("Error", errTelegram))
		}

		time.Sleep(3 * time.Second)
		os.Exit(1)
	}

	// DEBUG
	msg := fmt.Sprintf("Message text: %+v", msgResult.Messages)
	log.Debug(msg)

	if len(msgResult.Messages) > 0 {
		body := *msgResult.Messages[0].Body
		recipient := *msgResult.Messages[0].MessageAttributes["Recipient"].StringValue
		subject := *msgResult.Messages[0].MessageAttributes["Subject"].StringValue

		log.Debug("Get a message", rz.String("To", recipient), rz.String("Subject", subject))

		errSES := sendMailSES(emailJSON{
			Body:      body,
			Recipient: recipient,
			Subject:   subject,
		})

		if errSES == nil {

			log.Info("Message was send", rz.String("To", recipient), rz.String("Subject", subject))
			// remove item from the queue
			_, delErr := svc.DeleteMessage(&sqs.DeleteMessageInput{
				QueueUrl:      &sqsURL,
				ReceiptHandle: msgResult.Messages[0].ReceiptHandle,
			})

			if delErr != nil {
				log.Error("Cannot remove message from the queue", rz.Error("Error", delErr))
			}

		} else {
			log.Error("Cannot send email", rz.String("To", recipient), rz.String("Subject", subject), rz.Error("Error", errSES))
		}

		goto consume

	} else {
		log.Debug("Nothing to consume. Wait and try again")
		time.Sleep(2 * time.Second)
		goto consume
	}
}
