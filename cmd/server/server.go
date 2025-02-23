package main

import (
	"log"
	"pow/internal/interrupt_handler"
	"pow/internal/server_config"
	"pow/internal/server_wiring"
)

func main() {
	cfg, err := server_config.NewServerConfig()

	if err != nil {
		log.Fatalf("error loading server config: %s\n", err)
	}

	wire := server_wiring.NewServerWiring(cfg)

	go func() {
		defer wire.ContextCancel()

		if err := wire.Server.Run(); err != nil {
			wire.Log.WithError(err).Errorf("server failed")
		}
	}()

	interrupt_handler.WaitForInterrupt(wire.Log, wire.Context, wire.ContextCancel)
}
