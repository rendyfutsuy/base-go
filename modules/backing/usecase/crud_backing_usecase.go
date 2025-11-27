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
	mod "github.com/rendyfutsuy/base-go/modules/backing"
	"github.com/rendyfutsuy/base-go/modules/backing/dto"
	typeRepo "github.com/rendyfutsuy/base-go/modules/type"
	"github.com/rendyfutsuy/base-go/utils"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type backingUsecase struct {
	repo     mod.Repository
	typeRepo typeRepo.Repository
}

func NewBackingUsecase(repo mod.Repository, typeRepo typeRepo.Repository) mod.Usecase {
	return &backingUsecase{repo: repo, typeRepo: typeRepo}
}

func (u *backingUsecase) Create(c echo.Context, reqBody *dto.ReqCreateBacking, authId string) (*models.Backing, error) {
	ctx := c.Request().Context()
	user := c.Get("user")
	userID := ""
	if user != nil {
		if userModel, ok := user.(models.User); ok {
			userID = userModel.ID.String()
		}
	}

	// Check if type_id exists
	typeObject, err := u.typeRepo.GetByID(ctx, reqBody.TypeID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(constants.BackingTypeNotFound)
		}
		return nil, err
	}
	// Additional check: ensure typeObject is valid
	if typeObject == nil || typeObject.ID == uuid.Nil {
		return nil, errors.New(constants.BackingTypeNotFound)
	}

	exists, err := u.repo.ExistsByNameInType(ctx, reqBody.TypeID, reqBody.Name, uuid.Nil)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New(constants.BackingNameAlreadyExistsInType)
	}
	return u.repo.Create(ctx, reqBody.TypeID, reqBody.Name, userID)
}

func (u *backingUsecase) Update(c echo.Context, id string, reqBody *dto.ReqUpdateBacking, authId string) (*models.Backing, error) {
	ctx := c.Request().Context()
	bid, err := utils.StringToUUID(id)
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

	// Check if type_id exists
	typeObject, err := u.typeRepo.GetByID(ctx, reqBody.TypeID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(constants.BackingTypeNotFound)
		}
		return nil, err
	}
	// Additional check: ensure typeObject is valid
	if typeObject == nil || typeObject.ID == uuid.Nil {
		return nil, errors.New(constants.BackingTypeNotFound)
	}

	exists, err := u.repo.ExistsByNameInType(ctx, reqBody.TypeID, reqBody.Name, bid)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New(constants.BackingNameAlreadyExistsInType)
	}
	res, err := u.repo.Update(ctx, bid, reqBody.TypeID, reqBody.Name, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf(constants.BackingNotFound, id)
		}
		return nil, err
	}
	return res, nil
}

func (u *backingUsecase) Delete(c echo.Context, id string, authId string) error {
	ctx := c.Request().Context()
	bid, err := utils.StringToUUID(id)
	if err != nil {
		return err
	}
	return u.repo.Delete(ctx, bid)
}

func (u *backingUsecase) GetByID(c echo.Context, id string) (*models.Backing, error) {
	ctx := c.Request().Context()
	bid, err := utils.StringToUUID(id)
	if err != nil {
		return nil, err
	}
	return u.repo.GetByID(ctx, bid)
}

func (u *backingUsecase) GetIndex(c echo.Context, req request.PageRequest, filter dto.ReqBackingIndexFilter) ([]models.Backing, int, error) {
	ctx := c.Request().Context()
	// Search is already set in req.Search from PageRequest middleware
	return u.repo.GetIndex(ctx, req, filter)
}

func (u *backingUsecase) GetAll(c echo.Context, filter dto.ReqBackingIndexFilter) ([]models.Backing, error) {
	ctx := c.Request().Context()
	return u.repo.GetAll(ctx, filter)
}

func (u *backingUsecase) Export(c echo.Context, filter dto.ReqBackingIndexFilter) ([]byte, error) {
	ctx := c.Request().Context()
	// Extract search from query param and set to filter (filter.Search is already set from Bind)
	// Use GetAll for export without pagination
	list, err := u.repo.GetAll(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Create Excel file
	f := excelize.NewFile()
	sheet := "Backings"
	f.SetSheetName("Sheet1", sheet)

	// Header
	f.SetCellValue(sheet, "A1", "Kode Backing")
	f.SetCellValue(sheet, "B1", "Nama Backing")
	f.SetCellValue(sheet, "C1", "Nama Jenis")
	f.SetCellValue(sheet, "D1", "Nama Sub Golongan")
	f.SetCellValue(sheet, "E1", "Nama Golongan")
	f.SetCellValue(sheet, "F1", "Updated Date")

	// Rows
	for i, b := range list {
		row := i + 2
		f.SetCellValue(sheet, "A"+strconv.Itoa(row), b.BackingCode)
		f.SetCellValue(sheet, "B"+strconv.Itoa(row), b.Name)
		f.SetCellValue(sheet, "C"+strconv.Itoa(row), b.TypeName)
		f.SetCellValue(sheet, "D"+strconv.Itoa(row), b.SubgroupName)
		f.SetCellValue(sheet, "E"+strconv.Itoa(row), b.GroupName)
		f.SetCellValue(sheet, "F"+strconv.Itoa(row), b.UpdatedAt.Local().Format("2006/01/02"))
	}

	// Calculate total rows for border styling
	totalRows := len(list) + 1 // 1 header row + data rows

	// Define border configuration
	borderDefinition := []excelize.Border{
		{Type: "left", Color: "000000", Style: 1},
		{Type: "top", Color: "000000", Style: 1},
		{Type: "bottom", Color: "000000", Style: 1},
		{Type: "right", Color: "000000", Style: 1},
	}

	// Create border style
	borderStyle, err := f.NewStyle(&excelize.Style{
		Border: borderDefinition,
	})
	if err != nil {
		return nil, err
	}

	// Apply border style to all cells
	startCell := "A1"
	endCell := "F" + strconv.Itoa(totalRows)
	if err := f.SetCellStyle(sheet, startCell, endCell, borderStyle); err != nil {
		return nil, err
	}

	// Create header style with bold font and border
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
		Border: borderDefinition,
	})
	if err != nil {
		return nil, err
	}

	// Apply header style to header row
	if err := f.SetCellStyle(sheet, "A1", "F1", headerStyle); err != nil {
		return nil, err
	}

	// Write to buffer and return bytes
	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
