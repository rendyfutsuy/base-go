package constants

const (
	// Province validation errors
	ProvinceNameAlreadyExists    = "Province name already exists"
	ProvinceCreateFailedIDNotSet = "failed to create province: ID not set"
	ProvinceNotFound             = "province with id %s not found"

	// City validation errors
	CityNameAlreadyExists    = "City name already exists in this province"
	CityCreateFailedIDNotSet = "failed to create city: ID not set"
	CityProvinceNotFound     = "Province not found"
	CityNotFound             = "city with id %s not found"

	// District validation errors
	DistrictNameAlreadyExists    = "District name already exists in this city"
	DistrictCreateFailedIDNotSet = "failed to create district: ID not set"
	DistrictCityNotFound         = "City not found"
	DistrictNotFound             = "district with id %s not found"

	// Subdistrict validation errors
	SubdistrictNameAlreadyExists    = "Subdistrict name already exists in this district"
	SubdistrictCreateFailedIDNotSet = "failed to create subdistrict: ID not set"
	SubdistrictDistrictNotFound     = "District not found"
	SubdistrictNotFound             = "subdistrict with id %s not found"

	// Success messages
	ProvinceDeleteSuccess    = "Successfully deleted Province"
	CityDeleteSuccess        = "Successfully deleted City"
	DistrictDeleteSuccess    = "Successfully deleted District"
	SubdistrictDeleteSuccess = "Successfully deleted Subdistrict"
)
