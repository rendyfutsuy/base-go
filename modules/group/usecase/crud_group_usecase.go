package usecase

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/constants"
	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/models"
	mod "github.com/rendyfutsuy/base-go/modules/group"
	"github.com/rendyfutsuy/base-go/modules/group/dto"
	"github.com/rendyfutsuy/base-go/utils"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type groupUsecase struct {
	repo mod.Repository
}

func NewGroupUsecase(repo mod.Repository) mod.Usecase {
	return &groupUsecase{repo: repo}
}

func (u *groupUsecase) Create(ctx context.Context, reqBody *dto.ReqCreateGroup, userID string) (*models.GoodsGroup, error) {
	exists, err := u.repo.ExistsByName(ctx, reqBody.Name, uuid.Nil)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New(constants.GroupNameAlreadyExists)
	}

	return u.repo.Create(ctx, reqBody.Name, userID)
}

func (u *groupUsecase) Update(ctx context.Context, id string, reqBody *dto.ReqUpdateGroup, userID string) (*models.GoodsGroup, error) {
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

	res, err := u.repo.Update(ctx, gid, reqBody.Name, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf(constants.GroupNotFound, id)
		}
		return nil, err
	}
	return res, nil
}

func (u *groupUsecase) Delete(ctx context.Context, id string, userID string) error {
	gid, err := utils.StringToUUID(id)
	if err != nil {
		return err
	}

	// Check if group is still used in sub-groups (not deleted)
	exists, err := u.repo.ExistsInSubGroups(ctx, gid)
	if err != nil {
		return err
	}
	if exists {
		return errors.New(constants.GroupStillUsedInSubGroups)
	}

	return u.repo.Delete(ctx, gid, userID)
}

func (u *groupUsecase) GetByID(ctx context.Context, id string) (*models.GoodsGroup, error) {
	gid, err := utils.StringToUUID(id)
	if err != nil {
		return nil, err
	}
	return u.repo.GetByID(ctx, gid)
}

func (u *groupUsecase) GetIndex(ctx context.Context, req request.PageRequest, filter dto.ReqGroupIndexFilter) ([]models.GoodsGroup, int, error) {
	// Search is already set in req.Search from PageRequest middleware
	return u.repo.GetIndex(ctx, req, filter)
}

func (u *groupUsecase) GetAll(ctx context.Context, filter dto.ReqGroupIndexFilter) ([]models.GoodsGroup, error) {
	return u.repo.GetAll(ctx, filter)
}

func (u *groupUsecase) Export(ctx context.Context, filter dto.ReqGroupIndexFilter) ([]byte, error) {
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
	endCell := "C" + strconv.Itoa(totalRows)
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
	if err := f.SetCellStyle(sheet, "A1", "C1", headerStyle); err != nil {
		return nil, err
	}

	// Write to buffer and return bytes
	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (u *groupUsecase) ExistsInSubGroups(ctx context.Context, groupID uuid.UUID) (bool, error) {
	return u.repo.ExistsInSubGroups(ctx, groupID)
}
