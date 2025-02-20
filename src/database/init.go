package database

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"os"
	"smsAlert/logger"
	"strings"
)

type Database struct {
	Driver string
	Source string
}

var config = new(Database)

func Init(path string) (*gorm.DB, error) {
	config, err := initDatabase(path)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database configuration: %w", err)
	}

	db, err := connectMysql(config.Driver, config.Source)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

// 截取settings.yml中需要的子节点，生成数据库连接所需的参数：Driver、Source
func initDatabase(path string) (*Database, error) {
	viper.SetConfigFile(path)

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)

	}

	if err := viper.ReadConfig(strings.NewReader(os.ExpandEnv(string(content)))); err != nil {
		log.Fatal(fmt.Sprintf("Parse config file fail: %s", err.Error()))
	}

	// 读取yaml子节点settings.database
	cfgDatabase := viper.Sub("settings.database")
	if cfgDatabase == nil {
		return nil, fmt.Errorf("missing 'settings.database' configuration")
	}

	// 读取Viper类型的指针变量cfgDatabase，生成Config结构体，结构体中只含有连接所需的Driver和Source字段
	db := &Database{
		Driver: cfgDatabase.GetString("driver"),
		Source: cfgDatabase.GetString("source"),
	}

	if db.Driver == "" || db.Source == "" {
		return nil, fmt.Errorf("invalid database configuration: driver or source is empty")
	}

	return db, nil
}

// 连接db，生成*gorm.DB类型的指针变量
func connectMysql(driver, source string) (*gorm.DB, error) {
	if driver != "mysql" {
		return nil, fmt.Errorf("unsupported database driver: %s", driver)
	}

	db, err := gorm.Open(mysql.Open(source), &gorm.Config{})

	if err != nil {
		logger.Log.Error("failed to connect to database: " + err.Error())
		return nil, fmt.Errorf("failed to connect to MySQL database: %w", err)
	}

	return db, nil
}
