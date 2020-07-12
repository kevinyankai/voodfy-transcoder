package main

import (
	"log"
	"time"

	"github.com/Voodfy/voodfy-transcoder/internal/logging"
	"github.com/Voodfy/voodfy-transcoder/internal/models"
	"github.com/Voodfy/voodfy-transcoder/internal/queue"
	"github.com/Voodfy/voodfy-transcoder/internal/settings"
	st "github.com/getsentry/sentry-go"
)

func sentry() {
	err := st.Init(st.ClientOptions{
		Dsn: settings.AppSetting.SentryDNS,
	})
	log.Println(settings.AppSetting.SentryDNS)
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	defer st.Flush(2 * time.Second)
}

func init() {
	sentry()
	settings.Setup()
	models.InitDB()
	logging.Setup()
}

func main() {
	wrk := queue.NewWorker()
	if settings.AppSetting.QueueEnabled {
		wrk.Launch()
	}

}
