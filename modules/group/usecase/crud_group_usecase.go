package usecase

import (
	"errors"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rendyfutsuy/base-go/constants"
	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/models"
	mod "github.com/rendyfutsuy/base-go/modules/group"
	"github.com/rendyfutsuy/base-go/modules/group/dto"
	"github.com/rendyfutsuy/base-go/utils"
	"github.com/xuri/excelize/v2"
)

type groupUsecase struct {
	repo mod.Repository
}

func NewGroupUsecase(repo mod.Repository) mod.Usecase {
	return &groupUsecase{repo: repo}
}

func (u *groupUsecase) Create(c echo.Context, reqBody *dto.ReqCreateGroup, authId string) (*models.GoodsGroup, error) {
	ctx := c.Request().Context()
	exists, err := u.repo.ExistsByName(ctx, reqBody.Name, uuid.Nil)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New(constants.GroupNameAlreadyExists)
	}
	return u.repo.Create(ctx, reqBody.Name)
}

func (u *groupUsecase) Update(c echo.Context, id string, reqBody *dto.ReqUpdateGroup, authId string) (*models.GoodsGroup, error) {
	ctx := c.Request().Context()
	gid, err := utils.StringToUUID(id)
	if err != nil {
		return nil, err
	}
	exists, err := u.repo.ExistsByName(ctx, reqBody.Name, gid)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New(constants.GroupNameAlreadyExists)
	}
	return u.repo.Update(ctx, gid, reqBody.Name)
}

func (u *groupUsecase) Delete(c echo.Context, id string, authId string) error {
	ctx := c.Request().Context()
	gid, err := utils.StringToUUID(id)
	if err != nil {
		return err
	}
	return u.repo.Delete(ctx, gid)
}

func (u *groupUsecase) GetByID(c echo.Context, id string) (*models.GoodsGroup, error) {
	ctx := c.Request().Context()
	gid, err := utils.StringToUUID(id)
	if err != nil {
		return nil, err
	}
	return u.repo.GetByID(ctx, gid)
}

func (u *groupUsecase) GetIndex(c echo.Context, req request.PageRequest, filter dto.ReqGroupIndexFilter) ([]models.GoodsGroup, int, error) {
	ctx := c.Request().Context()
	// Search is already set in req.Search from PageRequest middleware
	return u.repo.GetIndex(ctx, req, filter)
}

func (u *groupUsecase) GetAll(c echo.Context, filter dto.ReqGroupIndexFilter) ([]models.GoodsGroup, error) {
	ctx := c.Request().Context()
	return u.repo.GetAll(ctx, filter)
}

func (u *groupUsecase) Export(c echo.Context, filter dto.ReqGroupIndexFilter) ([]byte, error) {
	ctx := c.Request().Context()
	// Extract search from query param and set to filter (filter.Search is already set from Bind)
	// Use GetAll for export without pagination
	list, err := u.repo.GetAll(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Create Excel file
	f := excelize.NewFile()
	sheet := "Groups"
	f.SetSheetName("Sheet1", sheet)

	// Header
	f.SetCellValue(sheet, "A1", "Kode Golongan")
	f.SetCellValue(sheet, "B1", "Nama Golongan")
	f.SetCellValue(sheet, "C1", "Update Date")

	// Rows
	for i, g := range list {
		row := i + 2
		f.SetCellValue(sheet, "A"+strconv.Itoa(row), g.GroupCode)
		f.SetCellValue(sheet, "B"+strconv.Itoa(row), g.Name)
		f.SetCellValue(sheet, "C"+strconv.Itoa(row), g.UpdatedAt.Local().Format("2006/01/02"))
	}

	// Write to buffer and return bytes
	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
