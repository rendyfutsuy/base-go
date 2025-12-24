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

func (u *subGroupUsecase) Create(ctx context.Context, reqBody *dto.ReqCreateSubGroup, userID string) (*models.SubGroup, error) {
	// Check if groups_id exists
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

func (u *subGroupUsecase) Update(ctx context.Context, id string, reqBody *dto.ReqUpdateSubGroup, userID string) (*models.SubGroup, error) {
	sgid, err := utils.StringToUUID(id)
	if err != nil {
		return nil, err
	}

	// Check if groups_id exists
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
	res, err := u.repo.Update(ctx, sgid, reqBody.GoodsGroupID, reqBody.Name, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf(constants.SubGroupNotFound, id)
		}
		return nil, err
	}
	return res, nil
}

func (u *subGroupUsecase) Delete(ctx context.Context, id string, userID string) error {
	sgid, err := utils.StringToUUID(id)
	if err != nil {
		return err
	}

	// Check if sub-group is still used in types
	exists, err := u.repo.ExistsInTypes(ctx, sgid)
	if err != nil {
		return err
	}
	if exists {
		return errors.New(constants.SubGroupStillUsedInTypes)
	}

	return u.repo.Delete(ctx, sgid, userID)
}

func (u *subGroupUsecase) ExistsInTypes(ctx context.Context, subGroupID string) (bool, error) {
	sgid, err := utils.StringToUUID(subGroupID)
	if err != nil {
		return false, err
	}
	return u.repo.ExistsInTypes(ctx, sgid)
}

func (u *subGroupUsecase) GetByID(ctx context.Context, id string) (*models.SubGroup, error) {
	sgid, err := utils.StringToUUID(id)
	if err != nil {
		return nil, err
	}
	return u.repo.GetByID(ctx, sgid)
}

func (u *subGroupUsecase) GetIndex(ctx context.Context, req request.PageRequest, filter dto.ReqSubGroupIndexFilter) ([]models.SubGroup, int, error) {
	// Search is already set in req.Search from PageRequest middleware
	return u.repo.GetIndex(ctx, req, filter)
}

func (u *subGroupUsecase) GetAll(ctx context.Context, filter dto.ReqSubGroupIndexFilter) ([]models.SubGroup, error) {
	return u.repo.GetAll(ctx, filter)
}

func (u *subGroupUsecase) Export(ctx context.Context, filter dto.ReqSubGroupIndexFilter) ([]byte, error) {
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
	headers := []string{
		"Kode Sub Golongan",
		"Nama Sub Golongan",
		"Nama Golongan",
		"Updated Date",
	}

	// Set headers
	for i, header := range headers {
		cell := string(rune('A'+i)) + "1"
		f.SetCellValue(sheet, cell, header)
	}

	// Rows
	for i, sg := range list {
		row := i + 2
		col := 0

		// Helper function to set cell value
		setCell := func(value interface{}) {
			cell := string(rune('A'+col)) + strconv.Itoa(row)
			f.SetCellValue(sheet, cell, value)
			col++
		}

		setCell(sg.SubgroupCode)
		setCell(sg.Name)
		if sg.GoodsGroupName != "" {
			setCell(sg.GoodsGroupName)
		} else {
			setCell("-")
		}
		setCell(sg.UpdatedAt.Format("2006-01-02"))
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
	endCell := "D" + strconv.Itoa(totalRows)
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
	if err := f.SetCellStyle(sheet, "A1", "D1", headerStyle); err != nil {
		return nil, err
	}

	// Write to buffer and return bytes
	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
