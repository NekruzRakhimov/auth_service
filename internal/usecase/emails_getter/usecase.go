package emails_getter

import (
	"context"
	"fmt"
	"github.com/NekruzRakhimov/auth_service/internal/config"
	"github.com/NekruzRakhimov/auth_service/internal/port/driven"
)

type UseCase struct {
	cfg         *config.Config
	userStorage driven.UserStorage
}

func New(cfg *config.Config, userStorage driven.UserStorage) *UseCase {
	return &UseCase{
		cfg:         cfg,
		userStorage: userStorage,
	}
}

func (uc *UseCase) GetAll(ctx context.Context) ([]string, error) {
	emails, err := uc.userStorage.GetAllUsersEmails(ctx)
	if err != nil {
		return nil, fmt.Errorf("usecase.GetAllUsersEmails: %w", err)
	}

	return emails, nil
}
