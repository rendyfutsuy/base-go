package usecase

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rendyfutsuy/base-go/constants"
	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/models"
	mod "github.com/rendyfutsuy/base-go/modules/parameter"
	"github.com/rendyfutsuy/base-go/modules/parameter/dto"
	"github.com/rendyfutsuy/base-go/utils"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type parameterUsecase struct {
	repo mod.Repository
}

func NewParameterUsecase(repo mod.Repository) mod.Usecase {
	return &parameterUsecase{repo: repo}
}

func (u *parameterUsecase) Create(c echo.Context, reqBody *dto.ReqCreateParameter, authId string) (*models.Parameter, error) {
	ctx := c.Request().Context()

	// Check if code already exists
	exists, err := u.repo.ExistsByCode(ctx, reqBody.Code, uuid.Nil)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New(constants.ParameterCodeAlreadyExists)
	}

	// Check if name already exists
	exists, err = u.repo.ExistsByName(ctx, reqBody.Name, uuid.Nil)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New(constants.ParameterNameAlreadyExists)
	}

	return u.repo.Create(ctx, reqBody.Code, reqBody.Name, reqBody.Value, reqBody.Type, reqBody.Desc)
}

func (u *parameterUsecase) Update(c echo.Context, id string, reqBody *dto.ReqUpdateParameter, authId string) (*models.Parameter, error) {
	ctx := c.Request().Context()
	pid, err := utils.StringToUUID(id)
	if err != nil {
		return nil, err
	}

	// Check if code already exists (excluding current id)
	exists, err := u.repo.ExistsByCode(ctx, reqBody.Code, pid)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New(constants.ParameterCodeAlreadyExists)
	}

	// Check if name already exists (excluding current id)
	exists, err = u.repo.ExistsByName(ctx, reqBody.Name, pid)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New(constants.ParameterNameAlreadyExists)
	}

	res, err := u.repo.Update(ctx, pid, reqBody.Code, reqBody.Name, reqBody.Value, reqBody.Type, reqBody.Desc)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf(constants.ParameterNotFound, id)
		}
		return nil, err
	}
	return res, nil
}

func (u *parameterUsecase) Delete(c echo.Context, id string, authId string) error {
	ctx := c.Request().Context()
	pid, err := utils.StringToUUID(id)
	if err != nil {
		return err
	}
	return u.repo.Delete(ctx, pid)
}

func (u *parameterUsecase) GetByID(c echo.Context, id string) (*models.Parameter, error) {
	ctx := c.Request().Context()
	pid, err := utils.StringToUUID(id)
	if err != nil {
		return nil, err
	}
	return u.repo.GetByID(ctx, pid)
}

func (u *parameterUsecase) GetIndex(c echo.Context, req request.PageRequest, filter dto.ReqParameterIndexFilter) ([]models.Parameter, int, error) {
	ctx := c.Request().Context()
	// Search is already set in req.Search from PageRequest middleware
	return u.repo.GetIndex(ctx, req, filter)
}

func (u *parameterUsecase) GetAll(c echo.Context, filter dto.ReqParameterIndexFilter) ([]models.Parameter, error) {
	ctx := c.Request().Context()
	return u.repo.GetAll(ctx, filter)
}

func (u *parameterUsecase) Export(c echo.Context, filter dto.ReqParameterIndexFilter) ([]byte, error) {
	ctx := c.Request().Context()
	// Use GetAll for export without pagination
	list, err := u.repo.GetAll(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Create Excel file
	f := excelize.NewFile()
	sheet := "Parameters"
	f.SetSheetName("Sheet1", sheet)

	// Header
	f.SetCellValue(sheet, "A1", "Code")
	f.SetCellValue(sheet, "B1", "Name")
	f.SetCellValue(sheet, "C1", "Value")
	f.SetCellValue(sheet, "D1", "Type")
	f.SetCellValue(sheet, "E1", "Description")

	// Rows
	for i, p := range list {
		row := i + 2
		f.SetCellValue(sheet, "A"+strconv.Itoa(row), p.Code)
		f.SetCellValue(sheet, "B"+strconv.Itoa(row), p.Name)
		if p.Value != nil {
			f.SetCellValue(sheet, "C"+strconv.Itoa(row), *p.Value)
		}
		if p.Type != nil {
			f.SetCellValue(sheet, "D"+strconv.Itoa(row), *p.Type)
		}
		if p.Description != nil {
			f.SetCellValue(sheet, "E"+strconv.Itoa(row), *p.Description)
		}
	}

	// Write to buffer and return bytes
	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
