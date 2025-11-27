package http

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rendyfutsuy/base-go/constants"
	"github.com/rendyfutsuy/base-go/helpers/middleware"
	_reqContext "github.com/rendyfutsuy/base-go/helpers/middleware/request"
	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/helpers/response"
	"github.com/rendyfutsuy/base-go/modules/regency"
	"github.com/rendyfutsuy/base-go/modules/regency/dto"
)

type ResponseError struct {
	Message string `json:"message"`
}

type Response struct {
	Message string `json:"message"`
}

type RegencyHandler struct {
	Usecase              regency.Usecase
	validator            *validator.Validate
	mwPageRequest        _reqContext.IMiddlewarePageRequest
	middlewareAuth       middleware.IMiddlewareAuth
	middlewarePermission middleware.IMiddlewarePermission
}

func NewRegencyHandler(e *echo.Echo, uc regency.Usecase, mwP _reqContext.IMiddlewarePageRequest, auth middleware.IMiddlewareAuth, middlewarePermission middleware.IMiddlewarePermission) {
	h := &RegencyHandler{Usecase: uc, validator: validator.New(), mwPageRequest: mwP, middlewareAuth: auth, middlewarePermission: middlewarePermission}

	// Permissions (shared for all regency entities)
	// View:   province.view
	// Create: province.create
	// Update: province.update
	// Delete: province.delete
	// Export: province.export
	permissionToView := []string{"province.view"}
	permissionToCreate := []string{"province.create"}
	permissionToUpdate := []string{"province.update"}
	permissionToDelete := []string{"province.delete"}
	permissionToExport := []string{"province.export"}

	// Province routes
	provinceGroup := e.Group("/v1/province")
	provinceGroup.Use(h.middlewareAuth.AuthorizationCheck)

	provinceGroup.GET("", h.GetProvinceIndex, middleware.RequireActivatedUser, h.mwPageRequest.PageRequestCtx, h.middlewarePermission.PermissionValidation(permissionToView))
	provinceGroup.GET("/export", h.ExportProvince, middleware.RequireActivatedUser, h.middlewarePermission.PermissionValidation(permissionToExport))
	provinceGroup.GET("/:id", h.GetProvinceByID, middleware.RequireActivatedUser, h.middlewarePermission.PermissionValidation(permissionToView))
	provinceGroup.POST("", h.CreateProvince, middleware.RequireActivatedUser, h.middlewarePermission.PermissionValidation(permissionToCreate))
	provinceGroup.PUT("/:id", h.UpdateProvince, middleware.RequireActivatedUser, h.middlewarePermission.PermissionValidation(permissionToUpdate))
	provinceGroup.DELETE("/:id", h.DeleteProvince, middleware.RequireActivatedUser, h.middlewarePermission.PermissionValidation(permissionToDelete))

	// City routes
	cityGroup := e.Group("/v1/city")
	cityGroup.Use(h.middlewareAuth.AuthorizationCheck)

	cityGroup.GET("", h.GetCityIndex, middleware.RequireActivatedUser, h.mwPageRequest.PageRequestCtx, h.middlewarePermission.PermissionValidation(permissionToView))
	cityGroup.GET("/area-codes", h.GetCityAreaCodes, middleware.RequireActivatedUser, h.middlewarePermission.PermissionValidation(permissionToView))
	cityGroup.GET("/export", h.ExportCity, middleware.RequireActivatedUser, h.middlewarePermission.PermissionValidation(permissionToExport))
	cityGroup.GET("/:id", h.GetCityByID, middleware.RequireActivatedUser, h.middlewarePermission.PermissionValidation(permissionToView))
	cityGroup.POST("", h.CreateCity, middleware.RequireActivatedUser, h.middlewarePermission.PermissionValidation(permissionToCreate))
	cityGroup.PUT("/:id", h.UpdateCity, middleware.RequireActivatedUser, h.middlewarePermission.PermissionValidation(permissionToUpdate))
	cityGroup.DELETE("/:id", h.DeleteCity, middleware.RequireActivatedUser, h.middlewarePermission.PermissionValidation(permissionToDelete))

	// District routes
	districtGroup := e.Group("/v1/district")
	districtGroup.Use(h.middlewareAuth.AuthorizationCheck)

	districtGroup.GET("", h.GetDistrictIndex, middleware.RequireActivatedUser, h.mwPageRequest.PageRequestCtx, h.middlewarePermission.PermissionValidation(permissionToView))
	districtGroup.GET("/export", h.ExportDistrict, middleware.RequireActivatedUser, h.middlewarePermission.PermissionValidation(permissionToExport))
	districtGroup.GET("/:id", h.GetDistrictByID, middleware.RequireActivatedUser, h.middlewarePermission.PermissionValidation(permissionToView))
	districtGroup.POST("", h.CreateDistrict, middleware.RequireActivatedUser, h.middlewarePermission.PermissionValidation(permissionToCreate))
	districtGroup.PUT("/:id", h.UpdateDistrict, middleware.RequireActivatedUser, h.middlewarePermission.PermissionValidation(permissionToUpdate))
	districtGroup.DELETE("/:id", h.DeleteDistrict, middleware.RequireActivatedUser, h.middlewarePermission.PermissionValidation(permissionToDelete))

	// Subdistrict routes
	subdistrictGroup := e.Group("/v1/subdistrict")
	subdistrictGroup.Use(h.middlewareAuth.AuthorizationCheck)

	subdistrictGroup.GET("", h.GetSubdistrictIndex, middleware.RequireActivatedUser, h.mwPageRequest.PageRequestCtx, h.middlewarePermission.PermissionValidation(permissionToView))
	subdistrictGroup.GET("/export", h.ExportSubdistrict, middleware.RequireActivatedUser, h.middlewarePermission.PermissionValidation(permissionToExport))
	subdistrictGroup.GET("/:id", h.GetSubdistrictByID, middleware.RequireActivatedUser, h.middlewarePermission.PermissionValidation(permissionToView))
	subdistrictGroup.POST("", h.CreateSubdistrict, middleware.RequireActivatedUser, h.middlewarePermission.PermissionValidation(permissionToCreate))
	subdistrictGroup.PUT("/:id", h.UpdateSubdistrict, middleware.RequireActivatedUser, h.middlewarePermission.PermissionValidation(permissionToUpdate))
	subdistrictGroup.DELETE("/:id", h.DeleteSubdistrict, middleware.RequireActivatedUser, h.middlewarePermission.PermissionValidation(permissionToDelete))
}

// Province Handlers

// CreateProvince godoc
// @Summary		Create a new province
// @Description	Create a new province with provided name
// @Tags			Regency - Province
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			request	body		dto.ReqCreateProvince	true	"Province creation data"
// @Success		200		{object}	response.NonPaginationResponse{data=dto.RespProvince}	"Successfully created province"
// @Failure		400		{object}	response.NonPaginationResponse	"Bad request - validation error"
// @Failure		401		{object}	response.NonPaginationResponse	"Unauthorized"
// @Router			/v1/province [post]
func (h *RegencyHandler) CreateProvince(c echo.Context) error {
	req := new(dto.ReqCreateProvince)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	res, err := h.Usecase.CreateProvince(c, req, "")
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(dto.ToRespProvince(*res))
	return c.JSON(http.StatusOK, resp)
}

// UpdateProvince godoc
// @Summary		Update province
// @Description	Update an existing province's information
// @Tags			Regency - Province
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			id		path	string					true	"Province UUID"
// @Param			request	body	dto.ReqUpdateProvince	true	"Updated province data"
// @Success		200		{object}	response.NonPaginationResponse{data=dto.RespProvince}	"Successfully updated province"
// @Failure		400		{object}	response.NonPaginationResponse	"Bad request - validation error"
// @Failure		401		{object}	response.NonPaginationResponse	"Unauthorized"
// @Failure		404		{object}	response.NonPaginationResponse	"Province not found"
// @Router			/v1/province/{id} [put]
func (h *RegencyHandler) UpdateProvince(c echo.Context) error {
	id := c.Param("id")
	req := new(dto.ReqUpdateProvince)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	res, err := h.Usecase.UpdateProvince(c, id, req, "")
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(dto.ToRespProvince(*res))
	return c.JSON(http.StatusOK, resp)
}

// DeleteProvince godoc
// @Summary		Delete province
// @Description	Delete an existing province by ID
// @Tags			Regency - Province
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			id		path	string	true	"Province UUID"
// @Success		200		{object}	response.NonPaginationResponse	"Successfully deleted province"
// @Failure		400		{object}	response.NonPaginationResponse	"Bad request - invalid UUID"
// @Failure		401		{object}	response.NonPaginationResponse	"Unauthorized"
// @Failure		404		{object}	response.NonPaginationResponse	"Province not found"
// @Router			/v1/province/{id} [delete]
func (h *RegencyHandler) DeleteProvince(c echo.Context) error {
	id := c.Param("id")
	if err := h.Usecase.DeleteProvince(c, id, ""); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(Response{Message: constants.ProvinceDeleteSuccess})
	return c.JSON(http.StatusOK, resp)
}

// GetProvinceIndex godoc
// @Summary		Get list of provinces with pagination
// @Description	Retrieve a paginated list of provinces with optional filters
// @Tags			Regency - Province
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			page		query		int						false	"Page number"
// @Param			per_page	query		int						false	"Items per page"
// @Param			search		query		string					false	"Search keyword"
// @Param			filter		query		dto.ReqProvinceIndexFilter	false	"Filter options"
// @Success		200			{object}	response.PaginationResponse{data=[]dto.RespProvinceIndex}	"Successfully retrieved provinces"
// @Failure		400			{object}	response.NonPaginationResponse	"Bad request"
// @Failure		401			{object}	response.NonPaginationResponse	"Unauthorized"
// @Router			/v1/province [get]
func (h *RegencyHandler) GetProvinceIndex(c echo.Context) error {
	pageRequest := c.Get("page_request").(*request.PageRequest)
	filter := new(dto.ReqProvinceIndexFilter)
	if err := c.Bind(filter); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	if err := c.Validate(filter); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	res, total, err := h.Usecase.GetProvinceIndex(c, *pageRequest, *filter)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	respProvince := []dto.RespProvinceIndex{}
	for _, v := range res {
		respProvince = append(respProvince, dto.ToRespProvinceIndex(v))
	}
	respPag := response.PaginationResponse{}
	respPag, err = respPag.SetResponse(respProvince, total, pageRequest.PerPage, pageRequest.Page)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	return c.JSON(http.StatusOK, respPag)
}

// GetProvinceByID godoc
// @Summary		Get province by ID
// @Description	Retrieve a single province by its ID
// @Tags			Regency - Province
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			id		path		string	true	"Province UUID"
// @Success		200		{object}	response.NonPaginationResponse{data=dto.RespProvince}	"Successfully retrieved province"
// @Failure		400		{object}	response.NonPaginationResponse	"Bad request - invalid UUID"
// @Failure		401		{object}	response.NonPaginationResponse	"Unauthorized"
// @Failure		404		{object}	response.NonPaginationResponse	"Province not found"
// @Router			/v1/province/{id} [get]
func (h *RegencyHandler) GetProvinceByID(c echo.Context) error {
	id := c.Param("id")
	res, err := h.Usecase.GetProvinceByID(c, id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(dto.ToRespProvince(*res))
	return c.JSON(http.StatusOK, resp)
}

// ExportProvince godoc
// @Summary		Export provinces to Excel
// @Description	Export provinces to Excel file (.xlsx) with optional search and filter
// @Tags			Regency - Province
// @Accept			json
// @Produce		application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Security		BearerAuth
// @Param			search		query		string					false	"Search keyword"
// @Param			filter		query		dto.ReqProvinceIndexFilter	false	"Filter options"
// @Success		200			{file}		binary	"Excel file with provinces data"
// @Failure		400			{object}	response.NonPaginationResponse	"Bad request"
// @Failure		401			{object}	response.NonPaginationResponse	"Unauthorized"
// @Router			/v1/province/export [get]
func (h *RegencyHandler) ExportProvince(c echo.Context) error {
	filter := new(dto.ReqProvinceIndexFilter)
	if err := c.Bind(filter); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	if err := c.Validate(filter); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	excelBytes, err := h.Usecase.ExportProvince(c, *filter)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	c.Response().Header().Set(echo.HeaderContentType, constants.ExcelContent)
	c.Response().Header().Set(constants.FieldContentDisposition, constants.ExcelContentDisposition("provinces.xlsx"))
	return c.Blob(http.StatusOK, constants.ExcelContent, excelBytes)
}

// City Handlers

// CreateCity godoc
// @Summary		Create a new city
// @Description	Create a new city with provided province_id, name, dan optional area_code
// @Tags			Regency - City
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			request	body		dto.ReqCreateCity	true	"City creation data"
// @Success		200		{object}	response.NonPaginationResponse{data=dto.RespCity}	"Successfully created city"
// @Failure		400		{object}	response.NonPaginationResponse	"Bad request - validation error"
// @Failure		401		{object}	response.NonPaginationResponse	"Unauthorized"
// @Router			/v1/city [post]
func (h *RegencyHandler) CreateCity(c echo.Context) error {
	req := new(dto.ReqCreateCity)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	res, err := h.Usecase.CreateCity(c, req, "")
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(dto.ToRespCity(*res))
	return c.JSON(http.StatusOK, resp)
}

// UpdateCity godoc
// @Summary		Update city
// @Description	Update an existing city's information termasuk area_code
// @Tags			Regency - City
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			id		path	string				true	"City UUID"
// @Param			request	body	dto.ReqUpdateCity	true	"Updated city data"
// @Success		200		{object}	response.NonPaginationResponse{data=dto.RespCity}	"Successfully updated city"
// @Failure		400		{object}	response.NonPaginationResponse	"Bad request - validation error"
// @Failure		401		{object}	response.NonPaginationResponse	"Unauthorized"
// @Failure		404		{object}	response.NonPaginationResponse	"City not found"
// @Router			/v1/city/{id} [put]
func (h *RegencyHandler) UpdateCity(c echo.Context) error {
	id := c.Param("id")
	req := new(dto.ReqUpdateCity)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	res, err := h.Usecase.UpdateCity(c, id, req, "")
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(dto.ToRespCity(*res))
	return c.JSON(http.StatusOK, resp)
}

// DeleteCity godoc
// @Summary		Delete city
// @Description	Delete an existing city by ID
// @Tags			Regency - City
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			id		path	string	true	"City UUID"
// @Success		200		{object}	response.NonPaginationResponse	"Successfully deleted city"
// @Failure		400		{object}	response.NonPaginationResponse	"Bad request - invalid UUID"
// @Failure		401		{object}	response.NonPaginationResponse	"Unauthorized"
// @Failure		404		{object}	response.NonPaginationResponse	"City not found"
// @Router			/v1/city/{id} [delete]
func (h *RegencyHandler) DeleteCity(c echo.Context) error {
	id := c.Param("id")
	if err := h.Usecase.DeleteCity(c, id, ""); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(Response{Message: constants.CityDeleteSuccess})
	return c.JSON(http.StatusOK, resp)
}

// GetCityIndex godoc
// @Summary		Get list of cities with pagination
// @Description	Retrieve a paginated list of cities with optional filters
// @Tags			Regency - City
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			page		query		int					false	"Page number"
// @Param			per_page	query		int					false	"Items per page"
// @Param			search		query		string				false	"Search keyword"
// @Param			filter		query		dto.ReqCityIndexFilter	false	"Filter options"
// @Success		200			{object}	response.PaginationResponse{data=[]dto.RespCityIndex}	"Successfully retrieved cities"
// @Failure		400			{object}	response.NonPaginationResponse	"Bad request"
// @Failure		401			{object}	response.NonPaginationResponse	"Unauthorized"
// @Router			/v1/city [get]
func (h *RegencyHandler) GetCityIndex(c echo.Context) error {
	pageRequest := c.Get("page_request").(*request.PageRequest)
	filter := new(dto.ReqCityIndexFilter)
	if err := c.Bind(filter); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	if err := c.Validate(filter); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	res, total, err := h.Usecase.GetCityIndex(c, *pageRequest, *filter)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	respCity := []dto.RespCityIndex{}
	for _, v := range res {
		respCity = append(respCity, dto.ToRespCityIndex(v))
	}
	respPag := response.PaginationResponse{}
	respPag, err = respPag.SetResponse(respCity, total, pageRequest.PerPage, pageRequest.Page)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	return c.JSON(http.StatusOK, respPag)
}

// GetCityAreaCodes godoc
// @Summary		Get distinct city area codes
// @Description	Retrieve distinct area_code values from city table with optional search by area_code
// @Tags			Regency - City
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			search	query		string	false	"Search area_code"
// @Success		200		{object}	response.NonPaginationResponse{data=[]dto.RespCityAreaCode}	"Successfully retrieved area codes"
// @Failure		400		{object}	response.NonPaginationResponse	"Bad request"
// @Failure		401		{object}	response.NonPaginationResponse	"Unauthorized"
// @Router			/v1/city/area-codes [get]
func (h *RegencyHandler) GetCityAreaCodes(c echo.Context) error {
	search := c.QueryParam("search")
	res, err := h.Usecase.GetCityAreaCodes(c, search)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	respData := make([]dto.RespCityAreaCode, 0, len(res))
	for _, code := range res {
		respData = append(respData, dto.RespCityAreaCode{AreaCode: code})
	}

	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(respData)
	return c.JSON(http.StatusOK, resp)
}

// GetCityByID godoc
// @Summary		Get city by ID
// @Description	Retrieve a single city by its ID
// @Tags			Regency - City
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			id		path		string	true	"City UUID"
// @Success		200		{object}	response.NonPaginationResponse{data=dto.RespCity}	"Successfully retrieved city"
// @Failure		400		{object}	response.NonPaginationResponse	"Bad request - invalid UUID"
// @Failure		401		{object}	response.NonPaginationResponse	"Unauthorized"
// @Failure		404		{object}	response.NonPaginationResponse	"City not found"
// @Router			/v1/city/{id} [get]
func (h *RegencyHandler) GetCityByID(c echo.Context) error {
	id := c.Param("id")
	res, err := h.Usecase.GetCityByID(c, id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(dto.ToRespCity(*res))
	return c.JSON(http.StatusOK, resp)
}

// ExportCity godoc
// @Summary		Export cities to Excel
// @Description	Export cities to Excel file (.xlsx) with optional search and filter
// @Tags			Regency - City
// @Accept			json
// @Produce		application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Security		BearerAuth
// @Param			search		query		string					false	"Search keyword"
// @Param			filter		query		dto.ReqCityIndexFilter	false	"Filter options"
// @Success		200			{file}		binary	"Excel file with cities data"
// @Failure		400			{object}	response.NonPaginationResponse	"Bad request"
// @Failure		401			{object}	response.NonPaginationResponse	"Unauthorized"
// @Router			/v1/city/export [get]
func (h *RegencyHandler) ExportCity(c echo.Context) error {
	filter := new(dto.ReqCityIndexFilter)
	if err := c.Bind(filter); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	if err := c.Validate(filter); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	excelBytes, err := h.Usecase.ExportCity(c, *filter)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	c.Response().Header().Set(echo.HeaderContentType, constants.ExcelContent)
	c.Response().Header().Set(constants.FieldContentDisposition, constants.ExcelContentDisposition("cities.xlsx"))
	return c.Blob(http.StatusOK, constants.ExcelContent, excelBytes)
}

// District Handlers

// CreateDistrict godoc
// @Summary		Create a new district
// @Description	Create a new district with provided city_id and name
// @Tags			Regency - District
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			request	body		dto.ReqCreateDistrict	true	"District creation data"
// @Success		200		{object}	response.NonPaginationResponse{data=dto.RespDistrict}	"Successfully created district"
// @Failure		400		{object}	response.NonPaginationResponse	"Bad request - validation error"
// @Failure		401		{object}	response.NonPaginationResponse	"Unauthorized"
// @Router			/v1/district [post]
func (h *RegencyHandler) CreateDistrict(c echo.Context) error {
	req := new(dto.ReqCreateDistrict)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	res, err := h.Usecase.CreateDistrict(c, req, "")
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(dto.ToRespDistrict(*res))
	return c.JSON(http.StatusOK, resp)
}

// UpdateDistrict godoc
// @Summary		Update district
// @Description	Update an existing district's information
// @Tags			Regency - District
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			id		path	string					true	"District UUID"
// @Param			request	body	dto.ReqUpdateDistrict	true	"Updated district data"
// @Success		200		{object}	response.NonPaginationResponse{data=dto.RespDistrict}	"Successfully updated district"
// @Failure		400		{object}	response.NonPaginationResponse	"Bad request - validation error"
// @Failure		401		{object}	response.NonPaginationResponse	"Unauthorized"
// @Failure		404		{object}	response.NonPaginationResponse	"District not found"
// @Router			/v1/district/{id} [put]
func (h *RegencyHandler) UpdateDistrict(c echo.Context) error {
	id := c.Param("id")
	req := new(dto.ReqUpdateDistrict)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	res, err := h.Usecase.UpdateDistrict(c, id, req, "")
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(dto.ToRespDistrict(*res))
	return c.JSON(http.StatusOK, resp)
}

// DeleteDistrict godoc
// @Summary		Delete district
// @Description	Delete an existing district by ID
// @Tags			Regency - District
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			id		path	string	true	"District UUID"
// @Success		200		{object}	response.NonPaginationResponse	"Successfully deleted district"
// @Failure		400		{object}	response.NonPaginationResponse	"Bad request - invalid UUID"
// @Failure		401		{object}	response.NonPaginationResponse	"Unauthorized"
// @Failure		404		{object}	response.NonPaginationResponse	"District not found"
// @Router			/v1/district/{id} [delete]
func (h *RegencyHandler) DeleteDistrict(c echo.Context) error {
	id := c.Param("id")
	if err := h.Usecase.DeleteDistrict(c, id, ""); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(Response{Message: constants.DistrictDeleteSuccess})
	return c.JSON(http.StatusOK, resp)
}

// GetDistrictIndex godoc
// @Summary		Get list of districts with pagination
// @Description	Retrieve a paginated list of districts with optional filters
// @Tags			Regency - District
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			page		query		int						false	"Page number"
// @Param			per_page	query		int						false	"Items per page"
// @Param			search		query		string					false	"Search keyword"
// @Param			filter		query		dto.ReqDistrictIndexFilter	false	"Filter options"
// @Success		200			{object}	response.PaginationResponse{data=[]dto.RespDistrictIndex}	"Successfully retrieved districts"
// @Failure		400			{object}	response.NonPaginationResponse	"Bad request"
// @Failure		401			{object}	response.NonPaginationResponse	"Unauthorized"
// @Router			/v1/district [get]
func (h *RegencyHandler) GetDistrictIndex(c echo.Context) error {
	pageRequest := c.Get("page_request").(*request.PageRequest)
	filter := new(dto.ReqDistrictIndexFilter)
	if err := c.Bind(filter); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	if err := c.Validate(filter); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	res, total, err := h.Usecase.GetDistrictIndex(c, *pageRequest, *filter)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	respDistrict := []dto.RespDistrictIndex{}
	for _, v := range res {
		respDistrict = append(respDistrict, dto.ToRespDistrictIndex(v))
	}
	respPag := response.PaginationResponse{}
	respPag, err = respPag.SetResponse(respDistrict, total, pageRequest.PerPage, pageRequest.Page)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	return c.JSON(http.StatusOK, respPag)
}

// GetDistrictByID godoc
// @Summary		Get district by ID
// @Description	Retrieve a single district by its ID
// @Tags			Regency - District
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			id		path		string	true	"District UUID"
// @Success		200		{object}	response.NonPaginationResponse{data=dto.RespDistrict}	"Successfully retrieved district"
// @Failure		400		{object}	response.NonPaginationResponse	"Bad request - invalid UUID"
// @Failure		401		{object}	response.NonPaginationResponse	"Unauthorized"
// @Failure		404		{object}	response.NonPaginationResponse	"District not found"
// @Router			/v1/district/{id} [get]
func (h *RegencyHandler) GetDistrictByID(c echo.Context) error {
	id := c.Param("id")
	res, err := h.Usecase.GetDistrictByID(c, id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(dto.ToRespDistrict(*res))
	return c.JSON(http.StatusOK, resp)
}

// ExportDistrict godoc
// @Summary		Export districts to Excel
// @Description	Export districts to Excel file (.xlsx) with optional search and filter
// @Tags			Regency - District
// @Accept			json
// @Produce		application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Security		BearerAuth
// @Param			search		query		string					false	"Search keyword"
// @Param			filter		query		dto.ReqDistrictIndexFilter	false	"Filter options"
// @Success		200			{file}		binary	"Excel file with districts data"
// @Failure		400			{object}	response.NonPaginationResponse	"Bad request"
// @Failure		401			{object}	response.NonPaginationResponse	"Unauthorized"
// @Router			/v1/district/export [get]
func (h *RegencyHandler) ExportDistrict(c echo.Context) error {
	filter := new(dto.ReqDistrictIndexFilter)
	if err := c.Bind(filter); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	if err := c.Validate(filter); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	excelBytes, err := h.Usecase.ExportDistrict(c, *filter)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	c.Response().Header().Set(echo.HeaderContentType, constants.ExcelContent)
	c.Response().Header().Set(constants.FieldContentDisposition, constants.ExcelContentDisposition("districts.xlsx"))
	return c.Blob(http.StatusOK, constants.ExcelContent, excelBytes)
}

// Subdistrict Handlers

// CreateSubdistrict godoc
// @Summary		Create a new subdistrict
// @Description	Create a new subdistrict with provided district_id and name
// @Tags			Regency - Subdistrict
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			request	body		dto.ReqCreateSubdistrict	true	"Subdistrict creation data"
// @Success		200		{object}	response.NonPaginationResponse{data=dto.RespSubdistrict}	"Successfully created subdistrict"
// @Failure		400		{object}	response.NonPaginationResponse	"Bad request - validation error"
// @Failure		401		{object}	response.NonPaginationResponse	"Unauthorized"
// @Router			/v1/subdistrict [post]
func (h *RegencyHandler) CreateSubdistrict(c echo.Context) error {
	req := new(dto.ReqCreateSubdistrict)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	res, err := h.Usecase.CreateSubdistrict(c, req, "")
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(dto.ToRespSubdistrict(*res))
	return c.JSON(http.StatusOK, resp)
}

// UpdateSubdistrict godoc
// @Summary		Update subdistrict
// @Description	Update an existing subdistrict's information
// @Tags			Regency - Subdistrict
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			id		path	string					true	"Subdistrict UUID"
// @Param			request	body	dto.ReqUpdateSubdistrict	true	"Updated subdistrict data"
// @Success		200		{object}	response.NonPaginationResponse{data=dto.RespSubdistrict}	"Successfully updated subdistrict"
// @Failure		400		{object}	response.NonPaginationResponse	"Bad request - validation error"
// @Failure		401		{object}	response.NonPaginationResponse	"Unauthorized"
// @Failure		404		{object}	response.NonPaginationResponse	"Subdistrict not found"
// @Router			/v1/subdistrict/{id} [put]
func (h *RegencyHandler) UpdateSubdistrict(c echo.Context) error {
	id := c.Param("id")
	req := new(dto.ReqUpdateSubdistrict)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	res, err := h.Usecase.UpdateSubdistrict(c, id, req, "")
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(dto.ToRespSubdistrict(*res))
	return c.JSON(http.StatusOK, resp)
}

// DeleteSubdistrict godoc
// @Summary		Delete subdistrict
// @Description	Delete an existing subdistrict by ID
// @Tags			Regency - Subdistrict
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			id		path	string	true	"Subdistrict UUID"
// @Success		200		{object}	response.NonPaginationResponse	"Successfully deleted subdistrict"
// @Failure		400		{object}	response.NonPaginationResponse	"Bad request - invalid UUID"
// @Failure		401		{object}	response.NonPaginationResponse	"Unauthorized"
// @Failure		404		{object}	response.NonPaginationResponse	"Subdistrict not found"
// @Router			/v1/subdistrict/{id} [delete]
func (h *RegencyHandler) DeleteSubdistrict(c echo.Context) error {
	id := c.Param("id")
	if err := h.Usecase.DeleteSubdistrict(c, id, ""); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(Response{Message: constants.SubdistrictDeleteSuccess})
	return c.JSON(http.StatusOK, resp)
}

// GetSubdistrictIndex godoc
// @Summary		Get list of subdistricts with pagination
// @Description	Retrieve a paginated list of subdistricts with optional filters
// @Tags			Regency - Subdistrict
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			page		query		int							false	"Page number"
// @Param			per_page	query		int							false	"Items per page"
// @Param			search		query		string						false	"Search keyword"
// @Param			filter		query		dto.ReqSubdistrictIndexFilter	false	"Filter options"
// @Success		200			{object}	response.PaginationResponse{data=[]dto.RespSubdistrictIndex}	"Successfully retrieved subdistricts"
// @Failure		400			{object}	response.NonPaginationResponse	"Bad request"
// @Failure		401			{object}	response.NonPaginationResponse	"Unauthorized"
// @Router			/v1/subdistrict [get]
func (h *RegencyHandler) GetSubdistrictIndex(c echo.Context) error {
	pageRequest := c.Get("page_request").(*request.PageRequest)
	filter := new(dto.ReqSubdistrictIndexFilter)
	if err := c.Bind(filter); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	if err := c.Validate(filter); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	res, total, err := h.Usecase.GetSubdistrictIndex(c, *pageRequest, *filter)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	respSubdistrict := []dto.RespSubdistrictIndex{}
	for _, v := range res {
		respSubdistrict = append(respSubdistrict, dto.ToRespSubdistrictIndex(v))
	}
	respPag := response.PaginationResponse{}
	respPag, err = respPag.SetResponse(respSubdistrict, total, pageRequest.PerPage, pageRequest.Page)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	return c.JSON(http.StatusOK, respPag)
}

// GetSubdistrictByID godoc
// @Summary		Get subdistrict by ID
// @Description	Retrieve a single subdistrict by its ID
// @Tags			Regency - Subdistrict
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			id		path		string	true	"Subdistrict UUID"
// @Success		200		{object}	response.NonPaginationResponse{data=dto.RespSubdistrict}	"Successfully retrieved subdistrict"
// @Failure		400		{object}	response.NonPaginationResponse	"Bad request - invalid UUID"
// @Failure		401		{object}	response.NonPaginationResponse	"Unauthorized"
// @Failure		404		{object}	response.NonPaginationResponse	"Subdistrict not found"
// @Router			/v1/subdistrict/{id} [get]
func (h *RegencyHandler) GetSubdistrictByID(c echo.Context) error {
	id := c.Param("id")
	res, err := h.Usecase.GetSubdistrictByID(c, id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(dto.ToRespSubdistrict(*res))
	return c.JSON(http.StatusOK, resp)
}

// ExportSubdistrict godoc
// @Summary		Export subdistricts to Excel
// @Description	Export subdistricts to Excel file (.xlsx) with optional search and filter
// @Tags			Regency - Subdistrict
// @Accept			json
// @Produce		application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Security		BearerAuth
// @Param			search		query		string						false	"Search keyword"
// @Param			filter		query		dto.ReqSubdistrictIndexFilter	false	"Filter options"
// @Success		200			{file}		binary	"Excel file with subdistricts data"
// @Failure		400			{object}	response.NonPaginationResponse	"Bad request"
// @Failure		401			{object}	response.NonPaginationResponse	"Unauthorized"
// @Router			/v1/subdistrict/export [get]
func (h *RegencyHandler) ExportSubdistrict(c echo.Context) error {
	filter := new(dto.ReqSubdistrictIndexFilter)
	if err := c.Bind(filter); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	if err := c.Validate(filter); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	excelBytes, err := h.Usecase.ExportSubdistrict(c, *filter)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	c.Response().Header().Set(echo.HeaderContentType, constants.ExcelContent)
	c.Response().Header().Set(constants.FieldContentDisposition, constants.ExcelContentDisposition("subdistricts.xlsx"))
	return c.Blob(http.StatusOK, constants.ExcelContent, excelBytes)
}
