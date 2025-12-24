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
	mod "github.com/rendyfutsuy/base-go/modules/expedition"
	"github.com/rendyfutsuy/base-go/modules/expedition/dto"
	"github.com/rendyfutsuy/base-go/utils"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type expeditionUsecase struct {
	repo mod.Repository
}

func NewExpeditionUsecase(repo mod.Repository) mod.Usecase {
	return &expeditionUsecase{repo: repo}
}

func (u *expeditionUsecase) Create(ctx context.Context, reqBody *dto.ReqCreateExpedition, authId string) (*models.Expedition, error) {
	// Check if expedition name already exists
	exists, err := u.repo.ExistsByExpeditionName(ctx, reqBody.ExpeditionName, uuid.Nil)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New(constants.ExpeditionNameAlreadyExists)
	}

	return u.repo.Create(ctx, mod.CreateExpeditionParams{
		ExpeditionName: reqBody.ExpeditionName,
		Address:        reqBody.Address,
		TelpNumbers:    reqBody.TelpNumbers,
		PhoneNumbers:   reqBody.PhoneNumbers,
		Notes:          reqBody.Notes,
		CreatedBy:      authId,
	})
}

func (u *expeditionUsecase) Update(ctx context.Context, id string, reqBody *dto.ReqUpdateExpedition, authId string) (*models.Expedition, error) {
	eid, err := utils.StringToUUID(id)
	if err != nil {
		return nil, err
	}

	// Check if expedition name already exists (excluding current id)
	exists, err := u.repo.ExistsByExpeditionName(ctx, reqBody.ExpeditionName, eid)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New(constants.ExpeditionNameAlreadyExists)
	}

	res, err := u.repo.Update(ctx, eid, mod.UpdateExpeditionParams{
		ExpeditionName: reqBody.ExpeditionName,
		Address:        reqBody.Address,
		TelpNumbers:    reqBody.TelpNumbers,
		PhoneNumbers:   reqBody.PhoneNumbers,
		Notes:          reqBody.Notes,
		UpdatedBy:      authId,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf(constants.ExpeditionNotFound, id)
		}
		return nil, err
	}
	return res, nil
}

func (u *expeditionUsecase) Delete(ctx context.Context, id string, authId string) error {
	eid, err := utils.StringToUUID(id)
	if err != nil {
		return err
	}
	return u.repo.Delete(ctx, eid, authId)
}

func (u *expeditionUsecase) GetByID(ctx context.Context, id string) (*models.Expedition, error) {
	eid, err := utils.StringToUUID(id)
	if err != nil {
		return nil, err
	}
	return u.repo.GetByID(ctx, eid)
}

func (u *expeditionUsecase) GetContactsByExpeditionID(ctx context.Context, id string) ([]models.ExpeditionContact, error) {
	eid, err := utils.StringToUUID(id)
	if err != nil {
		return nil, err
	}
	return u.repo.GetContactsByExpeditionID(ctx, eid)
}

func (u *expeditionUsecase) GetIndex(ctx context.Context, req request.PageRequest, filter dto.ReqExpeditionIndexFilter) ([]models.Expedition, int, error) {
	return u.repo.GetIndex(ctx, req, filter)
}

func (u *expeditionUsecase) GetAll(ctx context.Context, filter dto.ReqExpeditionIndexFilter) ([]models.Expedition, error) {
	return u.repo.GetAll(ctx, filter)
}

func (u *expeditionUsecase) Export(ctx context.Context, filter dto.ReqExpeditionIndexFilter) ([]byte, error) {
	// Use GetAllForExport for export with all phone numbers
	list, err := u.repo.GetAllForExport(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Find maximum number of phone numbers and telp numbers to determine header width
	maxPhoneNumbers := 0
	maxTelpNumbers := 0
	for _, expedition := range list {
		if len(expedition.PhoneNumbers) > maxPhoneNumbers {
			maxPhoneNumbers = len(expedition.PhoneNumbers)
		}
		if len(expedition.TelpNumbers) > maxTelpNumbers {
			maxTelpNumbers = len(expedition.TelpNumbers)
		}
	}
	// Ensure at least 1 column for each type
	if maxPhoneNumbers == 0 {
		maxPhoneNumbers = 1
	}
	if maxTelpNumbers == 0 {
		maxTelpNumbers = 1
	}

	// Create Excel file
	f := excelize.NewFile()
	sheet := "Expeditions"
	f.SetSheetName("Sheet1", sheet)

	// Helper function to convert column index (0-based) to Excel column letter
	colToLetter := func(col int) string {
		result := ""
		col++ // Convert to 1-based for calculation
		for col > 0 {
			col--
			result = string(rune('A'+(col%26))) + result
			col = col / 26
		}
		return result
	}

	// Header columns: Kode Ekspedisi, Nama Ekspedisi, Alamat Ekspedisi, No HP (merged), No Telp (merged), Update Date
	headerCol := 0

	// Set headers
	f.SetCellValue(sheet, colToLetter(headerCol)+"1", "Kode Ekspedisi")
	headerCol++
	f.SetCellValue(sheet, colToLetter(headerCol)+"1", "Nama Ekspedisi")
	headerCol++
	f.SetCellValue(sheet, colToLetter(headerCol)+"1", "Alamat Ekspedisi")
	headerCol++

	// Merge and center "No HP" header
	noHPStartCol := headerCol
	noHPEndCol := headerCol + maxPhoneNumbers - 1
	noHPStartCell := colToLetter(noHPStartCol) + "1"
	noHPEndCell := colToLetter(noHPEndCol) + "1"
	f.SetCellValue(sheet, noHPStartCell, "No HP")

	// Merge cells for "No HP" header
	if maxPhoneNumbers > 1 {
		if err := f.MergeCell(sheet, noHPStartCell, noHPEndCell); err != nil {
			return nil, err
		}
	}

	headerCol = noHPEndCol + 1

	// Merge and center "No Telp" header
	noTelpStartCol := headerCol
	noTelpEndCol := headerCol + maxTelpNumbers - 1
	noTelpStartCell := colToLetter(noTelpStartCol) + "1"
	noTelpEndCell := colToLetter(noTelpEndCol) + "1"
	f.SetCellValue(sheet, noTelpStartCell, "No Telp")

	// Merge cells for "No Telp" header
	if maxTelpNumbers > 1 {
		if err := f.MergeCell(sheet, noTelpStartCell, noTelpEndCell); err != nil {
			return nil, err
		}
	}

	headerCol = noTelpEndCol + 1
	f.SetCellValue(sheet, colToLetter(headerCol)+"1", "Update Date")

	// Rows
	for i, expedition := range list {
		row := i + 2
		col := 0

		// Helper function to set cell value
		setCell := func(value interface{}) {
			cell := colToLetter(col) + strconv.Itoa(row)
			f.SetCellValue(sheet, cell, value)
			col++
		}

		setCell(expedition.ExpeditionCode)
		setCell(expedition.ExpeditionName)
		setCell(expedition.Address)

		// Set phone numbers (HP) in separate cells
		phoneColStart := col
		if len(expedition.PhoneNumbers) > 0 {
			for _, phoneNumber := range expedition.PhoneNumbers {
				setCell(phoneNumber)
			}
		} else {
			setCell("-")
		}
		// Fill remaining cells with empty string if phone numbers less than max
		for col < phoneColStart+maxPhoneNumbers {
			setCell("-")
		}

		// Set telp numbers in separate cells
		telpColStart := col
		if len(expedition.TelpNumbers) > 0 {
			for _, telpNumber := range expedition.TelpNumbers {
				setCell(telpNumber)
			}
		} else {
			setCell("-")
		}
		// Fill remaining cells with empty string if telp numbers less than max
		for col < telpColStart+maxTelpNumbers {
			setCell("-")
		}

		setCell(expedition.UpdatedAt.Local().Format("2006/01/02"))
	}

	// Calculate total columns and rows for border styling
	// Columns: Kode Ekspedisi, Nama Ekspedisi, Alamat Ekspedisi, No HP (maxPhoneNumbers), No Telp (maxTelpNumbers), Update Date
	totalCols := 3 + maxPhoneNumbers + maxTelpNumbers + 1 // 3 fixed + maxPhoneNumbers + maxTelpNumbers + 1
	totalRows := len(list) + 1                            // 1 header row + data rows

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
	endCell := colToLetter(totalCols-1) + strconv.Itoa(totalRows)
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
	headerEndCell := colToLetter(totalCols-1) + "1"
	if err := f.SetCellStyle(sheet, "A1", headerEndCell, headerStyle); err != nil {
		return nil, err
	}

	// Update "No HP" and "No Telp" header styles to include bold and border
	noHPHeaderStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: borderDefinition,
	})
	if err != nil {
		return nil, err
	}
	if err := f.SetCellStyle(sheet, noHPStartCell, noHPEndCell, noHPHeaderStyle); err != nil {
		return nil, err
	}
	if err := f.SetCellStyle(sheet, noTelpStartCell, noTelpEndCell, noHPHeaderStyle); err != nil {
		return nil, err
	}

	// Write to buffer and return bytes
	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
