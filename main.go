package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/redis/go-redis/v9"
	"github.com/rendyfutsuy/base-go/database"
	"github.com/rendyfutsuy/base-go/router"
	"github.com/rendyfutsuy/base-go/utils"
	"gorm.io/gorm"
)

var (
	app struct {
		Database    *sql.DB
		GormDB      *gorm.DB
		RedisClient *redis.Client
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

	// Connect to database with raw SQL (keeping for backward compatibility)
	app.Database = database.ConnectToDB("Database")
	if app.Database == nil {
		panic("Can't connect to Postgres : Database!")
	}

	// Connect to database with GORM
	app.GormDB = database.ConnectToGORM("Database")
	if app.GormDB == nil {
		panic("Can't connect to Postgres with GORM : Database!")
	}

	// Connect to Redis
	app.RedisClient = database.ConnectToRedis()
	if app.RedisClient == nil {
		utils.Logger.Warn("Could not connect to Redis, continuing without Redis support")
	}

	app.Validator = validator.New()
	utils.RegisterCustomValidator(app.Validator)
}

func main() {

	utils.Logger.Info("Start the app")

	// Set a timeout for each endpoint
	timeoutContext := time.Duration(utils.ConfigVars.Int("context.timeout")) * time.Second

	app.Router = router.InitializedRouter(app.Database, app.GormDB, app.RedisClient, timeoutContext, app.Validator, app.NewRelicApp)

	app.Router.Validator = &utils.CustomValidator{Validator: app.Validator}

	app.Router.Logger.Fatal(app.Router.Start(fmt.Sprintf(":%s", utils.ConfigVars.String("app_port"))))

	defer func() {
		if app.Database != nil {
			app.Database.Close()
		}
		if app.GormDB != nil {
			sqlDB, _ := app.GormDB.DB()
			if sqlDB != nil {
				sqlDB.Close()
			}
		}
		if app.RedisClient != nil {
			app.RedisClient.Close()
		}
	}()
}
