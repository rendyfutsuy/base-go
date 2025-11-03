package controllers

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/rendyfutsuy/base-go/constants"
	"github.com/rendyfutsuy/base-go/utils"
)

type HomepageData struct {
	Version        string
	LastUpdated    string
	ShowSwaggerLink bool
}

func DefaultHomepage(c echo.Context) error {
	// Get the path to the template file
	// Try to find the template file relative to the executable or current working directory
	var templatePath string

	// First, try relative path from current working directory
	templatePath = filepath.Join("modules", "homepage", "delivery", "http", "templates", "homepage.html")
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		// If not found, try relative to the source file location
		// This works during development when running from project root
		templatePath = filepath.Join(".", "modules", "homepage", "delivery", "http", "templates", "homepage.html")
		if _, err := os.Stat(templatePath); os.IsNotExist(err) {
			return c.HTML(http.StatusInternalServerError, "<h1>Error: Template file not found</h1>")
		}
	}

	// Parse template
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return c.HTML(http.StatusInternalServerError, "<h1>Error loading template</h1>")
	}

	// Prepare data
	appEnv := utils.ConfigVars.String("app_env")
	showSwaggerLink := appEnv == "development"
	
	data := HomepageData{
		Version:        constants.Version,
		LastUpdated:    "2025/11/03 14:32 WIB",
		ShowSwaggerLink: showSwaggerLink,
	}

	// Execute template
	err = tmpl.Execute(c.Response().Writer, data)
	if err != nil {
		return c.HTML(http.StatusInternalServerError, "<h1>Error rendering template</h1>")
	}

	return nil
}
