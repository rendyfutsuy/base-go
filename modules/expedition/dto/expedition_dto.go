package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/constants"
	"github.com/rendyfutsuy/base-go/models"
)

// TelpNumberItem represents telp number with area code
type TelpNumberItem struct {
	AreaCode    *string `json:"area_code" form:"area_code"`
	PhoneNumber string  `json:"phone_number" form:"phone_number" validate:"required"`
}

type ReqCreateExpedition struct {
	ExpeditionName string           `form:"expedition_name" json:"expedition_name" validate:"required,max=255"`
	Address        string           `form:"address" json:"address" validate:"max=255"`
	TelpNumbers    []TelpNumberItem `form:"-" json:"telp_numbers" validate:"omitempty"`
	PhoneNumbers   []string         `form:"phone_numbers" json:"phone_numbers" validate:"omitempty"`
	Notes          *string          `form:"notes" json:"notes,omitempty"`
}

type ReqUpdateExpedition struct {
	ExpeditionName string           `form:"expedition_name" json:"expedition_name" validate:"required,max=255"`
	Address        string           `form:"address" json:"address" validate:"max=255"`
	TelpNumbers    []TelpNumberItem `form:"-" json:"telp_numbers"`
	PhoneNumbers   []string         `form:"phone_numbers" json:"phone_numbers"`
	Notes          *string          `form:"notes" json:"notes,omitempty"`
}

// ContactResponse represents contact response
type ContactResponse struct {
	ID          uuid.UUID `json:"id"`
	PhoneType   string    `json:"phone_type"`
	PhoneNumber string    `json:"phone_number"`
	IsPrimary   bool      `json:"is_primary"`
}

type RespExpedition struct {
	ID             uuid.UUID        `json:"id"`
	ExpeditionCode string           `json:"expedition_code"`
	ExpeditionName string           `json:"expedition_name"`
	Address        string           `json:"address"`
	TelpNumbers    []TelpNumberItem `json:"telp_numbers"`
	PhoneNumbers   []string         `json:"phone_numbers"`
	Notes          *string          `json:"notes,omitempty"`
	CreatedAt      string           `json:"created_at"`
	CreatedBy      string           `json:"created_by"`
	UpdatedAt      string           `json:"updated_at"`
	UpdatedBy      string           `json:"updated_by"`
	Deletable      bool             `json:"deletable"`
}

func ToRespExpedition(m models.Expedition, contacts []models.ExpeditionContact) RespExpedition {
	// Map contacts to response
	telpNumbers := []TelpNumberItem{}
	PhoneNumbers := []string{}
	for _, contact := range contacts {
		if contact.PhoneType == constants.ExpeditionContactTypeTelp {
			telpNumbers = append(telpNumbers, TelpNumberItem{
				AreaCode:    contact.AreaCode,
				PhoneNumber: contact.PhoneNumber,
			})
		} else {
			PhoneNumbers = append(PhoneNumbers, contact.PhoneNumber)
		}
	}

	return RespExpedition{
		ID:             m.ID,
		ExpeditionCode: m.ExpeditionCode,
		ExpeditionName: m.ExpeditionName,
		Address:        m.Address,
		TelpNumbers:    telpNumbers,
		PhoneNumbers:   PhoneNumbers,
		Notes:          m.Notes,
		CreatedAt:      m.CreatedAt.Format("2006-01-02 15:04:05"),
		CreatedBy:      m.CreatedBy,
		UpdatedAt:      m.UpdatedAt.Format("2006-01-02 15:04:05"),
		UpdatedBy:      m.UpdatedBy,
		Deletable:      m.Deletable,
	}
}

type RespExpeditionIndex struct {
	ID             uuid.UUID `json:"id"`
	ExpeditionCode string    `json:"expedition_code"`
	ExpeditionName string    `json:"expedition_name"`
	Address        string    `json:"address"`
	PhoneNumber    *string   `json:"phone_number"`
	TelpNumber     *string   `json:"telp_number"`
	CreatedAt      string    `json:"created_at"`
	UpdatedAt      string    `json:"updated_at"`
	Deletable      bool      `json:"deletable"`
}

func ToRespExpeditionIndex(m models.Expedition) RespExpeditionIndex {
	return RespExpeditionIndex{
		ID:             m.ID,
		ExpeditionCode: m.ExpeditionCode,
		ExpeditionName: m.ExpeditionName,
		Address:        m.Address,
		PhoneNumber:    m.PrimaryPhoneNumber,
		TelpNumber:     m.PrimaryTelpNumber,
		CreatedAt:      m.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:      m.UpdatedAt.Format("2006-01-02 15:04:05"),
		Deletable:      m.Deletable,
	}
}

// ExpeditionExport represents expedition data for export with all phone numbers
type ExpeditionExport struct {
	ExpeditionCode string
	ExpeditionName string
	Address        string
	PhoneNumbers   []string // All HP phone numbers
	TelpNumbers    []string // All Telp phone numbers
	UpdatedAt      time.Time
}

// ReqExpeditionIndexFilter for filtering expedition index (prepared for future use)
type ReqExpeditionIndexFilter struct {
	Search                 string   `query:"search" json:"search"` // Search keyword for filtering by expedition_name and expedition_code
	ExpeditionCodes        []string `query:"expedition_codes" json:"expedition_codes"`
	ExpeditionNames        []string `query:"expedition_names" json:"expedition_names"`
	ExpeditionCodesOrNames []string `query:"expedition_codes_or_names" json:"expedition_codes_or_names"`
	Addresses              []string `query:"addresses" json:"addresses"`
	TelpNumbers            []string `query:"telp_numbers" json:"telp_numbers"`
	PhoneNumbers           []string `query:"phone_numbers" json:"phone_numbers"`
	SortBy                 string   `query:"sort_by" json:"sort_by"`
	SortOrder              string   `query:"sort_order" json:"sort_order"`
}
