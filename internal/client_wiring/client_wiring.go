package client_wiring

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"pow/internal/client"
	"pow/internal/client_config"
	"pow/internal/logging"
	"pow/internal/proof_of_work"
)

type ClientWiring struct {
	Config        *client_config.ClientConfig
	Log           *logrus.Logger
	Context       context.Context
	ContextCancel context.CancelFunc
	ProofOfWork   *proof_of_work.ProofOfWork
	Client        *client.Client
}

func NewClientWiring(cfg *client_config.ClientConfig) *ClientWiring {
	log := logging.NewLogging(cfg.Logger.Level, cfg.Logger.Format)

	ctx, contextCancel := context.WithCancel(context.Background())

	proofOfWork := proof_of_work.NewProofOfWork()

	return &ClientWiring{
		Config:        cfg,
		Log:           log,
		Context:       ctx,
		ContextCancel: contextCancel,
		ProofOfWork:   proofOfWork,
		Client: client.NewClient(
			ctx,
			log.WithField("service", "client"),
			fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
			proofOfWork,
		),
	}
}
