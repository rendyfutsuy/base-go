package http

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rendyfutsuy/base-go/constants"
	"github.com/rendyfutsuy/base-go/helpers/response"
	"github.com/xuri/excelize/v2"
)

// ImportUsersFromExcel godoc
// @Summary		Import users from Excel file
// @Description	Import multiple users from an Excel file (.xlsx or .xls). The Excel file must have columns: email, full_name, username, nik, role_name. Validates for duplicate email, username, and NIK. Returns HTTP 400 if any row fails with detailed error information per row.
// @Tags			User Management
// @Accept			multipart/form-data
// @Produce		json
// @Security		BearerAuth
// @Param			file	formData	file	true	"Excel file (.xlsx or .xls) with columns: email, full_name, username, nik, role_name"
// @Success		200		{object}	response.NonPaginationResponse{data=dto.ResImportUsers}	"Successfully imported all users"
// @Failure		400		{object}	response.NonPaginationResponse{data=dto.ResImportUsers}	"Bad request - one or more rows failed validation. Response contains details for each row including row number, username, status, and error message"
// @Failure		401		{object}	response.NonPaginationResponse	"Unauthorized"
// @Failure		500		{object}	response.NonPaginationResponse	"Internal server error"
// @Router			/v1/user-management/user/import [post]
func (handler *UserManagementHandler) ImportUsersFromExcel(c echo.Context) error {
	// initialize context from echo
	ctx := c.Request().Context()

	// Get uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, constants.UserImportFileNotFound))
	}

	// Validate file extension
	ext := filepath.Ext(file.Filename)
	if ext != ".xlsx" && ext != ".xls" {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, constants.UserImportInvalidFileFormat))
	}

	// Open uploaded file
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.SetErrorResponse(http.StatusInternalServerError, fmt.Sprintf("%s: %v", constants.UserImportFileOpenFailed, err)))
	}
	defer src.Close()

	// Create temporary file
	tempDir := os.TempDir()
	tempFileName := fmt.Sprintf("import_users_%d%s", time.Now().Unix(), ext)
	tempFilePath := filepath.Join(tempDir, tempFileName)

	// Create temporary file
	dst, err := os.Create(tempFilePath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.SetErrorResponse(http.StatusInternalServerError, fmt.Sprintf("%s: %v", constants.UserImportTempFileCreateFailed, err)))
	}
	defer dst.Close()
	defer os.Remove(tempFilePath) // Clean up temporary file

	// Copy uploaded file to temporary file
	_, err = io.Copy(dst, src)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.SetErrorResponse(http.StatusInternalServerError, fmt.Sprintf("%s: %v", constants.UserImportFileSaveFailed, err)))
	}

	// Process Excel file
	res, err := handler.UserUseCase.ImportUsersFromExcel(ctx, tempFilePath)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	// If there are failed rows, return HTTP 400 with error details
	if res.FailedCount > 0 {
		resp := response.NonPaginationResponse{}
		resp, _ = resp.SetResponse(res)
		if res.FailedCount > 0 {
			resp.Message = constants.UserImportFailedPartial
			if res.SuccessCount == 0 {
				resp.Message = constants.UserImportFileSaveFailed
				resp.Status = http.StatusBadRequest
			}
		}
		return c.JSON(http.StatusBadRequest, resp)
	}

	// Return success response if all rows processed successfully
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
// @Failure		400	{object}	response.NonPaginationResponse	"Bad request"
// @Failure		401	{object}	response.NonPaginationResponse	"Unauthorized"
// @Failure		500	{object}	response.NonPaginationResponse	"Internal server error"
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
	c.Response().Header().Set(constants.FieldContentType, constants.ExcelContent)
	c.Response().Header().Set(constants.FieldContentDisposition, constants.ExcelContentDisposition("user_import_template.xlsx"))

	// Write Excel file to response
	err = f.Write(c.Response().Writer)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.SetErrorResponse(http.StatusInternalServerError, fmt.Sprintf("%s: %v", constants.UserImportTemplateCreateFailed, err)))
	}

	return nil
}
