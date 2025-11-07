package usecase

import (
	"errors"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rendyfutsuy/base-go/constants"
	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/models"
	groupRepo "github.com/rendyfutsuy/base-go/modules/group"
	mod "github.com/rendyfutsuy/base-go/modules/sub-group"
	"github.com/rendyfutsuy/base-go/modules/sub-group/dto"
	"github.com/rendyfutsuy/base-go/utils"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type subGroupUsecase struct {
	repo      mod.Repository
	groupRepo groupRepo.Repository
}

func NewSubGroupUsecase(repo mod.Repository, groupRepo groupRepo.Repository) mod.Usecase {
	return &subGroupUsecase{repo: repo, groupRepo: groupRepo}
}

func (u *subGroupUsecase) Create(c echo.Context, reqBody *dto.ReqCreateSubGroup, authId string) (*models.SubGroup, error) {
	ctx := c.Request().Context()

	// Get user ID from context
	user := c.Get("user")
	userID := ""
	if user != nil {
		if userModel, ok := user.(models.User); ok {
			userID = userModel.ID.String()
		}
	}

	// Check if goods_group_id exists
	if reqBody.GoodsGroupID != uuid.Nil {
		groupObject, err := u.groupRepo.GetByID(ctx, reqBody.GoodsGroupID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New(constants.SubGroupGoodsGroupNotFound)
			}
			return nil, err
		}
		// Additional check: ensure groupObject is valid
		if groupObject == nil || groupObject.ID == uuid.Nil {
			return nil, errors.New(constants.SubGroupGoodsGroupNotFound)
		}
	}

	exists, err := u.repo.ExistsByName(ctx, reqBody.GoodsGroupID, reqBody.Name, uuid.Nil)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New(constants.SubGroupNameAlreadyExists)
	}
	return u.repo.Create(ctx, reqBody.GoodsGroupID, reqBody.Name, userID)
}

func (u *subGroupUsecase) Update(c echo.Context, id string, reqBody *dto.ReqUpdateSubGroup, authId string) (*models.SubGroup, error) {
	ctx := c.Request().Context()
	sgid, err := utils.StringToUUID(id)
	if err != nil {
		return nil, err
	}

	// Get user ID from context
	user := c.Get("user")
	userID := ""
	if user != nil {
		if userModel, ok := user.(models.User); ok {
			userID = userModel.ID.String()
		}
	}

	// Check if goods_group_id exists
	if reqBody.GoodsGroupID != uuid.Nil {
		groupObject, err := u.groupRepo.GetByID(ctx, reqBody.GoodsGroupID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New(constants.SubGroupGoodsGroupNotFound)
			}
			return nil, err
		}
		// Additional check: ensure groupObject is valid
		if groupObject == nil || groupObject.ID == uuid.Nil {
			return nil, errors.New(constants.SubGroupGoodsGroupNotFound)
		}
	}

	exists, err := u.repo.ExistsByName(ctx, reqBody.GoodsGroupID, reqBody.Name, sgid)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New(constants.SubGroupNameAlreadyExists)
	}
	return u.repo.Update(ctx, sgid, reqBody.GoodsGroupID, reqBody.Name, userID)
}

func (u *subGroupUsecase) Delete(c echo.Context, id string, authId string) error {
	ctx := c.Request().Context()
	sgid, err := utils.StringToUUID(id)
	if err != nil {
		return err
	}

	// Get user ID from context
	user := c.Get("user")
	userID := ""
	if user != nil {
		if userModel, ok := user.(models.User); ok {
			userID = userModel.ID.String()
		}
	}

	return u.repo.Delete(ctx, sgid, userID)
}

func (u *subGroupUsecase) GetByID(c echo.Context, id string) (*models.SubGroup, error) {
	ctx := c.Request().Context()
	sgid, err := utils.StringToUUID(id)
	if err != nil {
		return nil, err
	}
	return u.repo.GetByID(ctx, sgid)
}

func (u *subGroupUsecase) GetIndex(c echo.Context, req request.PageRequest, filter dto.ReqSubGroupIndexFilter) ([]models.SubGroup, int, error) {
	ctx := c.Request().Context()
	// Search is already set in req.Search from PageRequest middleware
	return u.repo.GetIndex(ctx, req, filter)
}

func (u *subGroupUsecase) GetAll(c echo.Context, filter dto.ReqSubGroupIndexFilter) ([]models.SubGroup, error) {
	ctx := c.Request().Context()
	return u.repo.GetAll(ctx, filter)
}

func (u *subGroupUsecase) Export(c echo.Context, filter dto.ReqSubGroupIndexFilter) ([]byte, error) {
	ctx := c.Request().Context()
	// Use GetAll for export without pagination
	list, err := u.repo.GetAll(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Create Excel file
	f := excelize.NewFile()
	sheet := "SubGroups"
	f.SetSheetName("Sheet1", sheet)

	// Header
	f.SetCellValue(sheet, "A1", "Kode Sub Golongan")
	f.SetCellValue(sheet, "B1", "Nama Sub Golongan")
	f.SetCellValue(sheet, "C1", "Goods Group ID")
	f.SetCellValue(sheet, "D1", "Update Date")

	// Rows
	for i, sg := range list {
		row := i + 2
		f.SetCellValue(sheet, "A"+strconv.Itoa(row), sg.SubgroupCode)
		f.SetCellValue(sheet, "B"+strconv.Itoa(row), sg.Name)
		f.SetCellValue(sheet, "C"+strconv.Itoa(row), sg.GoodsGroupID.String())
		f.SetCellValue(sheet, "D"+strconv.Itoa(row), sg.UpdatedAt.Local().Format("2006/01/02"))
	}

	// Write to buffer and return bytes
	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
