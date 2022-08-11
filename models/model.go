package models

import (
	"fmt"

	"github.com/hollowdjj/course-selecting-sys/conf"
	"github.com/hollowdjj/course-selecting-sys/pkg/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var Db *gorm.DB //database

//connect database
func InitDb() {
	var err error
	dbConf := conf.GetDb()
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbConf.User,
		dbConf.Password,
		dbConf.Host,
		dbConf.Name)
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		logger.GetInstance().Fatalf("connect to database fail")
	}
	mysqlDB, _ := Db.DB()
	mysqlDB.SetMaxIdleConns(10)
	mysqlDB.SetMaxOpenConns(100)
}

func CloseDB() {
	mysqlDB, _ := Db.DB()
	mysqlDB.Close()
}
