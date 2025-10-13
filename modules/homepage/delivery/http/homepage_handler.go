package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rendyfutsuy/base-go/constants"
)

type Response struct {
	Message string `json:"message"`
	Version string `json:"version"`
}

func DefaultHomepage(c echo.Context) error {
	response := Response{
		Message: "Define version uses by Backend Resource, last updated 2025/08/11 08.11 WIB",
		Version: constants.Version,
	}
	return c.JSON(http.StatusOK, response)
}
