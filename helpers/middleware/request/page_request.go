package request

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/rendyfutsuy/base-go/helpers/request"
)

type ResponseError struct {
	Message string `json:"message"`
}

type IMiddlewarePageRequest interface {
	PageRequestCtx(next echo.HandlerFunc) echo.HandlerFunc
	PageRequestCtxWithoutLimitation(next echo.HandlerFunc) echo.HandlerFunc
}

type MiddlewarePageRequest struct{}

func NewMiddlewarePageRequest() IMiddlewarePageRequest {
	return &MiddlewarePageRequest{}
}

// Middleware to parse page request context
func (m *MiddlewarePageRequest) PageRequestCtx(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Parse pagination parameters
		page, _ := strconv.ParseInt(c.QueryParam("page"), 10, 64)
		perPage, _ := strconv.ParseInt(c.QueryParam("per_page"), 10, 64)

		if perPage == 0 {
			return c.JSON(400, ResponseError{Message: "per_page must be greater than 0"})
		}

		// Parse sorting parameters
		sortBy := c.QueryParam("sort_by")
		sortOrder := c.QueryParam("sort_order")

		// parse search
		search := c.QueryParam("search")

		// Parse filtering parameters
		var filters []request.Filter
		jsonFilters := c.QueryParam("filter[]")

		if jsonFilters != "" {
			err := json.Unmarshal([]byte(jsonFilters), &filters)

			if err != nil {
				return c.JSON(http.StatusBadRequest, ResponseError{Message: "Invalid filter parameter"})

			}
		}

		projections := strings.Split(c.QueryParam("projections"), "|")
		tempArray := []string{}
		for _, v := range projections {
			if v != "" {
				tempArray = append(tempArray, v)
			}
		}

		// Create and attach PageRequest to context

		p := request.NewPageRequest(int(page), int(perPage), search, sortBy, sortOrder, filters, projections)
		c.Set("page_request", p)

		return next(c)
	}
}

func (m *MiddlewarePageRequest) PageRequestCtxWithoutLimitation(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Parse pagination parameters
		page, _ := strconv.ParseInt(c.QueryParam("page"), 10, 64)
		perPage, _ := strconv.ParseInt(c.QueryParam("per_page"), 10, 64)

		// Parse sorting parameters
		sortBy := c.QueryParam("sort_by")
		sortOrder := c.QueryParam("sort_order")

		// parse search
		search := c.QueryParam("search")

		// Parse filtering parameters
		var filters []request.Filter
		jsonFilters := c.QueryParam("filter[]")

		if jsonFilters != "" {
			err := json.Unmarshal([]byte(jsonFilters), &filters)

			if err != nil {
				return c.JSON(http.StatusBadRequest, ResponseError{Message: "Invalid filter parameter"})

			}
		}

		projections := strings.Split(c.QueryParam("projections"), "|")
		tempArray := []string{}
		for _, v := range projections {
			if v != "" {
				tempArray = append(tempArray, v)
			}
		}

		// Create and attach PageRequest to context

		p := request.NewPageRequest(int(page), int(perPage), search, sortBy, sortOrder, filters, projections)
		c.Set("page_request", p)

		return next(c)
	}
}
