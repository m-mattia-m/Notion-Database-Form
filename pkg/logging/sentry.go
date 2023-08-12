package logging

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/getsentry/sentry-go"
	sentryotel "github.com/getsentry/sentry-go/otel"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"log"
	"time"
)

type Client struct {
	dsn                   *string
	env                   *string
	enableLogging         *string
	enableExternalLogging *string
}

func New(appEnv, envDns, enableLogging, enableExternalLogging string) (*Client, error) {
	if enableLogging == "" && appEnv != "PROD" {
		enableLogging = "true"
	}
	if enableLogging == "" && appEnv == "PROD" {
		enableLogging = "false"
	}

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	err := sentry.Init(sentry.ClientOptions{
		Dsn:              envDns,
		Debug:            true,
		EnableTracing:    true,
		AttachStacktrace: true,
		SampleRate:       0,
		TracesSampleRate: 1.0,
		//ServerName:       "",
		//Release:          "",
		Environment: appEnv,
		TracesSampler: func(ctx sentry.SamplingContext) float64 {
			return 1.0
		},
	})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("sentry.Init: %s", err))
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(sentryotel.NewSentrySpanProcessor()),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(sentryotel.NewSentryPropagator())

	defer sentry.Flush(2 * time.Second)

	return &Client{
		dsn:                   &envDns,
		env:                   &appEnv,
		enableLogging:         &enableLogging,
		enableExternalLogging: &enableExternalLogging,
	}, nil
}

func (loggingClient Client) Info(targetServiceName string, method string, message Message) {
	if *loggingClient.env == "" {
		log.Fatalln(fmt.Sprintf("[LOGGING-ERROR]: get env 'APP_ENV' was failed or env is empty"))
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

	if *loggingClient.enableLogging == "true" {
		log.Println(StringMsg)
	}

	if *loggingClient.env != "PROD" && *loggingClient.enableExternalLogging == "false" {
		return
	}

	sentry.ConfigureScope(func(scope *sentry.Scope) {
		if *loggingClient.env != "" {
			scope.SetTag("env", *loggingClient.env)
		}
		if msg.Message.ErrorLevel != "" {
			scope.SetTag("level", msg.Message.ErrorLevel.String())
		}
		if msg.Message.LogId != "" {
			scope.SetExtra("logId", msg.Message.LogId)
		}
	})
	sentry.CaptureMessage(StringMsg)
	sentry.Flush(time.Second * 5)
}

func (loggingClient Client) Error(targetServiceName string, method string, message Message) {
	if *loggingClient.env == "" {
		log.Fatalln(fmt.Sprintf("[LOGGING-ERROR]: get env 'APP_ENV' was failed or env is empty"))
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

	if *loggingClient.enableLogging == "true" {
		log.Println(StringMsg)
	}

	if *loggingClient.env != "PROD" && *loggingClient.enableExternalLogging == "false" {
		return
	}

	sentry.ConfigureScope(func(scope *sentry.Scope) {
		if *loggingClient.env != "" {
			scope.SetTag("env", *loggingClient.env)
		}
		if msg.Message.ErrorLevel != "" {
			scope.SetTag("level", msg.Message.ErrorLevel.String())
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

type errorLevel string

const (
	UltraHigh errorLevel = "UltraHigh"
	High      errorLevel = "High"
	Medium    errorLevel = "Medium"
	Low       errorLevel = "Low"
	Neglected errorLevel = "Neglected"
)

func (el errorLevel) String() string {
	return string(el)
}

type Message struct {
	Description string      `json:"description"`
	Detail      interface{} `json:"detail"`
	ErrorLevel  errorLevel  `json:"error_level"`
	LogId       string      `json:"log_id"`
}

type messageRequest struct {
	TargetServiceName string  `json:"target_service_name"`
	Method            string  `json:"method"`
	Message           Message `json:"message"`
}
