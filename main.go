package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/rendyfutsuy/base-go/database"
	"github.com/rendyfutsuy/base-go/router"
	"github.com/rendyfutsuy/base-go/utils"
)

var (
	app struct {
		DatabaseBlips *sql.DB
		Router        *echo.Echo
		// QueueClient          *asynq.Client
	}
)

func init() {
	utils.InitConfig("config.json")
	utils.InitializedLogger()
	log.Println("Starting service on port", utils.ConfigVars.String("app_port"))

	// Trying to connect to the database
	app.DatabaseBlips = database.ConnectToDB("Blips")
	if app.DatabaseBlips == nil {
		panic("Can't connect to Postgres : Blips!")
	}
}

func main() {
	utils.Logger.Info("Start the app")

	// Set a timeout for each endpoint
	timeoutContext := time.Duration(utils.ConfigVars.Int("context.timeout")) * time.Second

	app.Router = router.InitializedRouter(app.DatabaseBlips, timeoutContext)

	app.Router.Logger.Fatal(app.Router.Start(fmt.Sprintf(":%s", utils.ConfigVars.String("app_port"))))

	defer func() {
		if app.DatabaseBlips != nil {
			app.DatabaseBlips.Close()
		}
	}()
}
