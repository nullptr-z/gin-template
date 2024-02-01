package settings

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	// _ "gorm.io/driver/postgres"
	_ "github.com/lib/pq"
)

var (
	// Pdb *gorm.DB
	db *sqlx.DB
)

func InitSqlDB() error {
	var config DbConfig
	if err := viper.UnmarshalKey("postgres", &config); err != nil {
		zap.L().Error("Error load configure by viper.Unmarshal:", zap.Error(err))
		return err
	}
	connectStr := fmt.Sprintf("postgres://%s:%s@%s:%d?sslmode=disable",
		config.User, config.Password, config.Host, config.Port)
	// db, err := gorm.Open(postgres.Open(connectStr), &gorm.Config{})
	// 使用 sqlx 链接数据库
	if sql, err := sqlx.Connect("postgres", connectStr); err != nil {
		zap.L().Error("sqlx connect:", zap.Error(err))
		return err
	} else {
		db = sql
	}

	db.SetMaxOpenConns(config.MaxOpenConnect)
	db.SetMaxIdleConns(config.MaxIdleConnect)
	fmt.Println("Postgres initialized.......")
	return nil
}

func CloseDB() {
	db.Close()
}

type DbConfig struct {
	Host           string
	Port           int16
	User           string
	Password       string
	DbName         string `mapstructure:"db_name"`
	MaxOpenConnect int    `mapstructure:"max_open_connect"`
	MaxIdleConnect int    `mapstructure:"Max_idle_connect"`
}
