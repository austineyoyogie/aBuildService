package database

import (
	"aIBuildService/aPI/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var lc = config.LoadConfig()

func CDriver() *gorm.DB {
	db, err := gorm.Open(mysql.Open(lc.DBC.Username+":"+lc.DBC.Password+
		"@tcp("+lc.DBC.Hostname+":"+lc.DBC.Port+")/"+lc.DBC.Database+
		"?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		logFatal(err)
		return nil
	}
	return db
}

func logFatal(err error) {
	if err != nil {
		log.Fatal("\nError database connection mode.")
	}
}
