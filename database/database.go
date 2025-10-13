package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/rendyfutsuy/base-go.git/utils"
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

func setStringConnectionBlips() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		utils.ConfigVars.String("database.blips.host"),
		utils.ConfigVars.Int("database.blips.port"),
		utils.ConfigVars.String("database.blips.user"),
		utils.ConfigVars.String("database.blips.password"),
		utils.ConfigVars.String("database.blips.db_name"),
		utils.ConfigVars.String("database.blips.sslmode"),
	)
}

func ConnectToDB(destinationDB string) *sql.DB {
	var stringConnection string
	if destinationDB == "Blips" {
		stringConnection = setStringConnectionBlips()
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
