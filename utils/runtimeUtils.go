package utils

import (
	"context"
	"gorm.io/gorm"
	"log"
)

// SafeGo 捕获panic
func SafeGo(fn func()) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Recovered from panic")
		}
	}()
	fn()
}

// WithTx 事务context控制
func WithTx(ctx context.Context, db *gorm.DB, fn func(tx *gorm.DB) error) error {
	tx := db.WithContext(ctx).Begin()
	committed := false

	defer func() {
		if !committed {
			tx.Rollback()
		}
		if r := recover(); r != nil {
			panic(r)
		}
	}()

	if err := fn(tx); err != nil {
		return err
	}
	if err := tx.Commit().Error; err != nil {
		return err
	}
	committed = true
	return nil
}
