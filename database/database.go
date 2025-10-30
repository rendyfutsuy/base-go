package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/rendyfutsuy/base-go/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var counts int64

func openDB(psqlInfo string) (*sql.DB, error) {
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func setStringConnectionDatabase() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		utils.ConfigVars.String("database.host"),
		utils.ConfigVars.Int("database.port"),
		utils.ConfigVars.String("database.user"),
		utils.ConfigVars.String("database.password"),
		utils.ConfigVars.String("database.db_name"),
		utils.ConfigVars.String("database.sslmode"),
	)
}

func ConnectToDB(destinationDB string) *sql.DB {
	var stringConnection string
	if destinationDB == "Database" {
		stringConnection = setStringConnectionDatabase()
	}

	for {
		connection, err := openDB(stringConnection)
		if err != nil {
			utils.Logger.Error("Postgres not yet ready... : " + destinationDB)
			utils.Logger.Error(err.Error())
			counts++
		} else {
			utils.Logger.Info("Connected to Postgres : " + destinationDB)
			connection.SetMaxOpenConns(100)
			connection.SetMaxIdleConns(25)
			connection.SetConnMaxLifetime(5 * time.Minute)
			return connection
		}

		if counts > 10 {
			return nil
		}

		log.Println("backing off for two seconds...")
		time.Sleep(2 * time.Second)
		continue
	}
}

// ConnectToGORM creates a GORM database connection
func ConnectToGORM(destinationDB string) *gorm.DB {
	var stringConnection string
	if destinationDB == "Database" {
		stringConnection = setStringConnectionDatabase()
	}

	for {
		gormDB, err := gorm.Open(postgres.Open(stringConnection), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
			NowFunc: func() time.Time {
				return time.Now().UTC()
			},
			PrepareStmt: true,
		})

		if err != nil {
			utils.Logger.Error("Postgres not yet ready (GORM)... : " + destinationDB)
			utils.Logger.Error(err.Error())
			counts++
		} else {
			utils.Logger.Info("Connected to Postgres (GORM) : " + destinationDB)

			// Get underlying SQL database
			sqlDB, err := gormDB.DB()
			if err != nil {
				utils.Logger.Error("Failed to get underlying SQL DB: " + err.Error())
				counts++
				continue
			}

			// Set connection pool settings
			sqlDB.SetMaxOpenConns(100)
			sqlDB.SetMaxIdleConns(25)
			sqlDB.SetConnMaxLifetime(5 * time.Minute)

			return gormDB
		}

		if counts > 10 {
			return nil
		}

		log.Println("backing off for two seconds...")
		time.Sleep(2 * time.Second)
		continue
	}
}
