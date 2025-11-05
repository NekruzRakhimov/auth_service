package dbstore

import (
	"context"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"

	"github.com/NekruzRakhimov/auth_service/internal/domain"
)

type UserStorage struct {
	db *sqlx.DB
}

func NewUserStorage(db *sqlx.DB) *UserStorage {
	return &UserStorage{db: db}
}

type User struct {
	ID        int       `db:"id"`
	FullName  string    `db:"full_name"`
	Username  string    `db:"username"`
	Password  string    `db:"password"`
	Role      string    `db:"role"`
	Email     string    `db:"email"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (u *User) ToDomain() *domain.User {
	return &domain.User{
		ID:        u.ID,
		FullName:  u.FullName,
		Username:  u.Username,
		Password:  u.Password,
		Role:      domain.Role(u.Role),
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func (u *User) FromDomain(d domain.User) {
	u.ID = d.ID
	u.FullName = d.FullName
	u.Username = d.Username
	u.Password = d.Password
	u.Role = string(d.Role)
	u.Email = d.Email
	u.UpdatedAt = d.UpdatedAt
	u.CreatedAt = d.CreatedAt
}

func (u *UserStorage) CreateUser(ctx context.Context, user domain.User) (err error) {
	var dbUser User
	dbUser.FromDomain(user)

	logger := zerolog.New(os.Stdout).With().Timestamp().Str("func_name", "repository.CreateUser").Logger()
	_, err = u.db.ExecContext(ctx, `INSERT INTO users (full_name, username, password, role, email)
					VALUES ($1, $2, $3, $4, $5)`,
		dbUser.FullName,
		dbUser.Username,
		dbUser.Password,
		dbUser.Role,
		dbUser.Email)
	if err != nil {
		logger.Err(err).Msg("error inserting user")
		return u.translateError(err)
	}

	return nil
}

func (u *UserStorage) GetUserByID(ctx context.Context, id int) (domain.User, error) {
	logger := zerolog.New(os.Stdout).With().Timestamp().Str("func_name", "repository.GetUserByID").Logger()

	var dbUser User
	if err := u.db.GetContext(ctx, &dbUser, `
		SELECT id, full_name, username, password, role, created_at, updated_at, email
		FROM users
		WHERE id = $1`, id); err != nil {
		logger.Err(err).Msg("error selecting user")
		return domain.User{}, u.translateError(err)
	}

	return *dbUser.ToDomain(), nil
}

func (u *UserStorage) GetAllUsersEmails(ctx context.Context) (emails []string, err error) {
	logger := zerolog.New(os.Stdout).With().Timestamp().Str("func_name", "repository.GetUserByID").Logger()

	if err = u.db.SelectContext(ctx, &emails, `
		SELECT email 
		FROM users WHERE enable_notifications = true`); err != nil {
		logger.Err(err).Msg("error selecting user")
		return nil, u.translateError(err)
	}

	return emails, nil
}

func (u *UserStorage) GetUserByUsername(ctx context.Context, username string) (domain.User, error) {
	logger := zerolog.New(os.Stdout).With().Timestamp().Str("func_name", "repository.GetUserByUsername").Logger()
	var dbUser User

	if err := u.db.GetContext(ctx, &dbUser, `
		SELECT id, full_name, username, password, role, created_at, updated_at, email
		FROM users
		WHERE username = $1`, username); err != nil {
		logger.Err(err).Msg("error selecting user")
		return domain.User{}, u.translateError(err)
	}

	return *dbUser.ToDomain(), nil
}
