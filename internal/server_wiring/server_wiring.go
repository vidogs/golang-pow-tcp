package server_wiring

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"pow/internal/dao"
	"pow/internal/logging"
	"pow/internal/proof_of_work"
	"pow/internal/server"
	"pow/internal/server_config"
)

type ServerWiring struct {
	Config           *server_config.ServerConfig
	Log              *logrus.Logger
	Context          context.Context
	ContextCancel    context.CancelFunc
	WordsOfWisdomDao *dao.WordsOfWisdom
	ProofOfWork      *proof_of_work.ProofOfWork
	Server           *server.Server
}

func NewServerWiring(cfg *server_config.ServerConfig) *ServerWiring {
	log := logging.NewLogging(cfg.Logger.Level, cfg.Logger.Format)

	ctx, contextCancel := context.WithCancel(context.Background())

	wordsOfWisdomDao := dao.NewWordsOfWisdom(
		log.WithField("service", "dao").WithField("dao", "words_of_wisdon"),
		cfg.Dao.WordsOfWisdom,
	)

	proofOfWork := proof_of_work.NewProofOfWork()

	return &ServerWiring{
		Config:           cfg,
		Log:              log,
		Context:          ctx,
		ContextCancel:    contextCancel,
		WordsOfWisdomDao: wordsOfWisdomDao,
		ProofOfWork:      proofOfWork,
		Server: server.NewServer(
			ctx,
			log.WithField("service", "server"),
			fmt.Sprintf("%s:%d", cfg.Bind.Host, cfg.Bind.Port),
			cfg.Settings.SolveTimeout,
			cfg.Settings.Challenge,
			proofOfWork,
			wordsOfWisdomDao,
		),
	}
}
