package db

import (
  "gorm.io/driver/postgres"
  "gorm.io/gorm"
  "sync"
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
	  })
	  return instance.DB
}