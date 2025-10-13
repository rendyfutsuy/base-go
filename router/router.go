//	@title			Base Template API Documentation
//	@version		0.0-beta
//	@description	Welcome to the API documentation for the Base Template Web Application. This comprehensive guide is designed to help developers seamlessly integrate and interact with our platform's functionalities. Whether you're building new features, enhancing existing ones, or troubleshooting, this documentation provides all the necessary resources and information.

// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
// @description					Enter JWT token (ex: Bearer eyJhbGciOiJIU....)
package router

import (
	"database/sql"
	"net/http"
	"time"

	// "github.com/go-playground/validator/v10"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/newrelic/go-agent/v3/integrations/nrecho-v4"
	"github.com/newrelic/go-agent/v3/newrelic"
	_ "github.com/rendyfutsuy/base-go/docs"
	"github.com/rendyfutsuy/base-go/utils"
	"github.com/rendyfutsuy/base-go/utils/services"
	"github.com/rendyfutsuy/base-go/worker"
	echoSwagger "github.com/swaggo/echo-swagger"

	_homepageController "github.com/rendyfutsuy/base-go/modules/homepage/delivery/http"

	// middleware "github.com/rendyfutsuy/base-go/helpers/middleware"
	_reqContext "github.com/rendyfutsuy/base-go/helpers/middleware/request"
	// "github.com/rendyfutsuy/base-go/helpers/validations"

	_authController "github.com/rendyfutsuy/base-go/modules/auth/delivery/http"
	_authRepo "github.com/rendyfutsuy/base-go/modules/auth/repository"
	_authService "github.com/rendyfutsuy/base-go/modules/auth/usecase"

	authmiddleware "github.com/rendyfutsuy/base-go/helpers/middleware"
	roleMiddleware "github.com/rendyfutsuy/base-go/helpers/middleware"

	_userManagementController "github.com/rendyfutsuy/base-go/modules/user_management/delivery/http"
	_userManagementRepo "github.com/rendyfutsuy/base-go/modules/user_management/repository"
	_userManagementService "github.com/rendyfutsuy/base-go/modules/user_management/usecase"

	_roleManagementController "github.com/rendyfutsuy/base-go/modules/role_management/delivery/http"
	_roleManagementRepo "github.com/rendyfutsuy/base-go/modules/role_management/repository"
	_roleManagementService "github.com/rendyfutsuy/base-go/modules/role_management/usecase"
)

func InitializedRouter(db *sql.DB, timeoutContext time.Duration, v *validator.Validate, nrApp *newrelic.Application) *echo.Echo {
	router := echo.New()

	// queries := sqlc.New(db)

	// Config CORS
	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:          middleware.DefaultSkipper,
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, echo.HeaderXCSRFToken},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodOptions},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	router.Use(middleware.Logger())
	router.Use(middleware.Recover())
	router.Use(nrecho.Middleware(nrApp))

	// Config Rate Limiter allows 100 requests/sec
	router.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(100)))

	// Config Validator to Router
	// router.Validator = &utils.RequestValidator{Validator: validator.New()}
	// structValidator := validator.New()

	// Register RequestLog to Router Middleware
	// router.Use(utils.RequestLog)

	// Register HTTP Error Handler function
	// router.HTTPErrorHandler = helper.ErrorHandler

	router.GET("/", _homepageController.DefaultHomepage)

	router.GET("/swagger/*", echoSwagger.WrapHandler)
	// Services  ------------------------------------------------------------------------------------------------------------------------------------------------------
	emailServices, err := services.NewEmailService()
	if err != nil {
		panic(err)
	}

	// Repositories ------------------------------------------------------------------------------------------------------------------------------------------------------
	authRepo := _authRepo.NewAuthRepository(db, emailServices)
	roleManagementRepo := _roleManagementRepo.NewRoleManagementRepository(db)

	userManagementRepo := _userManagementRepo.NewUserManagementRepository(db)

	// Middlewares ------------------------------------------------------------------------------------------------------------------------------------------------------
	middlewareAuth := authmiddleware.NewMiddlewareAuth(authRepo)
	middlewarePermission := roleMiddleware.NewMiddlewarePermission(
		authRepo,
		roleManagementRepo,
	)

	middlewarePageRequest := _reqContext.NewMiddlewarePageRequest()

	// Auth
	authService := _authService.NewAuthUsecase(
		authRepo,
		timeoutContext,
		utils.ConfigVars.String("jwt_key"),
		[]byte(utils.ConfigVars.String("jwt_key")),
	)
	_authController.NewAuthHandler(
		router,
		authService,
		middlewareAuth,
		middlewarePageRequest,
	)

	// role management
	roleManagementService := _roleManagementService.NewRoleManagementUsecase(
		roleManagementRepo,
		authRepo,
		timeoutContext,
	)
	_roleManagementController.NewRoleManagementHandler(
		router,
		roleManagementService,
		middlewarePageRequest,
		middlewareAuth,
		middlewarePermission,
	)

	// user management
	userManagementService := _userManagementService.NewUserManagementUsecase(
		userManagementRepo,
		roleManagementRepo,
		authRepo,
		timeoutContext,
	)
	_userManagementController.NewUserManagementHandler(
		router,
		userManagementService,
		middlewarePageRequest,
		middlewareAuth,
		middlewarePermission,
	)

	usecaseRegistry := worker.UsecaseRegistry{
		// Add any other usecases that your background jobs might need
	}

	dispatcher := worker.NewDispatcher(10, usecaseRegistry) // Using 10 workers, for example
	dispatcher.Run()

	time.Sleep(1000 * time.Millisecond)
	return router
}
