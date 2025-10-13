package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/utils"

	"github.com/rendyfutsuy/base-go/modules/auth"
)

type GeneralResponse struct {
	Message string `json:"message"`
}

type IMiddlewareAuth interface {
	AuthorizationCheck(next echo.HandlerFunc) echo.HandlerFunc
}

type MiddlewareAuth struct {
	authRepository auth.Repository
}

func NewMiddlewareAuth(authRepository auth.Repository) IMiddlewareAuth {
	return &MiddlewareAuth{
		authRepository: authRepository,
	}
}

func (a *MiddlewareAuth) AuthorizationCheck(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// get authorization header
		authorization := c.Request().Header.Get("Authorization")

		// if token set in header
		if authorization == "" {
			return c.JSON(http.StatusUnauthorized, GeneralResponse{Message: "unauthorized"})
		}

		// if bearer token not set
		tokenString := strings.Split(authorization, "Bearer ")[1]

		// tokenString = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDQzMTUyMDAsInVzZXJfaWQiOiIwMTkwZDlkNi1kNDI3LTc5ZDctYjgyYy1jODAzN2EzYWQ0N2YifQ.e3lE58KU_NLE_hn4FJNLMfmgDkhLQL8xKRLJIxdUjGY"

		if tokenString == "" {
			return c.JSON(http.StatusUnauthorized, GeneralResponse{Message: "unauthorized"})
		}

		// get user data from token
		userData, err := a.getUserData(tokenString)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, GeneralResponse{Message: err.Error()})
		}

		// Parse and validate the token
		claims := &jwt.RegisteredClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// Provide the key for validating the token
			// This should be the same key used to sign the token
			return []byte(utils.ConfigVars.String("jwt_key")), nil
		})

		// Check if the token is expired
		if claims.ExpiresAt != nil && claims.ExpiresAt.Unix() < time.Now().Unix() {
			return c.JSON(http.StatusUnauthorized, GeneralResponse{Message: "Token Expired"})
		}

		// check if token valid
		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, GeneralResponse{Message: "unauthorized"})
		}

		// set token to context
		c.Set("token", tokenString)
		c.Set("user", userData)

		return next(c)
	}
}

func (a *MiddlewareAuth) getUserData(token string) (user models.User, err error) {

	tokenString := token

	user, err = a.authRepository.GetUserByAccessToken(tokenString)
	if err != nil {
		return user, err
	}

	return user, err
}
