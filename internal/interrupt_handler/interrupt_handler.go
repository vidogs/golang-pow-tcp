package interrupt_handler

import (
	"context"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
)

func WaitForInterrupt(log *logrus.Logger, ctx context.Context, cancel context.CancelFunc) {
	log.Info("waiting for interruption...")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill)

	select {
	case <-ctx.Done():
		log.Info("context done")
	case <-interrupt:
		log.Warn("interrupted")
		cancel()
	}
}
