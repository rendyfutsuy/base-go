package router

import (
	"database/sql"
	"net/http"
	"time"

	// "github.com/go-playground/validator/v10"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	// middleware "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/helper/middleware"
	_reqContext "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/helper/middleware/request"
	// "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/helper/validations"

	_roleController "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/role/delivery/http"
	_roleRepo "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/role/repository"
	_roleService "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/role/usecase"
	_shipyardController "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/shipyard/delivery/http"
	_shipyardRepo "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/shipyard/repository"
	_shipyardService "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/shipyard/usecase"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/utils"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/utils/services"

	_authController "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/auth/delivery/http"
	_authRepo "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/auth/repository"
	_authService "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/auth/usecase"

	authmiddleware "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/helpers/middleware"
	_classController "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/class/delivery/http"
	_classRepo "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/class/repository"
	_classService "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/class/usecase"

	_categoryController "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/category/delivery/http"
	_categoryRepo "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/category/repository"
	_categoryService "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/category/usecase"

	_cobsubcobController "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/cob-subcob/delivery/http"
	_cobsubcobRepo "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/cob-subcob/repository"
	_cobsubcobService "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/cob-subcob/usecase"

	_carriageController "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/carriage/delivery/http"
	_carriageRepo "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/carriage/repository"
	_carriageService "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/carriage/usecase"

	_conveyanceController "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/conveyance/delivery/http"
	_conveyanceRepo "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/conveyance/repository"
	_conveyanceService "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/conveyance/usecase"

	_contractorController "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/contractor/delivery/http"
	_contractorRepo "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/contractor/repository"
	_contractorService "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/contractor/usecase"
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
	router.Use(middleware.Logger())
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
	classRepo := _classRepo.NewClassRepository(dbBlips)
	authRepo := _authRepo.NewAuthRepository(dbBlips, emailServices)
	categoryRepo := _categoryRepo.NewCategoryRepository(dbBlips)
	cobsubcobRepo := _cobsubcobRepo.NewCobSubcobRepository(dbBlips)
	carriageRepo := _carriageRepo.NewCarriageRepository(dbBlips)
	conveyanceRepo := _conveyanceRepo.NewConveyanceRepository(dbBlips)
	contractorRepo := _contractorRepo.NewContractorRepository(dbBlips)

	shipyardRepo := _shipyardRepo.NewShipyardRepository(dbBlips)
	// Middlewares ------------------------------------------------------------------------------------------------------------------------------------------------------
	middlewareAuth := authmiddleware.NewMiddlewareAuth(authRepo)
	// middlewareAuth := middleware.NewMiddlewareAuth(
	// 	userRepo,
	// )
	middlewarePageRequest := _reqContext.NewMiddlewarePageRequest()
	roleService := _roleService.NewRoleUsecase(
		roleRepo,
		// dbValidations,
		timeoutContext,
	)
	_roleController.NewRoleHandler(
		router,
		roleService,
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
	//Classes
	classService := _classService.NewClassUsecase(
		classRepo,
		timeoutContext,
	)

	_classController.NewClassHandler(
		router,
		classService,
		middlewarePageRequest,
	)

	//Carriage
	carriageService := _carriageService.NewCarriageUsecase(
		carriageRepo,
		timeoutContext,
	)

	_carriageController.NewCarriageHandler(
		router,
		carriageService,
		middlewarePageRequest,
	)

	//Conveyance
	conveyanceService := _conveyanceService.NewConveyanceUsecase(
		conveyanceRepo,
		timeoutContext,
	)

	_conveyanceController.NewConveyanceHandler(
		router,
		conveyanceService,
		middlewarePageRequest,
	)

	//Contractor
	contractorService := _contractorService.NewContractorUsecase(
		contractorRepo,
		timeoutContext,
	)

	_contractorController.NewContractorHandler(
		router,
		contractorService,
		middlewarePageRequest,
	)

	// Shipyard
	shipyardService := _shipyardService.NewShipyardUsecase(shipyardRepo, timeoutContext)
	_shipyardController.NewShipyardHandler(router, shipyardService, middlewareAuth, middlewarePageRequest)

	//Categories
	categoryService := _categoryService.NewCategoryUsecase(
		categoryRepo,
		timeoutContext,
	)

	_categoryController.NewCategoryHandler(
		router,
		categoryService,
		middlewarePageRequest,
	)

	// cob-subcob
	cobsubcobService := _cobsubcobService.NewCobSubcobUsecase(
		categoryRepo,
		cobsubcobRepo,
		timeoutContext,
	)

	_cobsubcobController.NewCobSubcobHandler(
		router,
		cobsubcobService,
		middlewarePageRequest,
	)
	_shipyardController.NewShipyardHandler(router, shipyardService, middlewareAuth, middlewarePageRequest)

	return router

}
