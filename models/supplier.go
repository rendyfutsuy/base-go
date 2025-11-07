package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Supplier represents suppliers table
type Supplier struct {
	ID uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v7()" json:"id" validate:"required"`

	SubdistrictID       *uuid.UUID `gorm:"column:subdistrict_id;type:uuid" json:"subdistrict_id"`
	DistrictID          *uuid.UUID `gorm:"column:district_id;type:uuid" json:"district_id"`
	CityID              *uuid.UUID `gorm:"column:city_id;type:uuid" json:"city_id"`
	ExpeditionArrivesID *uuid.UUID `gorm:"column:expedition_arrives_id;type:uuid" json:"expedition_arrives_id"`

	CreatedAt time.Time      `gorm:"column:created_at;not null" json:"created_at"`
	CreatedBy string         `gorm:"column:created_by;type:varchar(255)" json:"created_by"`
	UpdatedAt time.Time      `gorm:"column:updated_at;not null" json:"updated_at"`
	UpdatedBy string         `gorm:"column:updated_by;type:varchar(255)" json:"updated_by"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`
	DeletedBy *string        `gorm:"column:deleted_by;type:varchar(255)" json:"deleted_by"`

	SupplierCode string `gorm:"column:supplier_code;type:varchar(255);unique;not null" json:"supplier_code"`

	IdentityType     uuid.UUID `gorm:"column:identity_type;type:uuid;not null" json:"identity_type"`
	IdentityName     *string   `gorm:"column:identity_name;type:varchar(255)" json:"identity_name"`
	IdentityNumber   string    `gorm:"column:identity_number;type:varchar(255);unique;not null" json:"identity_number"`
	IdentityDocument *string   `gorm:"column:identity_document;type:text" json:"identity_document"`

	SupplierName string  `gorm:"column:supplier_name;type:varchar(255);not null" json:"supplier_name"`
	Alias        *string `gorm:"column:alias;type:varchar(255)" json:"alias"`
	RT           *string `gorm:"column:rt;type:varchar(255)" json:"rt"`
	RW           *string `gorm:"column:rw;type:varchar(255)" json:"rw"`
	PostalCode   *string `gorm:"column:postal_code;type:varchar(255)" json:"postal_code"`
	Address      string  `gorm:"column:address;type:text;not null" json:"address"`

	Email *string `gorm:"column:email;type:varchar(255)" json:"email"`
	Notes *string `gorm:"column:notes;type:text" json:"notes"`

	RelationDate          *time.Time `gorm:"column:relation_date;type:date;default:CURRENT_DATE" json:"relation_date"`
	DeliveryOption        uuid.UUID  `gorm:"column:delivery_option;type:uuid;not null" json:"delivery_option"`
	ExpeditionPaidBy      *uuid.UUID `gorm:"column:expedition_paid_by;type:uuid" json:"expedition_paid_by"`
	ExpeditionCalculation *uuid.UUID `gorm:"column:expedition_calculation;type:uuid" json:"expedition_calculation"`

	// Fetched mutators for primary contacts (from joins)
	PrimaryTelpNumber  *string `gorm:"column:primary_telp_number;<-:false" json:"primary_telp_number"`
	PrimaryPhoneNumber *string `gorm:"column:primary_phone_number;<-:false" json:"primary_phone_number"`
}

func (Supplier) TableName() string {
	return "suppliers"
}
