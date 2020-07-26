package utils

import (
	"github.com/Voodfy/voodfy-transcoder/pkg/logging"
	"github.com/getsentry/sentry-go"
)

// SendError send error to sentry
func SendError(function string, err error) {
	if err != nil {
		logging.Error(function, err.Error())
		sentry.CaptureException(err)
	}
}
