package usercreater

import (
	"context"
	"errors"

	"github.com/NekruzRakhimov/auth_service/internal/config"
	"github.com/NekruzRakhimov/auth_service/internal/domain"
	"github.com/NekruzRakhimov/auth_service/internal/errs"
	"github.com/NekruzRakhimov/auth_service/internal/port/driven"
	"github.com/NekruzRakhimov/auth_service/utils"
)

type UseCase struct {
	cfg *config.Config
	userStorage driven.UserStorage
}

func New(cfg *config.Config, userStorage driven.UserStorage) *UseCase {
	return &UseCase{
		cfg: cfg,
		userStorage: userStorage,
	}
}

func (u *UseCase) CreateUser(ctx context.Context, user domain.User) (err error) {
	// Проверить существует ли пользователь с таким username'ом в бд
	_, err = u.userStorage.GetUserByUsername(ctx, user.Username)
	if err != nil {
		if !errors.Is(err, errs.ErrNotfound) {
			return err
		}
	} else {
		return errs.ErrUsernameAlreadyExists
	}

	// За хэшировать пароль
	user.Password, err = utils.GenerateHash(user.Password)
	if err != nil {
		return err
	}

	user.Role = domain.RoleUser

	// Добавить пользователя в бд
	if err = u.userStorage.CreateUser(ctx, user); err != nil {
		return err
	}

	return nil
}