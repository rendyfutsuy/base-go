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

	// Process data rows (skip header row at index 0)
	var results []dto.ResImportUserExcel
	successCount := 0
	failedCount := 0

	for i := 1; i < len(rows); i++ {
		row := rows[i]
		rowNum := i + 1 // Excel row number (1-indexed)

		result := dto.ResImportUserExcel{
			Row: rowNum,
		}

		// Check if row has enough columns
		if len(row) < 5 {
			result.Success = false
			result.ErrorMessage = "Row tidak memiliki cukup kolom (minimal 5 kolom: email, full_name, username, nik, role_name)"
			results = append(results, result)
			failedCount++
			continue
		}

		// Parse row data
		email := strings.TrimSpace(row[0])
		fullName := strings.TrimSpace(row[1])
		username := strings.TrimSpace(row[2])
		nik := strings.TrimSpace(row[3])
		roleName := strings.TrimSpace(row[4])

		result.Email = email
		result.FullName = fullName
		result.Username = username
		result.Nik = nik
		result.RoleName = roleName

		// Validate required fields
		var validationErrors []string
		if email == "" {
			validationErrors = append(validationErrors, "email tidak boleh kosong")
		}
		if fullName == "" {
			validationErrors = append(validationErrors, "full_name tidak boleh kosong")
		}
		if username == "" {
			validationErrors = append(validationErrors, "username tidak boleh kosong")
		}
		if nik == "" {
			validationErrors = append(validationErrors, "nik tidak boleh kosong")
		}
		if roleName == "" {
			validationErrors = append(validationErrors, "role_name tidak boleh kosong")
		}

		if len(validationErrors) > 0 {
			result.Success = false
			result.ErrorMessage = strings.Join(validationErrors, "; ")
			results = append(results, result)
			failedCount++
			continue
		}

		// Validate email format (basic check)
		if !strings.Contains(email, "@") {
			result.Success = false
			result.ErrorMessage = "Format email tidak valid"
			results = append(results, result)
			failedCount++
			continue
		}

		// Check if email already exists
		emailNotDuplicated, err := u.userRepo.EmailIsNotDuplicated(ctx, email, uuid.Nil)
		if err != nil {
			result.Success = false
			result.ErrorMessage = fmt.Sprintf("Error checking email: %v", err)
			results = append(results, result)
			failedCount++
			continue
		}
		if !emailNotDuplicated {
			result.Success = false
			result.ErrorMessage = "Email sudah terdaftar di database"
			results = append(results, result)
			failedCount++
			continue
		}

		// Check if username already exists
		usernameNotDuplicated, err := u.userRepo.UsernameIsNotDuplicated(ctx, username, uuid.Nil)
		if err != nil {
			result.Success = false
			result.ErrorMessage = fmt.Sprintf("Error checking username: %v", err)
			results = append(results, result)
			failedCount++
			continue
		}
		if !usernameNotDuplicated {
			result.Success = false
			result.ErrorMessage = "Username sudah terdaftar di database"
			results = append(results, result)
			failedCount++
			continue
		}

		// Check if NIK already exists
		nikNotDuplicated, err := u.userRepo.NikIsNotDuplicated(ctx, nik, uuid.Nil)
		if err != nil {
			result.Success = false
			result.ErrorMessage = fmt.Sprintf("Error checking NIK: %v", err)
			results = append(results, result)
			failedCount++
			continue
		}
		if !nikNotDuplicated {
			result.Success = false
			result.ErrorMessage = "NIK sudah terdaftar di database"
			results = append(results, result)
			failedCount++
			continue
		}

		// Get role by name
		role, err := u.roleManagement.GetRoleByName(ctx, roleName)
		if err != nil {
			result.Success = false
			result.ErrorMessage = fmt.Sprintf("Role dengan nama '%s' tidak ditemukan", roleName)
			results = append(results, result)
			failedCount++
			continue
		}

		// Create user
		userDb := dto.ToDBCreateUser{
			FullName: fullName,
			Username: username,
			RoleId:   role.ID,
			Email:    email,
			Nik:      nik,
			IsActive: true,
			Gender:   "", // Default empty, bisa ditambahkan di Excel jika diperlukan
			Password: passwordTemplate,
		}

		_, err = u.userRepo.CreateUser(ctx, userDb)
		if err != nil {
			result.Success = false
			result.ErrorMessage = fmt.Sprintf("Error creating user: %v", err)
			results = append(results, result)
			failedCount++
			continue
		}

		// Success
		result.Success = true
		results = append(results, result)
		successCount++
	}

	return &dto.ResImportUsers{
		TotalRows:    len(rows) - 1, // Exclude header
		SuccessCount: successCount,
		FailedCount:  failedCount,
		Results:      results,
	}, nil
}

