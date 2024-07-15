package http

import (
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/models"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// contextShipyard is a middleware function that retrieves a shipyard by its ID or code.
func contextShipyard(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		// Get the ID or code from the request parameters.
		id := c.Param("id")

		// Declare a variable to hold the shipyard data.
		var data models.Shipyard

		// If the handler is not nil, find the shipyard by its ID or code.
		if handler != nil {
			data, err = handler.ShipyardUsecase.FindShipyardByUUIDOrCode(id)
			if err != nil {
				// If an error occurred, log it and return a bad request response.
				zap.S().Error(err)
				return c.JSON(returnBadRequest(err.Error()))
			}
		}

		// Set the shipyard data in the context.
		c.Set("shipyard", data)

		// Call the next handler in the chain.
		return next(c)
	}
}
