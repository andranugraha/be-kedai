package connection

import (
	"fmt"
	"kedai/backend/be-kedai/config"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	dbConfig = config.DB
	db       *gorm.DB
)

func getLogger() logger.Interface {
	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			LogLevel: logger.Info,
		},
	)
}

func ConnectDB() (err error) {
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable",
		dbConfig.Host,
		dbConfig.Username,
		dbConfig.Password,
		dbConfig.DbName,
		dbConfig.Port,
	)

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: getLogger(),
	})

	return
}

func GetDB() *gorm.DB {
	return db
}
