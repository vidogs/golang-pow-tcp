package main

import (
	"log"
	"pow/internal/client_config"
	"pow/internal/client_wiring"
	"pow/internal/interrupt_handler"
)

func main() {
	cfg, err := client_config.NewClientConfig()

	if err != nil {
		log.Fatalf("error loading client config: %s\n", err)
	}

	wire := client_wiring.NewClientWiring(cfg)

	go func() {
		defer wire.ContextCancel()

		if err := wire.Client.Run(); err != nil {
			wire.Log.WithError(err).Errorf("client failed")
		}
	}()

	interrupt_handler.WaitForInterrupt(wire.Log, wire.Context, wire.ContextCancel)
}
