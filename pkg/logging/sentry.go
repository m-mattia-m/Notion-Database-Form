package logging

import (
	"Notion-Forms/global"
	"Notion-Forms/pkg/helper"
	"encoding/json"
	"errors"
	"github.com/getsentry/sentry-go"
	"log"
	"time"
)

var (
	dns *string
)

func init() {
	envDns, err := helper.GetEnv("LOGGING_DNS")
	if err != nil {
		log.Fatalln("" + err.Error())
	}
	dns = &envDns
}

func Client() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn: *dns,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	defer sentry.Flush(2 * time.Second)
}

func Info(targetServiceName string, method string, message Message) {
	msg := MessageRequest{
		TargetServiceName: targetServiceName,
		Method:            method,
		Message:           message,
	}
	StringMsg, err := stringifyMessage(msg)
	if err != nil {
		log.Println("can't convert logging-message to string: " + err.Error())
	}
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTag("level", "")
	})
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		appEnv, err := helper.GetEnv("APP_ENV")
		if err != nil || appEnv == "" {
			Error("logging/glitchtip", "Info", Message{
				Description: "[LOGGING]: get env was failed or env was empty",
				Detail:      err,
			})
		}
		scope.SetTag("env", appEnv)
	})
	// if you use zitadel, you can add this selection
	//if !reflect.DeepEqual(global.User, models.ZitadelUserinfo{}) {
	//	sentry.ConfigureScope(func(scope *sentry.Scope) {
	//		scope.SetUser(sentry.User{
	//			Email: global.User.Email,
	//			Name:  global.User.Name + " " + global.User.FamilyName,
	//		})
	//	})
	//}

	// remove this selection, if you use zitadel
	if global.Email != "" && global.Name != "" {
		sentry.ConfigureScope(func(scope *sentry.Scope) {
			scope.SetUser(sentry.User{
				Email: global.Email,
				Name:  global.Name,
			})
		})
	}
	sentry.CaptureMessage(StringMsg)
}

func Error(targetServiceName string, method string, message Message) {
	msg := MessageRequest{
		TargetServiceName: targetServiceName,
		Method:            method,
		Message:           message,
	}
	StringMsg, err := stringifyMessage(msg)
	if err != nil {
		log.Println("can't convert logging-message to string: " + err.Error())
	}
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTag("level", msg.Message.ErrorLevel)
	})
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		appEnv, err := helper.GetEnv("APP_ENV")
		if err != nil || appEnv == "" {
			Error("logging/glitchtip", "Error", Message{
				Description: "[LOGGING]: get env was failed or env was empty",
				Detail:      err,
			})
		}
		scope.SetTag("env", appEnv)
	})
	// if you use zitadel, you can add this selection
	//if !reflect.DeepEqual(global.User, models.ZitadelUserinfo{}) {
	//	sentry.ConfigureScope(func(scope *sentry.Scope) {
	//		scope.SetUser(sentry.User{
	//			Email: global.User.Email,
	//			Name:  global.User.Name + " " + global.User.FamilyName,
	//		})
	//	})
	//}

	// remove this selection, if you use zitadel
	if global.Email != "" && global.Name != "" {
		sentry.ConfigureScope(func(scope *sentry.Scope) {
			scope.SetUser(sentry.User{
				Email: global.Email,
				Name:  global.Name,
			})
		})
	}

	sentry.CaptureException(errors.New(StringMsg))
}

func stringifyMessage(msg MessageRequest) (string, error) {
	m, err := json.Marshal(msg)
	if err != nil {
		return "", err
	}
	return string(m), nil
}

type Message struct {
	Description string      `json:"description"`
	Detail      interface{} `json:"detail"`
	ErrorLevel  string      `json:"errorLevel"`
}

type MessageRequest struct {
	TargetServiceName string  `json:"targetServiceName"`
	Method            string  `json:"method"`
	Message           Message `json:"message"`
}
