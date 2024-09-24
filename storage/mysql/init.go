package mysql

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
)

var DB *gorm.DB

type DBConfig struct {
	User            string        `yaml:"user"`
	Password        string        `yaml:"password"`
	Host            string        `yaml:"host"`
	Port            int           `yaml:"port"`
	DBName          string        `yaml:"dbname"`
	Charset         string        `yaml:"charset"`
	ParseTime       bool          `yaml:"parseTime"`
	Loc             string        `yaml:"loc"`
	MaxIdleConns    int           `yaml:"maxIdleConns"`
	MaxOpenConns    int           `yaml:"maxOpenConns"`
	ConnMaxLifetime time.Duration `yaml:"connMaxLifetime"`
	LogLevel        string        `yaml:"logLevel"`
}

func (config *DBConfig) GetLogLevel() logger.LogLevel {
	switch config.LogLevel {
	case "error":
		return logger.Error
	case "warn":
		return logger.Warn
	case "info":
		return logger.Info
	default:
		return logger.Info
	}
}

func InitDB(configPath string) {
	viper.SetConfigFile(configPath)
	err := viper.ReadInConfig()
	if err != nil {
		log.Panicf("Error reading config file, %s", err)
	}
	var config DBConfig
	if err := viper.UnmarshalKey("db", &config); err != nil {
		log.Panicf("Unable to decode into struct, %v", err)
	}

	// 解析ConnMaxLifetime
	duration, err := time.ParseDuration(viper.GetString("db.connMaxLifetime"))
	if err != nil {
		log.Panicf("Invalid ConnMaxLifetime format, %v", err)
	}
	config.ConnMaxLifetime = duration
	db, err := initDB(&config)
	if err != nil {
		log.Panicf("Failed to connect database, %v", err)
	}
	DB = db

}

func initDB(config *DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.DBName,
		config.Charset,
		config.ParseTime,
		config.Loc,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.LogLevel(config.GetLogLevel())),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)

	return db, nil
}
