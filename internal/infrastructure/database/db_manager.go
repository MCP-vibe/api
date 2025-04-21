package database

import (
	"context"

	"gorm.io/gorm"
)

type DBManager struct {
	db *gorm.DB
}

type contextKey int

const txKey contextKey = 1

func NewDBManager(db *gorm.DB) *DBManager {
	return &DBManager{db: db}
}

func (tm *DBManager) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	tx := tm.db.Begin()
	ctx = context.WithValue(ctx, txKey, tx)
	if err := fn(ctx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (tm *DBManager) With(ctx context.Context) *gorm.DB {
	tx := ctx.Value(txKey)
	if tx != nil {
		return tx.(*gorm.DB)
	}
	return tm.db
}
