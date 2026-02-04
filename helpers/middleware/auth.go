package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rendyfutsuy/base-go/helpers/response"
	"github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/utils"
	"github.com/rendyfutsuy/base-go/utils/token_storage"
)

type GeneralResponse struct {
	Message string `json:"message"`
}

type IMiddlewareAuth interface {
	AuthorizationCheck(next echo.HandlerFunc) echo.HandlerFunc
}

type MiddlewareAuth struct {
}

func NewMiddlewareAuth() IMiddlewareAuth {
	return &MiddlewareAuth{}
}

func (a *MiddlewareAuth) AuthorizationCheck(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		// get authorization header
		authorization := c.Request().Header.Get("Authorization")

		// if token set in header
		if authorization == "" {
			return c.JSON(http.StatusUnauthorized, response.SetErrorResponse(http.StatusUnauthorized, "unauthorized"))
		}

		// if bearer token not set
		tokenString := strings.Split(authorization, "Bearer ")[1]

		// tokenString = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDQzMTUyMDAsInVzZXJfaWQiOiIwMTkwZDlkNi1kNDI3LTc5ZDctYjgyYy1jODAzN2EzYWQ0N2YifQ.e3lE58KU_NLE_hn4FJNLMfmgDkhLQL8xKRLJIxdUjGY"

		if tokenString == "" {
			return c.JSON(http.StatusUnauthorized, response.SetErrorResponse(http.StatusUnauthorized, "unauthorized"))
		}

		// Check if this is a refresh token endpoint
		// For refresh token, we allow expired access tokens as long as Redis session still exists
		isRefreshTokenEndpoint := strings.HasSuffix(c.Path(), "/refresh-token") || c.Path() == "/v1/auth/refresh-token"

		// Parse and validate the token
		claims := &jwt.RegisteredClaims{}

		// For refresh token endpoint, we need to allow expired tokens
		// We'll parse with a custom parser that skips expiration validation
		var token *jwt.Token
		var err error

		if isRefreshTokenEndpoint {
			// Parse without validating expiration (we'll validate via Redis session TTL)
			// Create parser that skips claims validation (including expiration)
			parser := jwt.NewParser(jwt.WithoutClaimsValidation())
			token, err = parser.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(utils.ConfigVars.String("jwt_key")), nil
			})

			// For refresh token, we only fail if signature is invalid
			// Expired token is fine as long as Redis session exists (validated in getUserData)
			if err != nil || (token != nil && !token.Valid) {
				return c.JSON(http.StatusUnauthorized, response.SetErrorResponse(http.StatusUnauthorized, "unauthorized"))
			}
		} else {
			// Standard parsing with expiration validation
			token, err = jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(utils.ConfigVars.String("jwt_key")), nil
			})

			// Check if the token is expired (for non-refresh endpoints)
			if claims.ExpiresAt != nil && claims.ExpiresAt.Unix() < time.Now().Unix() {
				return c.JSON(http.StatusUnauthorized, response.SetErrorResponse(http.StatusUnauthorized, "Token Expired"))
			}

			// check if token valid (signature and structure)
			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, response.SetErrorResponse(http.StatusUnauthorized, "unauthorized"))
			}
		}

		// get user data from token (validates Redis session exists)
		// For refresh token: if session exists in Redis (not expired), allow refresh
		// For other endpoints: session must exist and token must not be expired
		userData, err := a.getUserData(ctx, tokenString)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, response.SetErrorResponse(http.StatusUnauthorized, err.Error()))
		}

		// set token and user data to context
		c.Set("token", tokenString)
		c.Set("user", userData)
		c.Set("userId", userData.ID.String())

		return next(c)
	}
}

func (a *MiddlewareAuth) getUserData(ctx context.Context, token string) (user models.User, err error) {

	tokenString := token

	user, err = token_storage.ValidateAccessToken(ctx, tokenString)
	if err != nil {
		return user, err
	}

	if user.ID == uuid.Nil {
		return user, errors.New("user session is unauthorized")
	}

	return user, err
}
