package usecase

import (
	"github.com/NekruzRakhimov/auth_service/internal/adapter/driven/dbstore"
	"github.com/NekruzRakhimov/auth_service/internal/config"
	"github.com/NekruzRakhimov/auth_service/internal/port/driven"
	"github.com/NekruzRakhimov/auth_service/internal/port/usecase"
	authenticate "github.com/NekruzRakhimov/auth_service/internal/usecase/authenticator"
	"github.com/NekruzRakhimov/auth_service/internal/usecase/emails_getter"
	usercreater "github.com/NekruzRakhimov/auth_service/internal/usecase/user_creater"
)

type UseCases struct {
	UserCreater   usecase.UserCreater
	Authenticator usecase.Authenticate
	EmailsGetter  usecase.EmailsGetter
}

func New(cfg config.Config,
	store *dbstore.DBStore,
	amqp driven.AmqpProducer) *UseCases {
	return &UseCases{
		UserCreater:   usercreater.New(&cfg, store.UserStorage, amqp),
		Authenticator: authenticate.New(&cfg, store.UserStorage),
		EmailsGetter:  emails_getter.New(&cfg, store.UserStorage),
	}
}
