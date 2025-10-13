package router

import (
	"database/sql"
	"net/http"
	"time"

	// "github.com/go-playground/validator/v10"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	// middleware "github.com/rendyfutsuy/base-go/helper/middleware"
	_reqContext "github.com/rendyfutsuy/base-go/helper/middleware/request"
	// "github.com/rendyfutsuy/base-go/helper/validations"

	_roleController "github.com/rendyfutsuy/base-go/modules/role/delivery/http"
	_roleRepo "github.com/rendyfutsuy/base-go/modules/role/repository"
	_roleService "github.com/rendyfutsuy/base-go/modules/role/usecase"
	"github.com/rendyfutsuy/base-go/utils"
	"github.com/rendyfutsuy/base-go/utils/services"

	_authController "github.com/rendyfutsuy/base-go/modules/auth/delivery/http"
	_authRepo "github.com/rendyfutsuy/base-go/modules/auth/repository"
	_authService "github.com/rendyfutsuy/base-go/modules/auth/usecase"

	authmiddleware "github.com/rendyfutsuy/base-go/helpers/middleware"
	_accountController "github.com/rendyfutsuy/base-go/modules/account/delivery/http"
	_accountRepo "github.com/rendyfutsuy/base-go/modules/account/repository"
	_accountService "github.com/rendyfutsuy/base-go/modules/account/usecase"
)

func InitializedRouter(dbBlips *sql.DB, timeoutContext time.Duration) *echo.Echo {
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
	router.Use(middleware.Recover())

	// Config Rate Limiter allows 100 requests/sec
	router.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(100)))

	// Config Validator to Router
	// router.Validator = &utils.RequestValidator{Validator: validator.New()}
	// structValidator := validator.New()

	// Register RequestLog to Router Middleware
	// router.Use(utils.RequestLog)

	// Register HTTP Error Handler function
	// router.HTTPErrorHandler = helper.ErrorHandler

	router.GET("/", func(ec echo.Context) error {
		return ec.JSON(http.StatusOK, map[string]string{
			"message": "Default",
		})
	})
	// Services  ------------------------------------------------------------------------------------------------------------------------------------------------------
	emailServices, err := services.NewEmailService()
	if err != nil {
		panic(err)
	}

	// Repositories ------------------------------------------------------------------------------------------------------------------------------------------------------
	roleRepo := _roleRepo.NewRoleRepository(dbBlips)
	authRepo := _authRepo.NewAuthRepository(dbBlips, emailServices)
	accountRepo := _accountRepo.NewAccountRepository(dbBlips)

	// Middlewares ------------------------------------------------------------------------------------------------------------------------------------------------------
	middlewareAuth := authmiddleware.NewMiddlewareAuth(authRepo)
	middlewarePageRequest := _reqContext.NewMiddlewarePageRequest()

	//Roles
	roleService := _roleService.NewRoleUsecase(
		roleRepo,
		// dbValidations,
		timeoutContext,
	)
	_roleController.NewRoleHandler(
		router,
		roleService,
		middlewarePageRequest,
		middlewareAuth,
	)

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
	)

	// account
	accountService := _accountService.NewAccountUsecase(
		accountRepo,
		timeoutContext,
	)
	_accountController.NewAccountHandler(
		router,
		accountService,
		middlewarePageRequest,
		middlewareAuth,
	)

	return router

}
