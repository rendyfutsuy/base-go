package http

import (
	"encoding/json"
	"io"
	"net/http"

	_reqContext "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/helper/middleware/request"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/helper/request"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/helper/response"
	cobsubcob "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/cob-subcob"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/cob-subcob/dto"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type ResponseError struct {
	Message string `json:"message"`
}

type CobSubcobHandler struct {
	CobSubcobUsecase cobsubcob.Usecase
	validator        *validator.Validate
	mwPageRequest    _reqContext.IMiddlewarePageRequest
}

func NewCobSubcobHandler(e *echo.Echo, us cobsubcob.Usecase, mwP _reqContext.IMiddlewarePageRequest) {
	handler := &CobSubcobHandler{
		CobSubcobUsecase: us,
		validator:        validator.New(),
		mwPageRequest:    mwP,
	}

	r := e.Group("v1/cobsubcob")

	r.POST("", handler.CreateCob)
	r.GET("/index-cob", handler.FetchIndexCob, handler.mwPageRequest.PageRequestCtx)
	r.GET("/index-subcob", handler.FetchIndexSubcob, handler.mwPageRequest.PageRequestCtx)
	r.GET("/all-cob", handler.FetchAllCob)
	r.GET("/all-subcob", handler.FetchAllSubcob)
	r.GET("/cob/:id", handler.FetchCobByID)
	r.GET("/subcob/:id", handler.FetchSubcobByID)
}

func (handler *CobSubcobHandler) CreateCob(c echo.Context) error {

	authID := "authID" //this should be from middleware

	// Read the category file
	catFile, err := c.FormFile("category_file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "Category file is required"})
	}

	catFileSrc, err := catFile.Open()
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "Failed to open category file"})
	}
	defer catFileSrc.Close()

	catFileBytes, err := io.ReadAll(catFileSrc)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "Failed to read category file"})
	}

	var categories []dto.CategoryJson
	err = json.Unmarshal(catFileBytes, &categories)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "Failed to unmarshal category file"})
	}

	// Read the COB file
	cobFile, err := c.FormFile("cob_file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "COB file is required"})
	}

	cobFileSrc, err := cobFile.Open()
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "Failed to open COB file"})
	}

	cobFileBytes, err := io.ReadAll(cobFileSrc)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "Failed to read COB file"})
	}

	var cobs []dto.CobJson
	err = json.Unmarshal(cobFileBytes, &cobs)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "Failed to unmarshal COB file"})
	}

	// Insert the categories and COBs
	err = handler.CobSubcobUsecase.InsertCategoryCobSubcob(categories, cobs, authID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{Message: "Failed to insert categories and COBs"})
	}

	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse("Success insert categories and COBs")

	return c.JSON(http.StatusOK, resp)
}

func (handler *CobSubcobHandler) FetchIndexCob(c echo.Context) error {
	pageRequest := c.Get("page_request").(*request.PageRequest)

	res, total, err := handler.CobSubcobUsecase.GetIndexCob(pageRequest)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	respCob := []dto.RespCob{}

	for _, v := range res {
		respCob = append(respCob, dto.ToRespCob(v))
	}

	respPag := response.PaginationResponse{}
	respPag, err = respPag.SetResponse(respCob, total, pageRequest.PerPage, pageRequest.Page)

	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, respPag)
}

func (handler *CobSubcobHandler) FetchIndexSubcob(c echo.Context) error {
	pageRequest := c.Get("page_request").(*request.PageRequest)

	res, total, err := handler.CobSubcobUsecase.GetIndexSubcob(pageRequest)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	respSubcob := []dto.RespSubcob{}

	for _, v := range res {
		respSubcob = append(respSubcob, dto.ToRespSubcob(v))
	}

	respPag := response.PaginationResponse{}
	respPag, err = respPag.SetResponse(respSubcob, total, pageRequest.PerPage, pageRequest.Page)

	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, respPag)
}

func (handler *CobSubcobHandler) FetchCobByID(c echo.Context) error {
	id := c.Param("id")

	res, err := handler.CobSubcobUsecase.GetCobByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	resp := dto.ToRespCob(*res)

	return c.JSON(http.StatusOK, resp)
}

func (handler *CobSubcobHandler) FetchSubcobByID(c echo.Context) error {
	id := c.Param("id")

	res, err := handler.CobSubcobUsecase.GetSubcobByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	resp := dto.ToRespSubcob(*res)

	return c.JSON(http.StatusOK, resp)
}

func (handler *CobSubcobHandler) FetchAllCob(c echo.Context) error {
	res, err := handler.CobSubcobUsecase.GetAllCob()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	resp := []dto.RespCob{}

	for _, v := range res {
		resp = append(resp, dto.ToRespCob(v))
	}

	return c.JSON(http.StatusOK, resp)
}

func (handler *CobSubcobHandler) FetchAllSubcob(c echo.Context) error {
	res, err := handler.CobSubcobUsecase.GetAllSubcob()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	resp := []dto.RespSubcob{}

	for _, v := range res {
		resp = append(resp, dto.ToRespSubcob(v))
	}

	return c.JSON(http.StatusOK, resp)
}
