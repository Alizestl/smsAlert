package database

import (
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"smsAlert/logger"
)

type Database struct {
	Driver string
	Source string
}

var Config = new(Database)
var DB *gorm.DB

func Init(path string) {
	// 截取settings.yml中需要的子节点，生成包内全局变量cfgDatabase
	GenerateDatabaseConfig(path)

	// 读取Viper类型的指针变量cfgDatabase，生成Config结构体，结构体中只含有连接所需的Driver和Source字段
	Config = initDatabase(cfgDatabase)

	// 连接db，生成*gorm.DB类型的指针变量
	DB = connectMysql(Config.Driver, Config.Source)
}

func initDatabase(cfg *viper.Viper) *Database {
	db := &Database{
		Driver: cfg.GetString("driver"),
		Source: cfg.GetString("source"),
	}
	return db
}

func connectMysql(driver, source string) *gorm.DB {
	if driver == "mysql" {
		db, err := gorm.Open(mysql.Open(source), &gorm.Config{})
		if err != nil {
			logger.Log.Fatal("连接数据库失败, error=" + err.Error())
			panic("连接数据库失败, error=" + err.Error())
		}
		return db
	} else {
		panic("数据库类型有误！")
	}
}
