package configs

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDataBase() *gorm.DB {
	var err error
	dbConnURL, _ := ConnectionURLBuilder(GetEnv("DB_DRIVER", "postgres"))

	var loggerInfo logger.Interface

	if GetEnv("ENVIRONMENT", "development") == "development" {
		loggerInfo = logger.Default.LogMode(logger.Info)
	}
	switch GetEnv("DB_DRIVER", "postgres") {
	case "mysql":
		DB, err = gorm.Open(mysql.Open(dbConnURL), &gorm.Config{
			SkipDefaultTransaction:                   true,
			PrepareStmt:                              true,
			DisableForeignKeyConstraintWhenMigrating: true,
			Logger:                                   loggerInfo,
		})
	case "postgres":
		DB, err = gorm.Open(postgres.Open(dbConnURL), &gorm.Config{
			SkipDefaultTransaction: true,
			PrepareStmt:            true,
			Logger:                 loggerInfo,
		})
	}
	if err != nil {
		panic(err.Error())
	}

	return DB
}

func ConnectionURLBuilder(n string) (string, error) {
	// Define URL to connection.
	var url string

	// Switch given names.
	switch n {
	case "mysql":
		// URL for mysql connection.
		url = fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s",
			AppEnv.DBUser,     // GetEnv("DB_USER", "root"),
			AppEnv.DBPassword, // GetEnv("DB_PASSWORD", ""),
			AppEnv.DBHost,     // GetEnv("DB_HOST", "locahost"),
			AppEnv.DBPort,     // GetEnv("DB_PORT", "3306"),
			AppEnv.DBName,     // GetEnv("DB_NAME", "db_name"),
		)
	case "postgres":
		// URL for postgres connection.
		url = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			AppEnv.DBHost,     // GetEnv("DB_HOST", "localhost"),
			AppEnv.DBUser,     // GetEnv("DB_USER", "root"),
			AppEnv.DBPassword, // GetEnv("DB_PASSWORD", ""),
			AppEnv.DBName,     // GetEnv("DB_NAME", "db_name"),
			AppEnv.DBPort,     // GetEnv("DB_PORT", "5432"),
		)
	default:
		// Return error message.
		return "", fmt.Errorf("connection name '%v' is not supported", n)
	}

	// Return connection URL.
	return url, nil
}
