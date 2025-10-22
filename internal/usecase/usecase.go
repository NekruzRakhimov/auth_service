package usecase

import (
	"github.com/NekruzRakhimov/auth_service/internal/adapter/driven/dbstore"
	"github.com/NekruzRakhimov/auth_service/internal/config"
	"github.com/NekruzRakhimov/auth_service/internal/port/usecase"
	"github.com/NekruzRakhimov/auth_service/internal/usecase/authenticator"
	usercreater "github.com/NekruzRakhimov/auth_service/internal/usecase/user_creater"
)

type UseCases struct {
	UserCreater usecase.UserCreater
	Authenticator usecase.Authenticate
}

func New(cfg config.Config, store *dbstore.DBStore,) *UseCases {
	return &UseCases{
		UserCreater: usercreater.New(&cfg, store.UserStorage),
		Authenticator: authenticate.New(&cfg, store.UserStorage),
	}
}