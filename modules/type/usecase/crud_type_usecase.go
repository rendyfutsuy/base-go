package usecase

import (
	"errors"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rendyfutsuy/base-go/constants"
	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/models"
	subGroupRepo "github.com/rendyfutsuy/base-go/modules/sub-group"
	mod "github.com/rendyfutsuy/base-go/modules/type"
	"github.com/rendyfutsuy/base-go/modules/type/dto"
	"github.com/rendyfutsuy/base-go/utils"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type typeUsecase struct {
	repo         mod.Repository
	subGroupRepo subGroupRepo.Repository
}

func NewTypeUsecase(repo mod.Repository, subGroupRepo subGroupRepo.Repository) mod.Usecase {
	return &typeUsecase{repo: repo, subGroupRepo: subGroupRepo}
}

func (u *typeUsecase) Create(c echo.Context, reqBody *dto.ReqCreateType, authId string) (*models.Type, error) {
	ctx := c.Request().Context()
	user := c.Get("user")
	userID := ""
	if user != nil {
		if userModel, ok := user.(models.User); ok {
			userID = userModel.ID.String()
		}
	}

	// Check if subgroup_id exists
	if reqBody.SubgroupID != uuid.Nil {
		subGroupObject, err := u.subGroupRepo.GetByID(ctx, reqBody.SubgroupID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New(constants.TypeSubGroupNotFound)
			}
			return nil, err
		}
		// Additional check: ensure subGroupObject is valid
		if subGroupObject == nil || subGroupObject.ID == uuid.Nil {
			return nil, errors.New(constants.TypeSubGroupNotFound)
		}
	}

	exists, err := u.repo.ExistsByNameInSubgroup(ctx, reqBody.SubgroupID, reqBody.Name, uuid.Nil)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New(constants.TypeNameAlreadyExists)
	}
	return u.repo.Create(ctx, reqBody.SubgroupID, reqBody.Name, userID)
}

func (u *typeUsecase) Update(c echo.Context, id string, reqBody *dto.ReqUpdateType, authId string) (*models.Type, error) {
	ctx := c.Request().Context()
	tid, err := utils.StringToUUID(id)
	if err != nil {
		return nil, err
	}
	user := c.Get("user")
	userID := ""
	if user != nil {
		if userModel, ok := user.(models.User); ok {
			userID = userModel.ID.String()
		}
	}

	// Check if subgroup_id exists
	if reqBody.SubgroupID != uuid.Nil {
		subGroupObject, err := u.subGroupRepo.GetByID(ctx, reqBody.SubgroupID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New(constants.TypeSubGroupNotFound)
			}
			return nil, err
		}
		// Additional check: ensure subGroupObject is valid
		if subGroupObject == nil || subGroupObject.ID == uuid.Nil {
			return nil, errors.New(constants.TypeSubGroupNotFound)
		}
	}

	exists, err := u.repo.ExistsByNameInSubgroup(ctx, reqBody.SubgroupID, reqBody.Name, tid)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New(constants.TypeNameAlreadyExists)
	}
	return u.repo.Update(ctx, tid, reqBody.SubgroupID, reqBody.Name, userID)
}

func (u *typeUsecase) Delete(c echo.Context, id string, authId string) error {
	ctx := c.Request().Context()
	tid, err := utils.StringToUUID(id)
	if err != nil {
		return err
	}
	user := c.Get("user")
	userID := ""
	if user != nil {
		if userModel, ok := user.(models.User); ok {
			userID = userModel.ID.String()
		}
	}
	return u.repo.Delete(ctx, tid, userID)
}

func (u *typeUsecase) GetByID(c echo.Context, id string) (*models.Type, error) {
	ctx := c.Request().Context()
	tid, err := utils.StringToUUID(id)
	if err != nil {
		return nil, err
	}
	return u.repo.GetByID(ctx, tid)
}

func (u *typeUsecase) GetIndex(c echo.Context, req request.PageRequest, filter dto.ReqTypeIndexFilter) ([]models.Type, int, error) {
	ctx := c.Request().Context()
	// Search is already set in req.Search from PageRequest middleware
	return u.repo.GetIndex(ctx, req, filter)
}

func (u *typeUsecase) GetAll(c echo.Context, filter dto.ReqTypeIndexFilter) ([]models.Type, error) {
	ctx := c.Request().Context()
	return u.repo.GetAll(ctx, filter)
}

func (u *typeUsecase) Export(c echo.Context, filter dto.ReqTypeIndexFilter) ([]byte, error) {
	ctx := c.Request().Context()
	// Use GetAll for export without pagination
	list, err := u.repo.GetAll(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Create Excel file
	f := excelize.NewFile()
	sheet := "Types"
	f.SetSheetName("Sheet1", sheet)

	// Header
	f.SetCellValue(sheet, "A1", "Kode Jenis")
	f.SetCellValue(sheet, "B1", "Nama Sub Golongan")
	f.SetCellValue(sheet, "C1", "Nama Jenis")
	f.SetCellValue(sheet, "D1", "Update Date")

	// Rows
	for i, t := range list {
		row := i + 2
		f.SetCellValue(sheet, "A"+strconv.Itoa(row), t.TypeCode)
		f.SetCellValue(sheet, "B"+strconv.Itoa(row), t.SubgroupName)
		f.SetCellValue(sheet, "C"+strconv.Itoa(row), t.Name)
		f.SetCellValue(sheet, "D"+strconv.Itoa(row), t.UpdatedAt.Local().Format("2006/01/02"))
	}

	// Write to buffer and return bytes
	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
