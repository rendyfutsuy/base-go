package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rendyfutsuy/base-go/helpers/response"
	"github.com/rendyfutsuy/base-go/utils"
	"golang.org/x/time/rate"
)

type IThrottleMiddleware interface {
	Throttle() echo.MiddlewareFunc
}

type ThrottleMiddleware struct {
	limit int
}

func NewThrottleMiddleware() IThrottleMiddleware {
	// Get throttle limit from config, default to 1000 if not set
	limit := utils.ConfigVars.Int("throttle.limit")
	if limit <= 0 {
		limit = 1000 // Default limit
	}

	return &ThrottleMiddleware{
		limit: limit,
	}
}

// Throttle creates a rate limiting middleware with configurable limit
func (t *ThrottleMiddleware) Throttle() echo.MiddlewareFunc {
	// Create rate limiter with the configured limit
	config := middleware.RateLimiterConfig{
		Skipper: middleware.DefaultSkipper,
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{
				Rate:      rate.Limit(t.limit),
				Burst:     t.limit * 2, // Allow burst up to 2x the rate
				ExpiresIn: 3 * time.Minute,
			},
		),
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			// Extract identifier from Authorization header or IP address
			identifier := ctx.Request().Header.Get("Authorization")
			if identifier == "" {
				// Fallback to IP address if no auth header
				identifier = ctx.RealIP()
			} else {
				// Use a hash of the token to avoid storing full tokens
				identifier = "auth:" + identifier[:20] // Use first 20 chars of token
			}
			return identifier, nil
		},
		ErrorHandler: func(context echo.Context, err error) error {
			return context.JSON(http.StatusTooManyRequests, response.SetErrorResponse(
				http.StatusTooManyRequests,
				"Too many requests, please try again later",
			))
		},
		DenyHandler: func(context echo.Context, identifier string, err error) error {
			return context.JSON(http.StatusTooManyRequests, response.SetErrorResponse(
				http.StatusTooManyRequests,
				"Rate limit exceeded. Maximum "+strconv.Itoa(t.limit)+" requests per second allowed",
			))
		},
	}

	return middleware.RateLimiterWithConfig(config)
}
