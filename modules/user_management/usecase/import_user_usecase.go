package usecase

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rendyfutsuy/base-go/modules/user_management/dto"
	"github.com/rendyfutsuy/base-go/utils"
	"github.com/xuri/excelize/v2"
)

func (u *userUsecase) ImportUsersFromExcel(c echo.Context, filePath string) (res *dto.ResImportUsers, err error) {
	ctx := c.Request().Context()

	// Open Excel file
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		utils.Logger.Error(err.Error())
		return nil, fmt.Errorf("failed to open Excel file: %v", err)
	}
	defer f.Close()

	// Get the first sheet
	sheetName := f.GetSheetName(0)
	if sheetName == "" {
		return nil, fmt.Errorf("Excel file has no sheets")
	}

	// Read all rows
	rows, err := f.GetRows(sheetName)
	if err != nil {
		utils.Logger.Error(err.Error())
		return nil, fmt.Errorf("failed to read Excel file: %v", err)
	}

	if len(rows) < 2 {
		return nil, fmt.Errorf("Excel file must have at least header row and one data row")
	}

	// Get password default from config
	passwordTemplate := "temp"
	if utils.ConfigVars.Exists("user.default_password_template") {
		passwordTemplate = utils.ConfigVars.String("user.default_password_template")
	}

	totalRows := len(rows) - 1 // Exclude header

	// Phase 1: Parse all rows and collect data for batch validation
	type parsedRowData struct {
		RowNum   int
		Email    string
		FullName string
		Username string
		Nik      string
		RoleName string
		Result   dto.ResImportUserExcel
	}

	parsedRows := make([]parsedRowData, 0, totalRows)

	// Parse all rows first
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		rowNum := i + 1 // Excel row number (1-indexed)

		parsedRow := parsedRowData{
			RowNum: rowNum,
			Result: dto.ResImportUserExcel{
				Row: rowNum,
			},
		}

		// Parse row data (try to parse even if incomplete)
		if len(row) > 0 {
			parsedRow.Email = strings.TrimSpace(row[0])
		}
		if len(row) > 1 {
			parsedRow.FullName = strings.TrimSpace(row[1])
		}
		if len(row) > 2 {
			parsedRow.Username = strings.TrimSpace(row[2])
		}
		if len(row) > 3 {
			parsedRow.Nik = strings.TrimSpace(row[3])
		}
		if len(row) > 4 {
			parsedRow.RoleName = strings.TrimSpace(row[4])
		}

		parsedRow.Result.Username = parsedRow.Username
		parsedRows = append(parsedRows, parsedRow)
	}

	// Phase 2: Collect all emails, usernames, niks, and role names for batch validation
	emails := make([]string, 0, totalRows)
	usernames := make([]string, 0, totalRows)
	niks := make([]string, 0, totalRows)
	roleNamesMap := make(map[string]bool) // Set untuk unique role names

	for i := range parsedRows {
		parsedRow := &parsedRows[i]

		// Basic validation first
		// Note: parsedRow.RowNum is 1-indexed Excel row number, but rows array is 0-indexed
		rowIndex := parsedRow.RowNum - 1
		if rowIndex >= 0 && rowIndex < len(rows) && len(rows[rowIndex]) < 5 {
			parsedRow.Result.Success = false
			parsedRow.Result.Status = "failed"
			parsedRow.Result.ErrorMessage = "Row tidak memiliki cukup kolom (minimal 5 kolom: email, full_name, username, nik, role_name)"
			continue
		}

		// Validate required fields
		var validationErrors []string
		if parsedRow.Email == "" {
			validationErrors = append(validationErrors, "email tidak boleh kosong")
		}
		if parsedRow.FullName == "" {
			validationErrors = append(validationErrors, "full_name tidak boleh kosong")
		}
		if parsedRow.Username == "" {
			validationErrors = append(validationErrors, "username tidak boleh kosong")
		}
		if parsedRow.Nik == "" {
			validationErrors = append(validationErrors, "nik tidak boleh kosong")
		}
		if parsedRow.RoleName == "" {
			validationErrors = append(validationErrors, "role_name tidak boleh kosong")
		}

		if len(validationErrors) > 0 {
			parsedRow.Result.Success = false
			parsedRow.Result.Status = "failed"
			parsedRow.Result.ErrorMessage = strings.Join(validationErrors, "; ")
			continue
		}

		// Validate email format (basic check)
		if !strings.Contains(parsedRow.Email, "@") {
			parsedRow.Result.Success = false
			parsedRow.Result.Status = "failed"
			parsedRow.Result.ErrorMessage = "Format email tidak valid"
			continue
		}

		// Collect for batch validation (only if passed basic validation)
		emails = append(emails, parsedRow.Email)
		usernames = append(usernames, parsedRow.Username)
		niks = append(niks, parsedRow.Nik)
		roleNamesMap[parsedRow.RoleName] = true
	}

	// Phase 3: Batch validation - check for duplicates (single query for all)
	duplicatedEmails, duplicatedUsernames, duplicatedNiks, err := u.userRepo.CheckBatchDuplication(ctx, emails, usernames, niks)
	if err != nil {
		return nil, fmt.Errorf("failed to check batch duplication: %v", err)
	}

	// Phase 4: Batch fetch roles (cache roles to avoid duplicate queries)
	roleNames := make([]string, 0, len(roleNamesMap))
	for roleName := range roleNamesMap {
		roleNames = append(roleNames, roleName)
	}

	roleMap := make(map[string]uuid.UUID)

	for _, roleName := range roleNames {
		role, err := u.roleManagement.GetRoleByName(ctx, roleName)
		if err != nil {
			// Role tidak ditemukan, akan ditangani per row nanti
			continue
		}
		roleMap[roleName] = role.ID
	}

	// Phase 5: Process each row with batch-validated data and prepare for batch insert
	var results []dto.ResImportUserExcel
	var validUsers []dto.ToDBCreateUser
	var validUserRowIndices []int // Track which rows correspond to valid users

	for i := range parsedRows {
		parsedRow := &parsedRows[i]
		result := &parsedRow.Result

		// Skip if already marked as failed in Phase 2
		if result.Status == "failed" {
			results = append(results, *result)
			continue
		}

		// Check duplicates using batch validation results
		if duplicatedEmails[parsedRow.Email] {
			result.Success = false
			result.Status = "failed"
			result.ErrorMessage = "Email sudah terdaftar di database"
			results = append(results, *result)
			continue
		}

		if duplicatedUsernames[parsedRow.Username] {
			result.Success = false
			result.Status = "failed"
			result.ErrorMessage = "Username sudah terdaftar di database"
			results = append(results, *result)
			continue
		}

		if duplicatedNiks[parsedRow.Nik] {
			result.Success = false
			result.Status = "failed"
			result.ErrorMessage = "NIK sudah terdaftar di database"
			results = append(results, *result)
			continue
		}

		// Check role
		roleId, exists := roleMap[parsedRow.RoleName]
		if !exists {
			result.Success = false
			result.Status = "failed"
			result.ErrorMessage = fmt.Sprintf("Role dengan nama '%s' tidak ditemukan", parsedRow.RoleName)
			results = append(results, *result)
			continue
		}

		// Add to valid users for batch insert
		userDb := dto.ToDBCreateUser{
			FullName: parsedRow.FullName,
			Username: parsedRow.Username,
			RoleId:   roleId,
			Email:    parsedRow.Email,
			Nik:      parsedRow.Nik,
			IsActive: true,
			Gender:   "",
			Password: passwordTemplate,
		}
		validUsers = append(validUsers, userDb)
		validUserRowIndices = append(validUserRowIndices, len(results))

		// Mark as success (will be validated after batch insert)
		result.Success = true
		result.Status = "success"
		results = append(results, *result)
	}

	// Phase 6: Batch insert valid users (single transaction)
	if len(validUsers) > 0 {
		err = u.userRepo.BulkCreateUsers(ctx, validUsers)
		if err != nil {
			// If batch insert fails, mark all pending users as failed
			for _, idx := range validUserRowIndices {
				if idx < len(results) {
					results[idx].Success = false
					results[idx].Status = "failed"
					results[idx].ErrorMessage = fmt.Sprintf("Error creating user in batch: %v", err)
				}
			}
		}
	}

	// Calculate final counts
	successCount := 0
	failedCount := 0
	for _, result := range results {
		if result.Status == "success" {
			successCount++
		} else {
			failedCount++
		}
	}

	return &dto.ResImportUsers{
		TotalRows:    totalRows,
		SuccessCount: successCount,
		FailedCount:  failedCount,
		Results:      results,
	}, nil
}
