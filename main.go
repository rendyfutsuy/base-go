package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rendyfutsuy/base-go/database"
	"github.com/rendyfutsuy/base-go/router"
	"github.com/rendyfutsuy/base-go/utils"
)

var (
	app struct {
		Database    *sql.DB
		Router      *echo.Echo
		NewRelicApp *newrelic.Application
		Validator   *validator.Validate
		// QueueClient          *asynq.Client
	}
)

func init() {
	utils.InitConfig("config.json")
	if utils.ConfigVars.Exists("newrelic.enable_new_relic_logging") {
		if utils.ConfigVars.Bool("newrelic.enable_new_relic_logging") {
			app.NewRelicApp = utils.InitializeNewRelic()
		}
	}

	utils.InitializedLogger(app.NewRelicApp)

	if err := app.NewRelicApp.WaitForConnection(5 * time.Second); nil != err {
		fmt.Println(err)
	}
	// Trying to connect to the database
	app.Database = database.ConnectToDB("Blips")
	if app.Database == nil {
		panic("Can't connect to Postgres : Blips!")
	}

	app.Validator = validator.New()
	utils.RegisterCustomValidator(app.Validator)
}

func main() {

	utils.Logger.Info("Start the app")

	// Set a timeout for each endpoint
	timeoutContext := time.Duration(utils.ConfigVars.Int("context.timeout")) * time.Second

	app.Router = router.InitializedRouter(app.Database, timeoutContext, app.Validator, app.NewRelicApp)

	app.Router.Validator = &utils.CustomValidator{Validator: app.Validator}

	app.Router.Logger.Fatal(app.Router.Start(fmt.Sprintf(":%s", utils.ConfigVars.String("app_port"))))

	defer func() {
		if app.Database != nil {
			app.Database.Close()
		}
	}()
}
