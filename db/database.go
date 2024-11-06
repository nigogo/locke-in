package db

import (
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	instance *gorm.DB
	once     sync.Once
)

func GetDB() *gorm.DB {
	once.Do(func() {
		db, err := gorm.Open(sqlite.Open("locke-in.db"), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}
		instance = db
	})
	return instance
}
