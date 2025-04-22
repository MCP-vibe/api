package repo

import (
	"api/internal/entities"
	"api/internal/infrastructure/database"
	"context"
	"time"

	"gorm.io/gorm"
)

type userGORM struct {
	ID         uint32    `gorm:"primary_key"`
	TelegramID int64     `gorm:"column:telegram_id"`
	Username   string    `gorm:"column:username"`
	FirstName  string    `gorm:"column:first_name"`
	LastName   string    `gorm:"column:last_name"`
	IsActive   bool      `gorm:"column:is_active"`
	CreatedAt  time.Time `gorm:"column:created_at"`
}

type UserDB struct {
	dbManager database.DBManager
	tableName string
}

func NewUserDB(dbManager database.DBManager) UserDB {
	return UserDB{
		dbManager: dbManager,
		tableName: "users",
	}
}

func (u UserDB) Create(ctx context.Context, user entities.User) (entities.User, error) {
	var userGORM = userGORM{
		TelegramID: user.TelegramID(),
		Username:   user.Username(),
		FirstName:  user.FirstName(),
		LastName:   user.LastName(),
		IsActive:   user.IsActive(),
		CreatedAt:  user.CreatedAt(),
	}

	if err := u.dbManager.With(ctx).Table(u.tableName).Create(&userGORM).Error; err != nil {
		return entities.User{}, err
	}

	user.SetID(userGORM.ID)
	return user, nil
}

func (u UserDB) Update(ctx context.Context, user entities.User) error {
	updatesMap := map[string]interface{}{
		"username":   user.Username(),
		"first_name": user.FirstName(),
		"last_name":  user.LastName(),
		"is_active":  user.IsActive(),
		"created_at": user.CreatedAt(),
	}

	if err := u.dbManager.With(ctx).Table(u.tableName).Where("id = ?", user.ID()).Updates(updatesMap).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return entities.ErrUserNotFound
		default:
			return err
		}
	}
	return nil
}

func (u UserDB) FindByID(ctx context.Context, id uint32) (entities.User, error) {
	var userGORM userGORM

	if err := u.dbManager.With(ctx).Table(u.tableName).Where("id = ?", id).First(&userGORM).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return entities.User{}, entities.ErrUserNotFound
		}
		return entities.User{}, err
	}

	return entities.NewUser(
		userGORM.ID,
		userGORM.TelegramID,
		userGORM.Username,
		userGORM.FirstName,
		userGORM.LastName,
		userGORM.IsActive,
		userGORM.CreatedAt,
	), nil
}

func (u UserDB) FindByTelegramID(ctx context.Context, telegramID int64) (entities.User, error) {
	var userGORM userGORM

	if err := u.dbManager.With(ctx).Table(u.tableName).Where("telegram_id = ?", telegramID).First(&userGORM).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return entities.User{}, entities.ErrUserNotFound
		}
		return entities.User{}, err
	}

	return entities.NewUser(
		userGORM.ID,
		userGORM.TelegramID,
		userGORM.Username,
		userGORM.FirstName,
		userGORM.LastName,
		userGORM.IsActive,
		userGORM.CreatedAt,
	), nil
}

func (u UserDB) FindAll(ctx context.Context) ([]entities.User, int64, error) {
	var (
		userGORMs []userGORM
		total     int64
	)

	if err := u.dbManager.With(ctx).Table(u.tableName).Count(&total).Find(&userGORMs).Error; err != nil {
		return nil, 0, err
	}

	users := make([]entities.User, 0, len(userGORMs))
	for _, userGORM := range userGORMs {
		users = append(users, entities.NewUser(
			userGORM.ID,
			userGORM.TelegramID,
			userGORM.Username,
			userGORM.FirstName,
			userGORM.LastName,
			userGORM.IsActive,
			userGORM.CreatedAt,
		))
	}
	return users, total, nil
}
