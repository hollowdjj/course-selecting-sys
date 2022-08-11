package conf

import (
	"log"

	"github.com/go-ini/ini"
)

type App struct {
	RunMode  string
	Host     string
	MainHost string
}

type Logger struct {
	MaxSize    int
	MaxBackups int
	MaxAge     int
}

type Server struct {
	HttpPort     int
	ReadTimeout  int
	WriteTimeout int
}

type Db struct {
	User     string
	Password string
	Host     string
	Name     string
}

var (
	config *ini.File
	app    App
	logger Logger
	server Server
	db     Db
)

//load config.ini
func LoadConfig(path string) {
	var err error
	config, err = ini.Load(path)
	if err != nil {
		log.Fatalf("load config file [%s] error: %v", path, err)
	}
	mapTo("app", app)
	mapTo("logger", logger)
	mapTo("server", server)
	mapTo("db", db)
}

//map .ini file's section to a go struct
func mapTo(section string, v interface{}) {
	err := config.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("map section [%s] fail: %v", section, err)
	}
}

//return a copy of conf.app
func GetApp() App {
	return app
}

//return a copy of conf.logger
func GetLogger() Logger {
	return logger
}

//return a copy of conf.server
func GetServer() Server {
	return server
}

//return a copy of conf.db
func GetDb() Db {
	return db
}
