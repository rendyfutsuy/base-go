//	@title			Base Template API Documentation
//	@version		0.0-beta
//	@description	Welcome to the API documentation for the Base Template Web Application. This comprehensive guide is designed to help developers seamlessly integrate and interact with our platform's functionalities. Whether you're building new features, enhancing existing ones, or troubleshooting, this documentation provides all the necessary resources and information.

// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
// @description					Enter JWT token (ex: Bearer eyJhbGciOiJIU....)
package router

import (
	"net/http"
	"time"

	// "github.com/go-playground/validator/v10"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/newrelic/go-agent/v3/integrations/nrecho-v4"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/redis/go-redis/v9"
	_ "github.com/rendyfutsuy/base-go/docs"
	"github.com/rendyfutsuy/base-go/utils"
	"github.com/rendyfutsuy/base-go/utils/services"
	"github.com/rendyfutsuy/base-go/worker"
	echoSwagger "github.com/swaggo/echo-swagger"
	"gorm.io/gorm"

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

	_groupController "github.com/rendyfutsuy/base-go/modules/group/delivery/http"
	_groupRepo "github.com/rendyfutsuy/base-go/modules/group/repository"
	_groupService "github.com/rendyfutsuy/base-go/modules/group/usecase"

	_parameterController "github.com/rendyfutsuy/base-go/modules/parameter/delivery/http"
	_parameterRepo "github.com/rendyfutsuy/base-go/modules/parameter/repository"
	_parameterService "github.com/rendyfutsuy/base-go/modules/parameter/usecase"

	_regencyController "github.com/rendyfutsuy/base-go/modules/regency/delivery/http"
	_regencyRepo "github.com/rendyfutsuy/base-go/modules/regency/repository"
	_regencyService "github.com/rendyfutsuy/base-go/modules/regency/usecase"

	_subGroupController "github.com/rendyfutsuy/base-go/modules/sub-group/delivery/http"
	_subGroupRepo "github.com/rendyfutsuy/base-go/modules/sub-group/repository"
	_subGroupService "github.com/rendyfutsuy/base-go/modules/sub-group/usecase"

	_typeController "github.com/rendyfutsuy/base-go/modules/type/delivery/http"
	_typeRepo "github.com/rendyfutsuy/base-go/modules/type/repository"
	_typeService "github.com/rendyfutsuy/base-go/modules/type/usecase"

	_expeditionController "github.com/rendyfutsuy/base-go/modules/expedition/delivery/http"
	_expeditionRepo "github.com/rendyfutsuy/base-go/modules/expedition/repository"
	_expeditionService "github.com/rendyfutsuy/base-go/modules/expedition/usecase"

	_backingController "github.com/rendyfutsuy/base-go/modules/backing/delivery/http"
	_backingRepo "github.com/rendyfutsuy/base-go/modules/backing/repository"
	_backingService "github.com/rendyfutsuy/base-go/modules/backing/usecase"

	_fileRepo "github.com/rendyfutsuy/base-go/modules/file/repository"
	_fileService "github.com/rendyfutsuy/base-go/modules/file/usecase"
	_postController "github.com/rendyfutsuy/base-go/modules/post/delivery/http"
	_postRepo "github.com/rendyfutsuy/base-go/modules/post/repository"
	_postService "github.com/rendyfutsuy/base-go/modules/post/usecase"
)

func InitializedRouter(gormDB *gorm.DB, redisClient *redis.Client, timeoutContext time.Duration, v *validator.Validate, nrApp *newrelic.Application) *echo.Echo {
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

	// Config Rate Limiter with configurable limit (default 1000 requests/sec)
	throttleMiddleware := authmiddleware.NewThrottleMiddleware()
	router.Use(throttleMiddleware.Throttle())

	router.GET("/", _homepageController.DefaultHomepage)
	router.GET("/health/storage", _homepageController.StorageHealth)

	// Swagger documentation - hanya tersedia di development environment
	if utils.ConfigVars.String("app_env") == "development" {
		router.GET("/swagger/*", echoSwagger.WrapHandler)
	}

	// uses to render files stored on local device.
	router.Static("/storage", "public/storage")

	// Services  ------------------------------------------------------------------------------------------------------------------------------------------------------
	emailServices, err := services.NewEmailService()
	if err != nil {
		panic(err)
	}

	// Initialize the Queue client for Asynq
	q := services.NewQueueService()
	queueClient, err := q.NewAsynqClient()

	// render error if failed to create queue client
	if err != nil {
		panic(err)
	}

	// Repositories ------------------------------------------------------------------------------------------------------------------------------------------------------
	authRepo := _authRepo.NewAuthRepository(gormDB, emailServices, queueClient)   // Using GORM for auth
	roleManagementRepo := _roleManagementRepo.NewRoleManagementRepository(gormDB) // Using GORM for role_management

	userManagementRepo := _userManagementRepo.NewUserManagementRepository(gormDB) // Using GORM for user_management

	groupRepo := _groupRepo.NewGroupRepository(gormDB) // Using GORM for group

	parameterRepo := _parameterRepo.NewParameterRepository(gormDB) // Using GORM for parameter

	regencyRepo := _regencyRepo.NewRegencyRepository(gormDB) // Using GORM for regency

	subGroupRepo := _subGroupRepo.NewSubGroupRepository(gormDB) // Using GORM for sub-group

	typeRepo := _typeRepo.NewTypeRepository(gormDB) // Using GORM for type

	backingRepo := _backingRepo.NewBackingRepository(gormDB) // Using GORM for backing

	expeditionRepo := _expeditionRepo.NewExpeditionRepository(gormDB) // Using GORM for expedition

	postRepo := _postRepo.NewPostRepository(gormDB) // Using GORM for Post
	fileRepo := _fileRepo.NewFileRepository(gormDB) // Using GORM for File

	// Middlewares ------------------------------------------------------------------------------------------------------------------------------------------------------
	middlewareAuth := authmiddleware.NewMiddlewareAuth()
	middlewarePermission := roleMiddleware.NewMiddlewarePermission(
		roleManagementRepo,
	)

	middlewarePageRequest := _reqContext.NewMiddlewarePageRequest()

	// Initialize race condition middleware
	raceConditionMiddleware := authmiddleware.NewRaceConditionMiddleware(redisClient)

	// Example: Add protected routes with race condition prevention
	// This demonstrates how to apply race condition middleware to specific routes
	raceProtectedGroup := router.Group("/v1/protected")
	raceProtectedGroup.Use(middlewareAuth.AuthorizationCheck)
	raceProtectedGroup.Use(raceConditionMiddleware.PreventRaceCondition("protected_operations"))
	raceProtectedGroup.GET("/safe-operation", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "This operation is protected from race conditions"})
	})

	// File usecase (shared by modules)
	fileService := _fileService.NewFileUsecase(fileRepo)

	// Auth
	authService := _authService.NewAuthUsecase(
		authRepo,
		roleManagementRepo,
		timeoutContext,
		utils.ConfigVars.String("jwt_key"),
		[]byte(utils.ConfigVars.String("jwt_key")),
		[]byte(utils.ConfigVars.String("jwt_refresh_key")),
		fileService,
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

	// group management
	groupService := _groupService.NewGroupUsecase(groupRepo)
	_groupController.NewGroupHandler(
		router,
		groupService,
		middlewarePageRequest,
		middlewareAuth,
		middlewarePermission,
	)

	// parameter management
	parameterService := _parameterService.NewParameterUsecase(parameterRepo)
	_parameterController.NewParameterHandler(
		router,
		parameterService,
		middlewarePageRequest,
		middlewareAuth,
		middlewarePermission,
	)

	// regency management
	regencyService := _regencyService.NewRegencyUsecase(regencyRepo)
	_regencyController.NewRegencyHandler(
		router,
		regencyService,
		middlewarePageRequest,
		middlewareAuth,
		middlewarePermission,
	)

	// sub-group management
	subGroupService := _subGroupService.NewSubGroupUsecase(subGroupRepo, groupRepo)
	_subGroupController.NewSubGroupHandler(
		router,
		subGroupService,
		middlewarePageRequest,
		middlewareAuth,
		middlewarePermission,
	)

	// type management
	typeService := _typeService.NewTypeUsecase(typeRepo, subGroupRepo)
	_typeController.NewTypeHandler(
		router,
		typeService,
		middlewarePageRequest,
		middlewareAuth,
		middlewarePermission,
	)

	// backing management
	backingService := _backingService.NewBackingUsecase(backingRepo, typeRepo)
	_backingController.NewBackingHandler(
		router,
		backingService,
		middlewarePageRequest,
		middlewareAuth,
		middlewarePermission,
	)

	// expedition management
	expeditionService := _expeditionService.NewExpeditionUsecase(expeditionRepo)
	_expeditionController.NewExpeditionHandler(
		router,
		expeditionService,
		middlewarePageRequest,
		middlewareAuth,
		middlewarePermission,
	)

	// post management (public index & detail, protected create/update/delete)
	postService := _postService.NewPostUsecase(postRepo, parameterRepo, fileService)
	_postController.NewPostHandler(
		router,
		postService,
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
