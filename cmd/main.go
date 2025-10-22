package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/NekruzRakhimov/auth_service/internal/bootstrap"
	"github.com/NekruzRakhimov/auth_service/internal/config"
	"github.com/sethvargo/go-envconfig"
)

func main() {
	var cfg config.Config

	err := envconfig.ProcessWith(context.TODO(), &envconfig.Config{Target: &cfg, Lookuper: envconfig.OsLookuper()})
	if err != nil {
		panic(err)
	}

	quitSignal := make(chan os.Signal, 1)
	signal.Notify(quitSignal, os.Interrupt)

	app := bootstrap.New(cfg)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-quitSignal
		cancel()
	}()

	app.Run(ctx)

}