package db

import (
  "github.com/glebarez/sqlite"
  "gorm.io/gorm"
  "sync"
)

var instance struct {	
	DB *gorm.DB
	once sync.Once
}

func GetDB() *gorm.DB {
	  instance.once.Do(func() {
		db, err := gorm.Open(sqlite.Open("gym.db"), &gorm.Config{})
		if err != nil {
			panic("failed to connect database" + err.Error())
		}
		instance.DB = db
	  })
	  return instance.DB
}