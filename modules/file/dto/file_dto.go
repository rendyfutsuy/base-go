package dto

import "github.com/google/uuid"

type ToDBFile struct {
	Name        string
	FilePath    *string
	Description *string
}

type UploadInput struct {
	Data             []byte
	OriginalFileName string
	DestRoot         string  // root directory to save file
	ExtraPath        *string // extra path to save file. based on case used in module, e.g. "post/123"
	Description      *string
}

type AssignFileItem struct {
	FileID uuid.UUID
	Type   *string // file type, e.g. "image", "video", "document"
}

type AssignFilesToModule struct {
	ModuleID   uuid.UUID
	ModuleType string
	Items      []AssignFileItem
}
