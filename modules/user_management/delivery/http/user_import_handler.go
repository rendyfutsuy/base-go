package http

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rendyfutsuy/base-go/helpers/response"
	"github.com/xuri/excelize/v2"
)

// ImportUsersFromExcel godoc
// @Summary		Import users from Excel file
// @Description	Import multiple users from an Excel file (.xlsx or .xls). The Excel file must have columns: email, full_name, username, nik, role_name. Validates for duplicate email, username, and NIK.
// @Tags			User Management
// @Accept			multipart/form-data
// @Produce		json
// @Security		BearerAuth
// @Param			file	formData	file	true	"Excel file (.xlsx or .xls) with columns: email, full_name, username, nik, role_name"
// @Success		200		{object}	response.NonPaginationResponse{data=dto.ResImportUsers}	"Successfully imported users"
// @Failure		400		{object}	ResponseError	"Bad request - invalid file or validation error"
// @Failure		401		{object}	ResponseError	"Unauthorized"
// @Failure		500		{object}	ResponseError	"Internal server error"
// @Router			/v1/user-management/user/import [post]
func (handler *UserManagementHandler) ImportUsersFromExcel(c echo.Context) error {
	// Get uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "File tidak ditemukan. Gunakan field 'file' untuk upload Excel file"})
	}

	// Validate file extension
	ext := filepath.Ext(file.Filename)
	if ext != ".xlsx" && ext != ".xls" {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "File harus berformat .xlsx atau .xls"})
	}

	// Open uploaded file
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{Message: fmt.Sprintf("Gagal membuka file: %v", err)})
	}
	defer src.Close()

	// Create temporary file
	tempDir := os.TempDir()
	tempFileName := fmt.Sprintf("import_users_%d%s", time.Now().Unix(), ext)
	tempFilePath := filepath.Join(tempDir, tempFileName)

	// Create temporary file
	dst, err := os.Create(tempFilePath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{Message: fmt.Sprintf("Gagal membuat file temporary: %v", err)})
	}
	defer dst.Close()
	defer os.Remove(tempFilePath) // Clean up temporary file

	// Copy uploaded file to temporary file
	_, err = io.Copy(dst, src)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{Message: fmt.Sprintf("Gagal menyimpan file: %v", err)})
	}

	// Process Excel file
	res, err := handler.UserUseCase.ImportUsersFromExcel(c, tempFilePath)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	// Return response
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(res)

	return c.JSON(http.StatusOK, resp)
}

// DownloadUserImportTemplate godoc
// @Summary		Download user import Excel template
// @Description	Download Excel template file for importing users. Template contains columns: email, full_name, username, nik, role_name with example data.
// @Tags			User Management
// @Accept			json
// @Produce		application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Security		BearerAuth
// @Success		200	{file}		file	"Excel template file"
// @Failure		400	{object}	ResponseError	"Bad request"
// @Failure		401	{object}	ResponseError	"Unauthorized"
// @Failure		500	{object}	ResponseError	"Internal server error"
// @Router			/v1/user-management/user/import/template [get]
func (handler *UserManagementHandler) DownloadUserImportTemplate(c echo.Context) error {
	// Create new Excel file
	f := excelize.NewFile()
	defer f.Close()

	// Set sheet name
	sheetName := "Import Users"
	f.SetSheetName("Sheet1", sheetName)

	// Set header row
	headers := []string{"Email", "Full Name", "Username", "NIK", "Role Name"}
	for i, header := range headers {
		cellName := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cellName, header)
	}

	// Style header row
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 12,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#E8E8E8"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err == nil {
		f.SetCellStyle(sheetName, "A1", fmt.Sprintf("%c1", 'A'+len(headers)-1), headerStyle)
	}

	// Set column widths
	f.SetColWidth(sheetName, "A", "A", 30) // Email
	f.SetColWidth(sheetName, "B", "B", 30) // Full Name
	f.SetColWidth(sheetName, "C", "C", 20) // Username
	f.SetColWidth(sheetName, "D", "D", 20) // NIK
	f.SetColWidth(sheetName, "E", "E", 25) // Role Name

	// Add example row
	exampleRow := []interface{}{"user@example.com", "John Doe", "johndoe", "1234567890123456", "Super Admin"}
	for i, value := range exampleRow {
		cellName := fmt.Sprintf("%c2", 'A'+i)
		f.SetCellValue(sheetName, cellName, value)
	}

	// Set response headers
	c.Response().Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Response().Header().Set("Content-Disposition", "attachment; filename=user_import_template.xlsx")

	// Write Excel file to response
	err = f.Write(c.Response().Writer)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{Message: fmt.Sprintf("Gagal membuat template: %v", err)})
	}

	return nil
}

