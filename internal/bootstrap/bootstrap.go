package bootstrap

import (
	"context"
	"fmt"
	"github.com/NekruzRakhimov/auth_service/internal/adapter/driven/amqp"

	"github.com/NekruzRakhimov/auth_service/internal/adapter/driven/dbstore"
	"github.com/NekruzRakhimov/auth_service/internal/config"
	"github.com/NekruzRakhimov/auth_service/internal/usecase"
)

func initLayers(cfg config.Config) *App {
	teardown := make([]func(), 0)
	//log := logger.New(cfg.LogLevel, config.ServiceLabel, zap.WithCaller(true))

	_, amqpCh := amqp.InitAMQPProducer(cfg.AMQPURL)
	teardown = append(teardown,
		func() {
			err := amqpCh.Close()
			if err != nil {
				fmt.Printf("Error closing AMQP consumer: %s\n", err)
				return
			}
		},
	)

	authQueue, err := amqp.InitQueue(amqpCh, "auth-queue")
	if err != nil {
		fmt.Println(err)
	}

	amqpProducers := amqp.NewProducersAMQP(authQueue, amqpCh)

	db, err := initDB(*cfg.Postgres)
	if err != nil {
		panic(err)
	}

	storage := dbstore.New(db)
	//log.Info("Database connection established")

	teardown = append(teardown, func() {
		if err := db.Close(); err != nil {
			fmt.Println(err)
			//log.Error(err.Error())
		}
	})

	uc := usecase.New(cfg, storage, amqpProducers)

	httpSrv := initHTTPService(&cfg, uc)

	teardown = append(teardown,
		func() {
			//log.Info("HTTP is shutting down")
			ctxShutDown, cancel := context.WithTimeout(context.Background(), gracefulDeadline)
			defer cancel()
			if err := httpSrv.Shutdown(ctxShutDown); err != nil {
				//log.Error(fmt.Sprintf("server Shutdown Failed:%s", err))
				return
			}
			//log.Info("HTTP is shut down")
		},
	)

	return &App{
		cfg:      cfg,
		rest:     httpSrv,
		teardown: teardown,
	}
}
