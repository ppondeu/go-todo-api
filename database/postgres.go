package database

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ppondeu/go-todo-api/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type postgresDatabase struct {
	Db *gorm.DB
}

var (
	once       sync.Once
	dbInstance *postgresDatabase
)

type sqlLogger struct {
	logger.Interface
}

func (s sqlLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sqlString, _ := fc()
	fmt.Printf("\n===============================\n%v\n===============================\n", sqlString)
}

func NewPostgresDatabase(conf *config.Config) Database {
	once.Do(func() {
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
			conf.Database.Host,
			conf.Database.User,
			conf.Database.Password,
			conf.Database.DBName,
			conf.Database.Port,
			conf.Database.SSLMode,
			conf.Database.TimeZone,
		)

		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: sqlLogger{logger.Default.LogMode(logger.Info)},
		})
		if err != nil {
			panic("failed to connect database")
		}

		db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")

		dbInstance = &postgresDatabase{Db: db}
	})

	return dbInstance
}

func (p *postgresDatabase) GetDb() *gorm.DB {
	return dbInstance.Db
}
