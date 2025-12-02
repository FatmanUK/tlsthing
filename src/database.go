package main

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/driver/postgres"
)

func GetDatabaseConnection(
		dsn string,
		logs chan Log) (*gorm.DB, error) {
	var err error
	var db *gorm.DB
	dialector := postgres.Open(dsn)
	gormConfig := gorm.Config{
		Logger: logger.Discard,
	}
	count := 0
	dbConnectMsg := "Connecting database (attempt %d)..."
	for ok := true; ok; ok = (err != nil && count < 9) {
		DoublingSleep(count, 2)
		logs <- Info(fmt.Sprintf(dbConnectMsg, count + 1))
		db, err = gorm.Open(dialector, &gormConfig)
		if err == nil {
			break
		}
		count++
	}
	if err != nil {
		logs <- Error("Connection failed.")
		return nil, err
	}
	logs <- Info("Connection established.")
	db.AutoMigrate(&RegistryDatum{})
	return db, nil
}
