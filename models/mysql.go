package models

import (
	"fmt"
	"os"

	"github.com/toolkits/pkg/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

type MysqlSection struct {
	Addr     string `yaml:"addr"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Charset  string `yaml:"charset"`
	Database string `ymal:"database"`
	Max      int    `yaml:"max"`
	Idle     int    `yaml:"idle"`
}

var MySQL MysqlSection

func InitMySQL(MySQL MysqlSection) {
	conf := MySQL

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=true&loc=Local", conf.Username,
		conf.Password, conf.Addr, conf.Database, conf.Charset)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("cannot connect mysql[%s]: %v", conf.Addr, err)
		os.Exit(1)
	}

	sqlDB, err := db.DB()
	if err != nil {
		fmt.Printf("cannot connect mysql[%s]: %v", conf.Addr, err)
		os.Exit(1)
	}
	sqlDB.SetMaxIdleConns(conf.Idle)
	sqlDB.SetMaxOpenConns(conf.Max)

	DB = db
}

func DBInsertOne(bean interface{}) error {
	err := DB.Create(bean).Error
	if err != nil {
		logger.Errorf("mysql.error: insert fail: %v, to insert object: %+v", err, bean)
		return internalServerError
	}

	return nil
}
