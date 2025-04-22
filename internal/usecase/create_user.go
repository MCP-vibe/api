package usecase

import (
	"api/internal/entities"
	"context"
	"errors"
	"time"
)

type CreateUserInput struct {
	TelegramID int64  `json:"telegram_id" validate:"required"`
	Username   string `json:"username" validate:"required"`
	FirstName  string `json:"first_name" validate:"required"`
	LastName   string `json:"last_name" validate:"required"`
}

type CreateUserUseCase struct {
	userRepo   entities.UserRepository
	ctxTimeout time.Duration
}

func NewCreateUserUseCase(
	userRepo entities.UserRepository,
	t time.Duration,
) CreateUserUseCase {
	return CreateUserUseCase{
		userRepo:   userRepo,
		ctxTimeout: t,
	}
}

func (uc CreateUserUseCase) Execute(ctx context.Context, input CreateUserInput) error {
	ctx, cancel := context.WithTimeout(ctx, uc.ctxTimeout)
	defer cancel()

	// check if user already exists
	_, err := uc.userRepo.FindByTelegramID(ctx, input.TelegramID)
	if err != nil {
		if !errors.Is(err, entities.ErrUserNotFound) {
			return err
		}
	} else {
		return entities.ErrUserAlreadyExist
	}

	user := entities.NewUserCreate(
		input.TelegramID,
		input.Username,
		input.FirstName,
		input.LastName,
		true,
		time.Now(),
	)

	_, err = uc.userRepo.Create(ctx, user)
	if err != nil {
		return err
	}

	return nil
}
