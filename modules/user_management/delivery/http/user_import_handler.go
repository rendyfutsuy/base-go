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

