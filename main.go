package main

import (
	"context"
	"net/http"
	"time"

	"github.com/rybit/lambda_example/util"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sirupsen/logrus"
)

var config configuration
var rootLogger *logrus.Entry
var client *http.Client

func main() {
	util.LoadConfig(&config)
	rootLogger = NewLogger(&config.LogConfig)
	if err := config.validate(); err != nil {
		rootLogger.WithError(err).Fatal("Invalid configuration")
	}

	client = new(http.Client)
	if config.TimeoutSec > 0 {
		client.Timeout = time.Second * time.Duration(config.TimeoutSec)
	}

	rootLogger.Debug("Startup completed")
	lambda.Start(handleEvent)
}

// handleEvent will decode the payload and send it to humio. Errors will only be returned if we could recover on retry
func handleEvent(ctx context.Context, input rawEvent) error {
	log := rootLogger.WithField("aws_id", util.RequestID(ctx))
	decoded, err := input.decode()
	if err != nil {
		log.WithError(err).Error("Failed to decode message")
		return nil // swallow because we don't want the msg again
	}
	log.Debug("Successfully decoded message")

	out := newHumioMsg(decoded)
	code, err := sendStrings(out)
	if err != nil {
		log.WithError(err).Error("Failed to post data to humio")
		return err
	}

	log.WithField("status_code", code).Info("Finished sending lines to humio")
	return nil
}
