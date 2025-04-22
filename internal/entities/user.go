package entities

import (
	"api/internal/errors"
	"context"
	"net/http"
	"time"
)

var (
	ErrUserNotFound     = errors.NewError(http.StatusNotFound, "user_not_found")
	ErrUserAlreadyExist = errors.NewError(http.StatusConflict, "user_already_exist")
)

type UserRepository interface {
	Create(ctx context.Context, user User) (User, error)
	Update(ctx context.Context, user User) error
	FindByID(ctx context.Context, id uint32) (User, error)
	FindByTelegramID(ctx context.Context, telegramID int64) (User, error)
	FindAll(ctx context.Context) ([]User, int64, error)
}

type User struct {
	id         uint32
	telegramID int64
	username   string
	firstName  string
	lastName   string
	isActive   bool
	createdAt  time.Time
}

func NewUser(
	id uint32,
	telegramID int64,
	username string,
	firstName string,
	lastName string,
	isActive bool,
	createdAt time.Time,
) User {
	return User{
		id:         id,
		telegramID: telegramID,
		username:   username,
		firstName:  firstName,
		lastName:   lastName,
		isActive:   isActive,
		createdAt:  createdAt,
	}
}

func NewUserCreate(
	telegramID int64,
	username string,
	firstName string,
	lastName string,
	isActive bool,
	createdAt time.Time,
) User {
	return User{
		telegramID: telegramID,
		username:   username,
		firstName:  firstName,
		lastName:   lastName,
		isActive:   isActive,
		createdAt:  createdAt,
	}
}

func (u User) ID() uint32 {
	return u.id
}

func (u User) TelegramID() int64 {
	return u.telegramID
}

func (u User) Username() string {
	return u.username
}

func (u User) FirstName() string {
	return u.firstName
}

func (u User) LastName() string {
	return u.lastName
}

func (u User) IsActive() bool {
	return u.isActive
}

func (u User) CreatedAt() time.Time {
	return u.createdAt
}

// -----------------------------------------------------------------------------

func (u *User) SetID(id uint32) {
	u.id = id
}

func (u *User) SetFirstName(firstName string) {
	u.firstName = firstName
}

func (u *User) SetLastName(lastName string) {
	u.lastName = lastName
}

func (u *User) SetIsActive(isActive bool) {
	u.isActive = isActive
}
