package dto

import (
	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/models"
)

// Province DTOs
type ReqCreateProvince struct {
	Name string `form:"name" json:"name" validate:"required,max=100"`
}

type ReqUpdateProvince struct {
	Name string `form:"name" json:"name" validate:"required,max=100"`
}

type RespProvince struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

func ToRespProvince(m models.Province) RespProvince {
	return RespProvince{
		ID:        m.ID,
		Name:      m.Name,
		CreatedAt: m.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: m.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

type RespProvinceIndex struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

func ToRespProvinceIndex(m models.Province) RespProvinceIndex {
	return RespProvinceIndex{
		ID:        m.ID,
		Name:      m.Name,
		CreatedAt: m.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: m.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

type ReqProvinceIndexFilter struct {
	Search string `query:"search" json:"search"` // Search keyword for filtering by name
}

// City DTOs
type ReqCreateCity struct {
	ProvinceID uuid.UUID `form:"province_id" json:"province_id" validate:"required"`
	Name       string    `form:"name" json:"name" validate:"required,max=255"`
}

type ReqUpdateCity struct {
	ProvinceID uuid.UUID `form:"province_id" json:"province_id" validate:"required"`
	Name       string    `form:"name" json:"name" validate:"required,max=255"`
}

type RespCity struct {
	ID         uuid.UUID     `json:"id"`
	ProvinceID uuid.UUID     `json:"province_id"`
	Province   *RespProvince `json:"province,omitempty"`
	Name       string        `json:"name"`
	CreatedAt  string        `json:"created_at"`
	UpdatedAt  string        `json:"updated_at"`
}

func ToRespCity(m models.City) RespCity {
	resp := RespCity{
		ID:         m.ID,
		ProvinceID: m.ProvinceID,
		Name:       m.Name,
		CreatedAt:  m.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:  m.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	if m.Province.ID != uuid.Nil {
		province := ToRespProvince(m.Province)
		resp.Province = &province
	}
	return resp
}

type RespCityIndex struct {
	ID         uuid.UUID `json:"id"`
	ProvinceID uuid.UUID `json:"province_id"`
	Name       string    `json:"name"`
	CreatedAt  string    `json:"created_at"`
	UpdatedAt  string    `json:"updated_at"`
}

func ToRespCityIndex(m models.City) RespCityIndex {
	return RespCityIndex{
		ID:         m.ID,
		ProvinceID: m.ProvinceID,
		Name:       m.Name,
		CreatedAt:  m.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:  m.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

type ReqCityIndexFilter struct {
	Search     string    `query:"search" json:"search"`           // Search keyword for filtering by name
	ProvinceID uuid.UUID `query:"province_id" json:"province_id"` // Filter by province_id
}

// District DTOs
type ReqCreateDistrict struct {
	CityID uuid.UUID `form:"city_id" json:"city_id" validate:"required"`
	Name   string    `form:"name" json:"name" validate:"required,max=255"`
}

type ReqUpdateDistrict struct {
	CityID uuid.UUID `form:"city_id" json:"city_id" validate:"required"`
	Name   string    `form:"name" json:"name" validate:"required,max=255"`
}

type RespDistrict struct {
	ID        uuid.UUID `json:"id"`
	CityID    uuid.UUID `json:"city_id"`
	City      *RespCity `json:"city,omitempty"`
	Name      string    `json:"name"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

func ToRespDistrict(m models.District) RespDistrict {
	resp := RespDistrict{
		ID:        m.ID,
		CityID:    m.CityID,
		Name:      m.Name,
		CreatedAt: m.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: m.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	if m.City.ID != uuid.Nil {
		city := ToRespCity(m.City)
		resp.City = &city
	}
	return resp
}

type RespDistrictIndex struct {
	ID        uuid.UUID `json:"id"`
	CityID    uuid.UUID `json:"city_id"`
	Name      string    `json:"name"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

func ToRespDistrictIndex(m models.District) RespDistrictIndex {
	return RespDistrictIndex{
		ID:        m.ID,
		CityID:    m.CityID,
		Name:      m.Name,
		CreatedAt: m.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: m.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

type ReqDistrictIndexFilter struct {
	Search string    `query:"search" json:"search"`   // Search keyword for filtering by name
	CityID uuid.UUID `query:"city_id" json:"city_id"` // Filter by city_id
}

// Subdistrict DTOs
type ReqCreateSubdistrict struct {
	DistrictID uuid.UUID `form:"district_id" json:"district_id" validate:"required"`
	Name       string    `form:"name" json:"name" validate:"required,max=255"`
}

type ReqUpdateSubdistrict struct {
	DistrictID uuid.UUID `form:"district_id" json:"district_id" validate:"required"`
	Name       string    `form:"name" json:"name" validate:"required,max=255"`
}

type RespSubdistrict struct {
	ID         uuid.UUID     `json:"id"`
	DistrictID uuid.UUID     `json:"district_id"`
	District   *RespDistrict `json:"district,omitempty"`
	Name       string        `json:"name"`
	CreatedAt  string        `json:"created_at"`
	UpdatedAt  string        `json:"updated_at"`
}

func ToRespSubdistrict(m models.Subdistrict) RespSubdistrict {
	resp := RespSubdistrict{
		ID:         m.ID,
		DistrictID: m.DistrictID,
		Name:       m.Name,
		CreatedAt:  m.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:  m.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	if m.District.ID != uuid.Nil {
		district := ToRespDistrict(m.District)
		resp.District = &district
	}
	return resp
}

type RespSubdistrictIndex struct {
	ID         uuid.UUID `json:"id"`
	DistrictID uuid.UUID `json:"district_id"`
	Name       string    `json:"name"`
	CreatedAt  string    `json:"created_at"`
	UpdatedAt  string    `json:"updated_at"`
}

func ToRespSubdistrictIndex(m models.Subdistrict) RespSubdistrictIndex {
	return RespSubdistrictIndex{
		ID:         m.ID,
		DistrictID: m.DistrictID,
		Name:       m.Name,
		CreatedAt:  m.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:  m.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

type ReqSubdistrictIndexFilter struct {
	Search     string    `query:"search" json:"search"`           // Search keyword for filtering by name
	DistrictID uuid.UUID `query:"district_id" json:"district_id"` // Filter by district_id
}
