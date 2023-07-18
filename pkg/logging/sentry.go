package logging

import (
	"Notion-Forms/pkg/logging/helper"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/getsentry/sentry-go"
	"log"
	"time"
)

type Client struct {
	Dsn           *string
	Env           *string
	EnableLogging *string
	//SlackClient   *slack.Client
}

func New() (*Client, error) {
	appEnv, err := helper.GetEnv("APP_ENV")
	if err != nil {
		return nil, err
	}

	envDns, err := helper.GetEnv("LOGGING_DNS")
	if err != nil {
		log.Fatalln("" + err.Error())
		return nil, err
	}

	enableLogging, err := helper.GetEnv("ENABLE_LOGGING")
	if err != nil {
		log.Fatalln("" + err.Error())
		return nil, err
	}

	if enableLogging == "" && appEnv != "PROD" {
		enableLogging = "true"
	}
	if enableLogging == "" && appEnv == "PROD" {
		enableLogging = "false"
	}

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	err = sentry.Init(sentry.ClientOptions{
		Dsn: envDns,
	})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("sentry.Init: %s", err))
	}
	defer sentry.Flush(2 * time.Second)

	return &Client{
		Dsn:           &envDns,
		Env:           &appEnv,
		EnableLogging: &enableLogging,
	}, nil
}

func (loggingClient Client) Info(targetServiceName string, method string, message Message) {
	if *loggingClient.Env == "" {
		log.Fatalln(fmt.Sprintf("[LOGGING-ERROR]: get Env 'APP_ENV' was failed or Env is empty"))
	}

	msg := messageRequest{
		TargetServiceName: targetServiceName,
		Method:            method,
		Message: Message{
			Description: message.Description,
			Detail:      fmt.Sprintf("%s", message.Detail),
			ErrorLevel:  message.ErrorLevel,
			LogId:       message.LogId,
		},
	}
	StringMsg, err := stringifyMessage(msg)
	if err != nil {
		log.Println("can't convert logging-message to string: " + err.Error())
	}

	if *loggingClient.EnableLogging == "true" {
		log.Println(StringMsg)
	}

	if *loggingClient.Env != "PROD" {
		return
	}

	sentry.ConfigureScope(func(scope *sentry.Scope) {
		if *loggingClient.Env != "" {
			scope.SetTag("Env", *loggingClient.Env)
		}
		if msg.Message.ErrorLevel != "" {
			scope.SetTag("level", msg.Message.ErrorLevel)
		}
		if msg.Message.LogId != "" {
			scope.SetExtra("logId", msg.Message.LogId)
		}
	})
	sentry.CaptureMessage(StringMsg)
	sentry.Flush(time.Second * 5)
}

func (loggingClient Client) Error(targetServiceName string, method string, message Message) {
	if *loggingClient.Env == "" {
		log.Fatalln(fmt.Sprintf("[LOGGING-ERROR]: get Env 'APP_ENV' was failed or Env is empty"))
	}

	msg := messageRequest{
		TargetServiceName: targetServiceName,
		Method:            method,
		Message: Message{
			Description: message.Description,
			Detail:      fmt.Sprintf("%s", message.Detail),
			ErrorLevel:  message.ErrorLevel,
			LogId:       message.LogId,
		},
	}
	StringMsg, err := stringifyMessage(msg)
	if err != nil {
		log.Println("can't convert logging-message to string: " + err.Error())
	}

	if *loggingClient.EnableLogging == "true" {
		log.Println(StringMsg)
	}

	if *loggingClient.Env != "PROD" {
		return
	}

	sentry.ConfigureScope(func(scope *sentry.Scope) {
		if *loggingClient.Env != "" {
			scope.SetTag("Env", *loggingClient.Env)
		}
		if msg.Message.ErrorLevel != "" {
			scope.SetTag("level", msg.Message.ErrorLevel)
		}
		if msg.Message.LogId != "" {
			scope.SetExtra("logId", msg.Message.LogId)
		}
	})
	sentry.CaptureException(errors.New(StringMsg))
	sentry.Flush(time.Second * 5)
}

func stringifyMessage(msg messageRequest) (string, error) {
	m, err := json.Marshal(msg)
	if err != nil {
		return "", err
	}
	return string(m), nil
}

type Message struct {
	Description string      `json:"description"`
	Detail      interface{} `json:"detail"`
	ErrorLevel  string      `json:"error_level"`
	LogId       string      `json:"log_id"`
}

type messageRequest struct {
	TargetServiceName string  `json:"target_service_name"`
	Method            string  `json:"method"`
	Message           Message `json:"message"`
}
