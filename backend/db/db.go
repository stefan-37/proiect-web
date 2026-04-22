package db

import (
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var instance struct {	
	DB *gorm.DB
	once sync.Once
}

func GetDB() *gorm.DB {
	  instance.once.Do(func() {
		db, err := gorm.Open(postgres.Open("postgresql://postgres:password@postgres:5432/mydatabase?sslmode=disable"), &gorm.Config{})
		if err != nil {
			panic("failed to connect database" + err.Error())
		}
		instance.DB = db

		sqlDB, err := db.DB()
		if err != nil {
			panic(err)
		}
		sqlDB.SetMaxOpenConns(20)
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetConnMaxLifetime(5 * time.Minute)

	  })
	  return instance.DB
}