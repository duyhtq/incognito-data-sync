package service

import (
	"log"
	"os"

	config "github.com/duyhtq/incognito-data-sync/config"
	zapsentry "github.com/plimble/zap-sentry"
	"go.uber.org/zap"
)

func NewLogger(conf *config.Config) *zap.Logger {
	exec, err := os.Executable()
	if err != nil {
		log.Fatalf("Can't figure executable: %+v", err)
	}
	return zapsentry.New(
		zapsentry.WithStage(conf.Env),
		zapsentry.WithSentry(conf.SentryDSN, map[string]string{"exe": exec}, nil),
	).Desugar()
}
