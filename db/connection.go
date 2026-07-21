package db

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/ppondeu/go-todo-api/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type database struct {
	db *gorm.DB
}

var (
	once       sync.Once
	dbInstance *database
)

func NewDatabase(conf *config.Config) DB {
	once.Do(func() {
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
			conf.DB.Host,
			conf.DB.User,
			conf.DB.Password,
			conf.DB.DBName,
			conf.DB.Port,
			conf.DB.SSLMode,
			conf.DB.TimeZone,
		)

		gormLogger := logger.New(log.New(os.Stdout, "", log.LstdFlags), logger.Config{
			LogLevel:             logger.Warn,
			ParameterizedQueries: true,
		})
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger:         gormLogger,
			TranslateError: true,
		})
		if err != nil {
			panic("failed to connect database")
		}

		if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
			panic(fmt.Sprintf("failed to create uuid extension: %v", err))
		}

		dbInstance = &database{db: db}
	})
	return dbInstance
}

func (p *database) GetInstance() *gorm.DB {
	return p.db
}
