package test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/models"
	regencyDto "github.com/rendyfutsuy/base-go/modules/regency/dto"
	"github.com/rendyfutsuy/base-go/modules/regency/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRegencyRepository is a mock implementation of regency.Repository
type MockRegencyRepository struct {
	mock.Mock
}

// Province methods
func (m *MockRegencyRepository) CreateProvince(ctx context.Context, name string) (*models.Province, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Province), args.Error(1)
}

func (m *MockRegencyRepository) UpdateProvince(ctx context.Context, id uuid.UUID, name string) (*models.Province, error) {
	args := m.Called(ctx, id, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Province), args.Error(1)
}

func (m *MockRegencyRepository) DeleteProvince(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRegencyRepository) GetProvinceByID(ctx context.Context, id uuid.UUID) (*models.Province, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Province), args.Error(1)
}

func (m *MockRegencyRepository) GetProvinceIndex(ctx context.Context, req request.PageRequest, filter regencyDto.ReqProvinceIndexFilter) ([]models.Province, int, error) {
	args := m.Called(ctx, req, filter)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]models.Province), args.Int(1), args.Error(2)
}

func (m *MockRegencyRepository) GetAllProvince(ctx context.Context, filter regencyDto.ReqProvinceIndexFilter) ([]models.Province, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Province), args.Error(1)
}

func (m *MockRegencyRepository) ExistsProvinceByName(ctx context.Context, name string, excludeID uuid.UUID) (bool, error) {
	args := m.Called(ctx, name, excludeID)
	return args.Bool(0), args.Error(1)
}

// City methods
func (m *MockRegencyRepository) CreateCity(ctx context.Context, provinceID uuid.UUID, name string, areaCode *string) (*models.City, error) {
	args := m.Called(ctx, provinceID, name, areaCode)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.City), args.Error(1)
}

func (m *MockRegencyRepository) UpdateCity(ctx context.Context, id uuid.UUID, provinceID uuid.UUID, name string, areaCode *string) (*models.City, error) {
	args := m.Called(ctx, id, provinceID, name, areaCode)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.City), args.Error(1)
}

func (m *MockRegencyRepository) DeleteCity(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRegencyRepository) GetCityByID(ctx context.Context, id uuid.UUID) (*models.City, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.City), args.Error(1)
}

func (m *MockRegencyRepository) GetCityIndex(ctx context.Context, req request.PageRequest, filter regencyDto.ReqCityIndexFilter) ([]models.City, int, error) {
	args := m.Called(ctx, req, filter)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]models.City), args.Int(1), args.Error(2)
}

func (m *MockRegencyRepository) GetAllCity(ctx context.Context, filter regencyDto.ReqCityIndexFilter) ([]models.City, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.City), args.Error(1)
}

func (m *MockRegencyRepository) ExistsCityByName(ctx context.Context, provinceID uuid.UUID, name string, excludeID uuid.UUID) (bool, error) {
	args := m.Called(ctx, provinceID, name, excludeID)
	return args.Bool(0), args.Error(1)
}

func (m *MockRegencyRepository) GetCityAreaCodes(ctx context.Context, search string) ([]string, error) {
	args := m.Called(ctx, search)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

// District methods
func (m *MockRegencyRepository) CreateDistrict(ctx context.Context, cityID uuid.UUID, name string) (*models.District, error) {
	args := m.Called(ctx, cityID, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.District), args.Error(1)
}

func (m *MockRegencyRepository) UpdateDistrict(ctx context.Context, id uuid.UUID, cityID uuid.UUID, name string) (*models.District, error) {
	args := m.Called(ctx, id, cityID, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.District), args.Error(1)
}

func (m *MockRegencyRepository) DeleteDistrict(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRegencyRepository) GetDistrictByID(ctx context.Context, id uuid.UUID) (*models.District, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.District), args.Error(1)
}

func (m *MockRegencyRepository) GetDistrictIndex(ctx context.Context, req request.PageRequest, filter regencyDto.ReqDistrictIndexFilter) ([]models.District, int, error) {
	args := m.Called(ctx, req, filter)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]models.District), args.Int(1), args.Error(2)
}

func (m *MockRegencyRepository) GetAllDistrict(ctx context.Context, filter regencyDto.ReqDistrictIndexFilter) ([]models.District, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.District), args.Error(1)
}

func (m *MockRegencyRepository) ExistsDistrictByName(ctx context.Context, cityID uuid.UUID, name string, excludeID uuid.UUID) (bool, error) {
	args := m.Called(ctx, cityID, name, excludeID)
	return args.Bool(0), args.Error(1)
}

// Subdistrict methods
func (m *MockRegencyRepository) CreateSubdistrict(ctx context.Context, districtID uuid.UUID, name string) (*models.Subdistrict, error) {
	args := m.Called(ctx, districtID, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Subdistrict), args.Error(1)
}

func (m *MockRegencyRepository) UpdateSubdistrict(ctx context.Context, id uuid.UUID, districtID uuid.UUID, name string) (*models.Subdistrict, error) {
	args := m.Called(ctx, id, districtID, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Subdistrict), args.Error(1)
}

func (m *MockRegencyRepository) DeleteSubdistrict(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRegencyRepository) GetSubdistrictByID(ctx context.Context, id uuid.UUID) (*models.Subdistrict, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Subdistrict), args.Error(1)
}

func (m *MockRegencyRepository) GetSubdistrictIndex(ctx context.Context, req request.PageRequest, filter regencyDto.ReqSubdistrictIndexFilter) ([]models.Subdistrict, int, error) {
	args := m.Called(ctx, req, filter)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]models.Subdistrict), args.Int(1), args.Error(2)
}

func (m *MockRegencyRepository) GetAllSubdistrict(ctx context.Context, filter regencyDto.ReqSubdistrictIndexFilter) ([]models.Subdistrict, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Subdistrict), args.Error(1)
}

func (m *MockRegencyRepository) ExistsSubdistrictByName(ctx context.Context, districtID uuid.UUID, name string, excludeID uuid.UUID) (bool, error) {
	args := m.Called(ctx, districtID, name, excludeID)
	return args.Bool(0), args.Error(1)
}

// Province Tests
func TestCreateProvince(t *testing.T) {
	e := echo.New()
	ctx := context.Background()

	tests := []struct {
		name           string
		req            *regencyDto.ReqCreateProvince
		authId         string
		setupMock      func(*MockRegencyRepository)
		expectedError  error
		expectedResult *models.Province
	}{
		{
			name: "success create province",
			req: &regencyDto.ReqCreateProvince{
				Name: "Jawa Barat",
			},
			authId: "test-auth-id",
			setupMock: func(m *MockRegencyRepository) {
				m.On("ExistsProvinceByName", ctx, "Jawa Barat", uuid.Nil).Return(false, nil).Once()
				m.On("CreateProvince", ctx, "Jawa Barat").Return(&models.Province{
					ID:        uuid.New(),
					Name:      "Jawa Barat",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil).Once()
			},
			expectedError: nil,
			expectedResult: &models.Province{
				Name: "Jawa Barat",
			},
		},
		{
			name: "error when name already exists",
			req: &regencyDto.ReqCreateProvince{
				Name: "Jawa Barat",
			},
			authId: "test-auth-id",
			setupMock: func(m *MockRegencyRepository) {
				m.On("ExistsProvinceByName", ctx, "Jawa Barat", uuid.Nil).Return(true, nil).Once()
			},
			expectedError:  errors.New("Province name already exists"),
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRegencyRepository)
			tt.setupMock(mockRepo)

			usecaseInstance := usecase.NewRegencyUsecase(mockRepo)

			req := httptest.NewRequest(http.MethodPost, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetRequest(req.WithContext(ctx))

			result, err := usecaseInstance.CreateProvince(c, tt.req, tt.authId)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				if tt.expectedResult != nil {
					assert.Equal(t, tt.expectedResult.Name, result.Name)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUpdateProvince(t *testing.T) {
	e := echo.New()
	ctx := context.Background()
	validID := uuid.New()

	tests := []struct {
		name           string
		id             string
		req            *regencyDto.ReqUpdateProvince
		authId         string
		setupMock      func(*MockRegencyRepository)
		expectedError  error
		expectedResult *models.Province
	}{
		{
			name: "success update province",
			id:   validID.String(),
			req: &regencyDto.ReqUpdateProvince{
				Name: "Jawa Barat Updated",
			},
			authId: "test-auth-id",
			setupMock: func(m *MockRegencyRepository) {
				m.On("ExistsProvinceByName", ctx, "Jawa Barat Updated", validID).Return(false, nil).Once()
				m.On("UpdateProvince", ctx, validID, "Jawa Barat Updated").Return(&models.Province{
					ID:        validID,
					Name:      "Jawa Barat Updated",
					UpdatedAt: time.Now(),
				}, nil).Once()
			},
			expectedError: nil,
			expectedResult: &models.Province{
				ID:   validID,
				Name: "Jawa Barat Updated",
			},
		},
		{
			name: "error when invalid UUID",
			id:   "invalid-uuid",
			req: &regencyDto.ReqUpdateProvince{
				Name: "Test Province",
			},
			authId: "test-auth-id",
			setupMock: func(m *MockRegencyRepository) {
				// No mock calls expected
			},
			expectedError:  errors.New("requested param is string"),
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRegencyRepository)
			tt.setupMock(mockRepo)

			usecaseInstance := usecase.NewRegencyUsecase(mockRepo)

			req := httptest.NewRequest(http.MethodPut, "/:id", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.id)
			c.SetRequest(req.WithContext(ctx))

			result, err := usecaseInstance.UpdateProvince(c, tt.id, tt.req, tt.authId)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				if tt.expectedResult != nil {
					assert.Equal(t, tt.expectedResult.Name, result.Name)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestDeleteProvince(t *testing.T) {
	e := echo.New()
	ctx := context.Background()
	validID := uuid.New()

	tests := []struct {
		name          string
		id            string
		authId        string
		setupMock     func(*MockRegencyRepository)
		expectedError error
	}{
		{
			name:   "success delete province",
			id:     validID.String(),
			authId: "test-auth-id",
			setupMock: func(m *MockRegencyRepository) {
				m.On("DeleteProvince", ctx, validID).Return(nil).Once()
			},
			expectedError: nil,
		},
		{
			name:   "error when invalid UUID",
			id:     "invalid-uuid",
			authId: "test-auth-id",
			setupMock: func(m *MockRegencyRepository) {
				// No mock calls expected
			},
			expectedError: errors.New("requested param is string"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRegencyRepository)
			tt.setupMock(mockRepo)

			usecaseInstance := usecase.NewRegencyUsecase(mockRepo)

			req := httptest.NewRequest(http.MethodDelete, "/:id", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.id)
			c.SetRequest(req.WithContext(ctx))

			err := usecaseInstance.DeleteProvince(c, tt.id, tt.authId)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetProvinceByID(t *testing.T) {
	e := echo.New()
	ctx := context.Background()
	validID := uuid.New()

	tests := []struct {
		name           string
		id             string
		setupMock      func(*MockRegencyRepository)
		expectedError  error
		expectedResult *models.Province
	}{
		{
			name: "success get province by id",
			id:   validID.String(),
			setupMock: func(m *MockRegencyRepository) {
				m.On("GetProvinceByID", ctx, validID).Return(&models.Province{
					ID:        validID,
					Name:      "Jawa Barat",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil).Once()
			},
			expectedError: nil,
			expectedResult: &models.Province{
				ID:   validID,
				Name: "Jawa Barat",
			},
		},
		{
			name: "error when invalid UUID",
			id:   "invalid-uuid",
			setupMock: func(m *MockRegencyRepository) {
				// No mock calls expected
			},
			expectedError:  errors.New("requested param is string"),
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRegencyRepository)
			tt.setupMock(mockRepo)

			usecaseInstance := usecase.NewRegencyUsecase(mockRepo)

			req := httptest.NewRequest(http.MethodGet, "/:id", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.id)
			c.SetRequest(req.WithContext(ctx))

			result, err := usecaseInstance.GetProvinceByID(c, tt.id)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				if tt.expectedResult != nil {
					assert.Equal(t, tt.expectedResult.ID, result.ID)
					assert.Equal(t, tt.expectedResult.Name, result.Name)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetProvinceIndex(t *testing.T) {
	e := echo.New()
	ctx := context.Background()

	tests := []struct {
		name           string
		req            request.PageRequest
		filter         regencyDto.ReqProvinceIndexFilter
		setupMock      func(*MockRegencyRepository)
		expectedError  error
		expectedResult []models.Province
		expectedTotal  int
	}{
		{
			name: "success get index with pagination",
			req: request.PageRequest{
				Page:    1,
				PerPage: 10,
			},
			filter: regencyDto.ReqProvinceIndexFilter{},
			setupMock: func(m *MockRegencyRepository) {
				provinces := []models.Province{
					{
						ID:        uuid.New(),
						Name:      "Jawa Barat",
						CreatedAt: time.Now(),
					},
					{
						ID:        uuid.New(),
						Name:      "Jawa Tengah",
						CreatedAt: time.Now(),
					},
				}
				m.On("GetProvinceIndex", ctx, request.PageRequest{Page: 1, PerPage: 10}, regencyDto.ReqProvinceIndexFilter{}).Return(provinces, 2, nil).Once()
			},
			expectedError: nil,
			expectedResult: []models.Province{
				{Name: "Jawa Barat"},
				{Name: "Jawa Tengah"},
			},
			expectedTotal: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRegencyRepository)
			tt.setupMock(mockRepo)

			usecaseInstance := usecase.NewRegencyUsecase(mockRepo)

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetRequest(req.WithContext(ctx))

			result, total, err := usecaseInstance.GetProvinceIndex(c, tt.req, tt.filter)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
				assert.Equal(t, 0, total)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedTotal, total)
				assert.Equal(t, len(tt.expectedResult), len(result))
				if len(tt.expectedResult) > 0 {
					assert.Equal(t, tt.expectedResult[0].Name, result[0].Name)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestExportProvince(t *testing.T) {
	e := echo.New()
	ctx := context.Background()

	tests := []struct {
		name           string
		search         string
		filter         regencyDto.ReqProvinceIndexFilter
		setupMock      func(*MockRegencyRepository)
		expectedError  error
		expectedResult []models.Province
	}{
		{
			name:   "success export provinces",
			search: "",
			filter: regencyDto.ReqProvinceIndexFilter{},
			setupMock: func(m *MockRegencyRepository) {
				provinces := []models.Province{
					{
						ID:        uuid.New(),
						Name:      "Jawa Barat",
						CreatedAt: time.Now(),
					},
					{
						ID:        uuid.New(),
						Name:      "Jawa Tengah",
						CreatedAt: time.Now(),
					},
				}
				m.On("GetAllProvince", ctx, regencyDto.ReqProvinceIndexFilter{}).Return(provinces, nil).Once()
			},
			expectedError: nil,
			expectedResult: []models.Province{
				{Name: "Jawa Barat"},
				{Name: "Jawa Tengah"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRegencyRepository)
			tt.setupMock(mockRepo)

			usecaseInstance := usecase.NewRegencyUsecase(mockRepo)

			u := &url.URL{Path: "/export"}
			if tt.search != "" {
				q := u.Query()
				q.Set("search", tt.search)
				u.RawQuery = q.Encode()
			}
			req := httptest.NewRequest(http.MethodGet, u.String(), nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetRequest(req.WithContext(ctx))

			result, err := usecaseInstance.ExportProvince(c, tt.filter)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Greater(t, len(result), 0)
				assert.Equal(t, []byte("PK")[0], result[0])
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// City Tests
func TestCreateCity(t *testing.T) {
	e := echo.New()
	ctx := context.Background()
	provinceID := uuid.New()
	areaCode := "022"

	tests := []struct {
		name           string
		req            *regencyDto.ReqCreateCity
		authId         string
		setupMock      func(*MockRegencyRepository)
		expectedError  error
		expectedResult *models.City
	}{
		{
			name: "success create city",
			req: &regencyDto.ReqCreateCity{
				ProvinceID: provinceID,
				Name:       "Bandung",
				AreaCode:   &areaCode,
			},
			authId: "test-auth-id",
			setupMock: func(m *MockRegencyRepository) {
				m.On("GetProvinceByID", ctx, provinceID).Return(&models.Province{
					ID:   provinceID,
					Name: "Jawa Barat",
				}, nil).Once()
				m.On("ExistsCityByName", ctx, provinceID, "Bandung", uuid.Nil).Return(false, nil).Once()
				m.On("CreateCity", ctx, provinceID, "Bandung", &areaCode).Return(&models.City{
					ID:         uuid.New(),
					ProvinceID: provinceID,
					Name:       "Bandung",
					AreaCode:   &areaCode,
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
				}, nil).Once()
			},
			expectedError: nil,
			expectedResult: &models.City{
				ProvinceID: provinceID,
				Name:       "Bandung",
			},
		},
		{
			name: "error when name already exists in province",
			req: &regencyDto.ReqCreateCity{
				ProvinceID: provinceID,
				Name:       "Bandung",
				AreaCode:   &areaCode,
			},
			authId: "test-auth-id",
			setupMock: func(m *MockRegencyRepository) {
				m.On("GetProvinceByID", ctx, provinceID).Return(&models.Province{
					ID:   provinceID,
					Name: "Jawa Barat",
				}, nil).Once()
				m.On("ExistsCityByName", ctx, provinceID, "Bandung", uuid.Nil).Return(true, nil).Once()
			},
			expectedError:  errors.New("City name already exists in this province"),
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRegencyRepository)
			tt.setupMock(mockRepo)

			usecaseInstance := usecase.NewRegencyUsecase(mockRepo)

			req := httptest.NewRequest(http.MethodPost, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetRequest(req.WithContext(ctx))

			result, err := usecaseInstance.CreateCity(c, tt.req, tt.authId)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				if tt.expectedResult != nil {
					assert.Equal(t, tt.expectedResult.Name, result.Name)
					assert.Equal(t, tt.expectedResult.ProvinceID, result.ProvinceID)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUpdateCity(t *testing.T) {
	e := echo.New()
	ctx := context.Background()
	validID := uuid.New()
	provinceID := uuid.New()
	areaCode := "022"

	tests := []struct {
		name           string
		id             string
		req            *regencyDto.ReqUpdateCity
		authId         string
		setupMock      func(*MockRegencyRepository)
		expectedError  error
		expectedResult *models.City
	}{
		{
			name: "success update city",
			id:   validID.String(),
			req: &regencyDto.ReqUpdateCity{
				ProvinceID: provinceID,
				Name:       "Bandung Updated",
				AreaCode:   &areaCode,
			},
			authId: "test-auth-id",
			setupMock: func(m *MockRegencyRepository) {
				m.On("GetProvinceByID", ctx, provinceID).Return(&models.Province{
					ID:   provinceID,
					Name: "Jawa Barat",
				}, nil).Once()
				m.On("ExistsCityByName", ctx, provinceID, "Bandung Updated", validID).Return(false, nil).Once()
				m.On("UpdateCity", ctx, validID, provinceID, "Bandung Updated", &areaCode).Return(&models.City{
					ID:         validID,
					ProvinceID: provinceID,
					Name:       "Bandung Updated",
					AreaCode:   &areaCode,
					UpdatedAt:  time.Now(),
				}, nil).Once()
			},
			expectedError: nil,
			expectedResult: &models.City{
				ID:         validID,
				ProvinceID: provinceID,
				Name:       "Bandung Updated",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRegencyRepository)
			tt.setupMock(mockRepo)

			usecaseInstance := usecase.NewRegencyUsecase(mockRepo)

			req := httptest.NewRequest(http.MethodPut, "/:id", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.id)
			c.SetRequest(req.WithContext(ctx))

			result, err := usecaseInstance.UpdateCity(c, tt.id, tt.req, tt.authId)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				if tt.expectedResult != nil {
					assert.Equal(t, tt.expectedResult.Name, result.Name)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestDeleteCity(t *testing.T) {
	e := echo.New()
	ctx := context.Background()
	validID := uuid.New()

	tests := []struct {
		name          string
		id            string
		authId        string
		setupMock     func(*MockRegencyRepository)
		expectedError error
	}{
		{
			name:   "success delete city",
			id:     validID.String(),
			authId: "test-auth-id",
			setupMock: func(m *MockRegencyRepository) {
				m.On("DeleteCity", ctx, validID).Return(nil).Once()
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRegencyRepository)
			tt.setupMock(mockRepo)

			usecaseInstance := usecase.NewRegencyUsecase(mockRepo)

			req := httptest.NewRequest(http.MethodDelete, "/:id", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.id)
			c.SetRequest(req.WithContext(ctx))

			err := usecaseInstance.DeleteCity(c, tt.id, tt.authId)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetCityByID(t *testing.T) {
	e := echo.New()
	ctx := context.Background()
	validID := uuid.New()
	provinceID := uuid.New()

	tests := []struct {
		name           string
		id             string
		setupMock      func(*MockRegencyRepository)
		expectedError  error
		expectedResult *models.City
	}{
		{
			name: "success get city by id",
			id:   validID.String(),
			setupMock: func(m *MockRegencyRepository) {
				m.On("GetCityByID", ctx, validID).Return(&models.City{
					ID:         validID,
					ProvinceID: provinceID,
					Name:       "Bandung",
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
				}, nil).Once()
			},
			expectedError: nil,
			expectedResult: &models.City{
				ID:         validID,
				ProvinceID: provinceID,
				Name:       "Bandung",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRegencyRepository)
			tt.setupMock(mockRepo)

			usecaseInstance := usecase.NewRegencyUsecase(mockRepo)

			req := httptest.NewRequest(http.MethodGet, "/:id", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.id)
			c.SetRequest(req.WithContext(ctx))

			result, err := usecaseInstance.GetCityByID(c, tt.id)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				if tt.expectedResult != nil {
					assert.Equal(t, tt.expectedResult.ID, result.ID)
					assert.Equal(t, tt.expectedResult.Name, result.Name)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetCityIndex(t *testing.T) {
	e := echo.New()
	ctx := context.Background()

	tests := []struct {
		name           string
		req            request.PageRequest
		filter         regencyDto.ReqCityIndexFilter
		setupMock      func(*MockRegencyRepository)
		expectedError  error
		expectedResult []models.City
		expectedTotal  int
	}{
		{
			name: "success get index with pagination",
			req: request.PageRequest{
				Page:    1,
				PerPage: 10,
			},
			filter: regencyDto.ReqCityIndexFilter{},
			setupMock: func(m *MockRegencyRepository) {
				cities := []models.City{
					{
						ID:        uuid.New(),
						Name:      "Bandung",
						CreatedAt: time.Now(),
					},
				}
				m.On("GetCityIndex", ctx, request.PageRequest{Page: 1, PerPage: 10}, regencyDto.ReqCityIndexFilter{}).Return(cities, 1, nil).Once()
			},
			expectedError: nil,
			expectedResult: []models.City{
				{Name: "Bandung"},
			},
			expectedTotal: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRegencyRepository)
			tt.setupMock(mockRepo)

			usecaseInstance := usecase.NewRegencyUsecase(mockRepo)

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetRequest(req.WithContext(ctx))

			result, total, err := usecaseInstance.GetCityIndex(c, tt.req, tt.filter)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
				assert.Equal(t, 0, total)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedTotal, total)
				assert.Equal(t, len(tt.expectedResult), len(result))
				if len(tt.expectedResult) > 0 {
					assert.Equal(t, tt.expectedResult[0].Name, result[0].Name)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestExportCity(t *testing.T) {
	e := echo.New()
	ctx := context.Background()

	tests := []struct {
		name           string
		search         string
		filter         regencyDto.ReqCityIndexFilter
		setupMock      func(*MockRegencyRepository)
		expectedError  error
		expectedResult []models.City
	}{
		{
			name:   "success export cities",
			search: "",
			filter: regencyDto.ReqCityIndexFilter{},
			setupMock: func(m *MockRegencyRepository) {
				cities := []models.City{
					{
						ID:        uuid.New(),
						Name:      "Bandung",
						CreatedAt: time.Now(),
					},
				}
				m.On("GetAllCity", ctx, regencyDto.ReqCityIndexFilter{}).Return(cities, nil).Once()
			},
			expectedError: nil,
			expectedResult: []models.City{
				{Name: "Bandung"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRegencyRepository)
			tt.setupMock(mockRepo)

			usecaseInstance := usecase.NewRegencyUsecase(mockRepo)

			u := &url.URL{Path: "/export"}
			if tt.search != "" {
				q := u.Query()
				q.Set("search", tt.search)
				u.RawQuery = q.Encode()
			}
			req := httptest.NewRequest(http.MethodGet, u.String(), nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetRequest(req.WithContext(ctx))

			result, err := usecaseInstance.ExportCity(c, tt.filter)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Greater(t, len(result), 0)
				assert.Equal(t, []byte("PK")[0], result[0])
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// District Tests
func TestCreateDistrict(t *testing.T) {
	e := echo.New()
	ctx := context.Background()
	cityID := uuid.New()

	tests := []struct {
		name           string
		req            *regencyDto.ReqCreateDistrict
		authId         string
		setupMock      func(*MockRegencyRepository)
		expectedError  error
		expectedResult *models.District
	}{
		{
			name: "success create district",
			req: &regencyDto.ReqCreateDistrict{
				CityID: cityID,
				Name:   "Cicendo",
			},
			authId: "test-auth-id",
			setupMock: func(m *MockRegencyRepository) {
				m.On("GetCityByID", ctx, cityID).Return(&models.City{
					ID:   cityID,
					Name: "Bandung",
				}, nil).Once()
				m.On("ExistsDistrictByName", ctx, cityID, "Cicendo", uuid.Nil).Return(false, nil).Once()
				m.On("CreateDistrict", ctx, cityID, "Cicendo").Return(&models.District{
					ID:        uuid.New(),
					CityID:    cityID,
					Name:      "Cicendo",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil).Once()
			},
			expectedError: nil,
			expectedResult: &models.District{
				CityID: cityID,
				Name:   "Cicendo",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRegencyRepository)
			tt.setupMock(mockRepo)

			usecaseInstance := usecase.NewRegencyUsecase(mockRepo)

			req := httptest.NewRequest(http.MethodPost, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetRequest(req.WithContext(ctx))

			result, err := usecaseInstance.CreateDistrict(c, tt.req, tt.authId)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				if tt.expectedResult != nil {
					assert.Equal(t, tt.expectedResult.Name, result.Name)
					assert.Equal(t, tt.expectedResult.CityID, result.CityID)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUpdateDistrict(t *testing.T) {
	e := echo.New()
	ctx := context.Background()
	validID := uuid.New()
	cityID := uuid.New()

	tests := []struct {
		name           string
		id             string
		req            *regencyDto.ReqUpdateDistrict
		authId         string
		setupMock      func(*MockRegencyRepository)
		expectedError  error
		expectedResult *models.District
	}{
		{
			name: "success update district",
			id:   validID.String(),
			req: &regencyDto.ReqUpdateDistrict{
				CityID: cityID,
				Name:   "Cicendo Updated",
			},
			authId: "test-auth-id",
			setupMock: func(m *MockRegencyRepository) {
				m.On("GetCityByID", ctx, cityID).Return(&models.City{
					ID:   cityID,
					Name: "Bandung",
				}, nil).Once()
				m.On("ExistsDistrictByName", ctx, cityID, "Cicendo Updated", validID).Return(false, nil).Once()
				m.On("UpdateDistrict", ctx, validID, cityID, "Cicendo Updated").Return(&models.District{
					ID:        validID,
					CityID:    cityID,
					Name:      "Cicendo Updated",
					UpdatedAt: time.Now(),
				}, nil).Once()
			},
			expectedError: nil,
			expectedResult: &models.District{
				ID:     validID,
				CityID: cityID,
				Name:   "Cicendo Updated",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRegencyRepository)
			tt.setupMock(mockRepo)

			usecaseInstance := usecase.NewRegencyUsecase(mockRepo)

			req := httptest.NewRequest(http.MethodPut, "/:id", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.id)
			c.SetRequest(req.WithContext(ctx))

			result, err := usecaseInstance.UpdateDistrict(c, tt.id, tt.req, tt.authId)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				if tt.expectedResult != nil {
					assert.Equal(t, tt.expectedResult.Name, result.Name)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestDeleteDistrict(t *testing.T) {
	e := echo.New()
	ctx := context.Background()
	validID := uuid.New()

	tests := []struct {
		name          string
		id            string
		authId        string
		setupMock     func(*MockRegencyRepository)
		expectedError error
	}{
		{
			name:   "success delete district",
			id:     validID.String(),
			authId: "test-auth-id",
			setupMock: func(m *MockRegencyRepository) {
				m.On("DeleteDistrict", ctx, validID).Return(nil).Once()
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRegencyRepository)
			tt.setupMock(mockRepo)

			usecaseInstance := usecase.NewRegencyUsecase(mockRepo)

			req := httptest.NewRequest(http.MethodDelete, "/:id", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.id)
			c.SetRequest(req.WithContext(ctx))

			err := usecaseInstance.DeleteDistrict(c, tt.id, tt.authId)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetDistrictByID(t *testing.T) {
	e := echo.New()
	ctx := context.Background()
	validID := uuid.New()
	cityID := uuid.New()

	tests := []struct {
		name           string
		id             string
		setupMock      func(*MockRegencyRepository)
		expectedError  error
		expectedResult *models.District
	}{
		{
			name: "success get district by id",
			id:   validID.String(),
			setupMock: func(m *MockRegencyRepository) {
				m.On("GetDistrictByID", ctx, validID).Return(&models.District{
					ID:        validID,
					CityID:    cityID,
					Name:      "Cicendo",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil).Once()
			},
			expectedError: nil,
			expectedResult: &models.District{
				ID:     validID,
				CityID: cityID,
				Name:   "Cicendo",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRegencyRepository)
			tt.setupMock(mockRepo)

			usecaseInstance := usecase.NewRegencyUsecase(mockRepo)

			req := httptest.NewRequest(http.MethodGet, "/:id", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.id)
			c.SetRequest(req.WithContext(ctx))

			result, err := usecaseInstance.GetDistrictByID(c, tt.id)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				if tt.expectedResult != nil {
					assert.Equal(t, tt.expectedResult.ID, result.ID)
					assert.Equal(t, tt.expectedResult.Name, result.Name)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetDistrictIndex(t *testing.T) {
	e := echo.New()
	ctx := context.Background()

	tests := []struct {
		name           string
		req            request.PageRequest
		filter         regencyDto.ReqDistrictIndexFilter
		setupMock      func(*MockRegencyRepository)
		expectedError  error
		expectedResult []models.District
		expectedTotal  int
	}{
		{
			name: "success get index with pagination",
			req: request.PageRequest{
				Page:    1,
				PerPage: 10,
			},
			filter: regencyDto.ReqDistrictIndexFilter{},
			setupMock: func(m *MockRegencyRepository) {
				districts := []models.District{
					{
						ID:        uuid.New(),
						Name:      "Cicendo",
						CreatedAt: time.Now(),
					},
				}
				m.On("GetDistrictIndex", ctx, request.PageRequest{Page: 1, PerPage: 10}, regencyDto.ReqDistrictIndexFilter{}).Return(districts, 1, nil).Once()
			},
			expectedError: nil,
			expectedResult: []models.District{
				{Name: "Cicendo"},
			},
			expectedTotal: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRegencyRepository)
			tt.setupMock(mockRepo)

			usecaseInstance := usecase.NewRegencyUsecase(mockRepo)

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetRequest(req.WithContext(ctx))

			result, total, err := usecaseInstance.GetDistrictIndex(c, tt.req, tt.filter)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
				assert.Equal(t, 0, total)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedTotal, total)
				assert.Equal(t, len(tt.expectedResult), len(result))
				if len(tt.expectedResult) > 0 {
					assert.Equal(t, tt.expectedResult[0].Name, result[0].Name)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestExportDistrict(t *testing.T) {
	e := echo.New()
	ctx := context.Background()

	tests := []struct {
		name           string
		search         string
		filter         regencyDto.ReqDistrictIndexFilter
		setupMock      func(*MockRegencyRepository)
		expectedError  error
		expectedResult []models.District
	}{
		{
			name:   "success export districts",
			search: "",
			filter: regencyDto.ReqDistrictIndexFilter{},
			setupMock: func(m *MockRegencyRepository) {
				districts := []models.District{
					{
						ID:        uuid.New(),
						Name:      "Cicendo",
						CreatedAt: time.Now(),
					},
				}
				m.On("GetAllDistrict", ctx, regencyDto.ReqDistrictIndexFilter{}).Return(districts, nil).Once()
			},
			expectedError: nil,
			expectedResult: []models.District{
				{Name: "Cicendo"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRegencyRepository)
			tt.setupMock(mockRepo)

			usecaseInstance := usecase.NewRegencyUsecase(mockRepo)

			u := &url.URL{Path: "/export"}
			if tt.search != "" {
				q := u.Query()
				q.Set("search", tt.search)
				u.RawQuery = q.Encode()
			}
			req := httptest.NewRequest(http.MethodGet, u.String(), nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetRequest(req.WithContext(ctx))

			result, err := usecaseInstance.ExportDistrict(c, tt.filter)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Greater(t, len(result), 0)
				assert.Equal(t, []byte("PK")[0], result[0])
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// Subdistrict Tests
func TestCreateSubdistrict(t *testing.T) {
	e := echo.New()
	ctx := context.Background()
	districtID := uuid.New()

	tests := []struct {
		name           string
		req            *regencyDto.ReqCreateSubdistrict
		authId         string
		setupMock      func(*MockRegencyRepository)
		expectedError  error
		expectedResult *models.Subdistrict
	}{
		{
			name: "success create subdistrict",
			req: &regencyDto.ReqCreateSubdistrict{
				DistrictID: districtID,
				Name:       "Arjuna",
			},
			authId: "test-auth-id",
			setupMock: func(m *MockRegencyRepository) {
				m.On("GetDistrictByID", ctx, districtID).Return(&models.District{
					ID:   districtID,
					Name: "Cicendo",
				}, nil).Once()
				m.On("ExistsSubdistrictByName", ctx, districtID, "Arjuna", uuid.Nil).Return(false, nil).Once()
				m.On("CreateSubdistrict", ctx, districtID, "Arjuna").Return(&models.Subdistrict{
					ID:         uuid.New(),
					DistrictID: districtID,
					Name:       "Arjuna",
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
				}, nil).Once()
			},
			expectedError: nil,
			expectedResult: &models.Subdistrict{
				DistrictID: districtID,
				Name:       "Arjuna",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRegencyRepository)
			tt.setupMock(mockRepo)

			usecaseInstance := usecase.NewRegencyUsecase(mockRepo)

			req := httptest.NewRequest(http.MethodPost, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetRequest(req.WithContext(ctx))

			result, err := usecaseInstance.CreateSubdistrict(c, tt.req, tt.authId)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				if tt.expectedResult != nil {
					assert.Equal(t, tt.expectedResult.Name, result.Name)
					assert.Equal(t, tt.expectedResult.DistrictID, result.DistrictID)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUpdateSubdistrict(t *testing.T) {
	e := echo.New()
	ctx := context.Background()
	validID := uuid.New()
	districtID := uuid.New()

	tests := []struct {
		name           string
		id             string
		req            *regencyDto.ReqUpdateSubdistrict
		authId         string
		setupMock      func(*MockRegencyRepository)
		expectedError  error
		expectedResult *models.Subdistrict
	}{
		{
			name: "success update subdistrict",
			id:   validID.String(),
			req: &regencyDto.ReqUpdateSubdistrict{
				DistrictID: districtID,
				Name:       "Arjuna Updated",
			},
			authId: "test-auth-id",
			setupMock: func(m *MockRegencyRepository) {
				m.On("GetDistrictByID", ctx, districtID).Return(&models.District{
					ID:   districtID,
					Name: "Cicendo",
				}, nil).Once()
				m.On("ExistsSubdistrictByName", ctx, districtID, "Arjuna Updated", validID).Return(false, nil).Once()
				m.On("UpdateSubdistrict", ctx, validID, districtID, "Arjuna Updated").Return(&models.Subdistrict{
					ID:         validID,
					DistrictID: districtID,
					Name:       "Arjuna Updated",
					UpdatedAt:  time.Now(),
				}, nil).Once()
			},
			expectedError: nil,
			expectedResult: &models.Subdistrict{
				ID:         validID,
				DistrictID: districtID,
				Name:       "Arjuna Updated",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRegencyRepository)
			tt.setupMock(mockRepo)

			usecaseInstance := usecase.NewRegencyUsecase(mockRepo)

			req := httptest.NewRequest(http.MethodPut, "/:id", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.id)
			c.SetRequest(req.WithContext(ctx))

			result, err := usecaseInstance.UpdateSubdistrict(c, tt.id, tt.req, tt.authId)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				if tt.expectedResult != nil {
					assert.Equal(t, tt.expectedResult.Name, result.Name)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestDeleteSubdistrict(t *testing.T) {
	e := echo.New()
	ctx := context.Background()
	validID := uuid.New()

	tests := []struct {
		name          string
		id            string
		authId        string
		setupMock     func(*MockRegencyRepository)
		expectedError error
	}{
		{
			name:   "success delete subdistrict",
			id:     validID.String(),
			authId: "test-auth-id",
			setupMock: func(m *MockRegencyRepository) {
				m.On("DeleteSubdistrict", ctx, validID).Return(nil).Once()
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRegencyRepository)
			tt.setupMock(mockRepo)

			usecaseInstance := usecase.NewRegencyUsecase(mockRepo)

			req := httptest.NewRequest(http.MethodDelete, "/:id", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.id)
			c.SetRequest(req.WithContext(ctx))

			err := usecaseInstance.DeleteSubdistrict(c, tt.id, tt.authId)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetSubdistrictByID(t *testing.T) {
	e := echo.New()
	ctx := context.Background()
	validID := uuid.New()
	districtID := uuid.New()

	tests := []struct {
		name           string
		id             string
		setupMock      func(*MockRegencyRepository)
		expectedError  error
		expectedResult *models.Subdistrict
	}{
		{
			name: "success get subdistrict by id",
			id:   validID.String(),
			setupMock: func(m *MockRegencyRepository) {
				m.On("GetSubdistrictByID", ctx, validID).Return(&models.Subdistrict{
					ID:         validID,
					DistrictID: districtID,
					Name:       "Arjuna",
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
				}, nil).Once()
			},
			expectedError: nil,
			expectedResult: &models.Subdistrict{
				ID:         validID,
				DistrictID: districtID,
				Name:       "Arjuna",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRegencyRepository)
			tt.setupMock(mockRepo)

			usecaseInstance := usecase.NewRegencyUsecase(mockRepo)

			req := httptest.NewRequest(http.MethodGet, "/:id", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.id)
			c.SetRequest(req.WithContext(ctx))

			result, err := usecaseInstance.GetSubdistrictByID(c, tt.id)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				if tt.expectedResult != nil {
					assert.Equal(t, tt.expectedResult.ID, result.ID)
					assert.Equal(t, tt.expectedResult.Name, result.Name)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetSubdistrictIndex(t *testing.T) {
	e := echo.New()
	ctx := context.Background()

	tests := []struct {
		name           string
		req            request.PageRequest
		filter         regencyDto.ReqSubdistrictIndexFilter
		setupMock      func(*MockRegencyRepository)
		expectedError  error
		expectedResult []models.Subdistrict
		expectedTotal  int
	}{
		{
			name: "success get index with pagination",
			req: request.PageRequest{
				Page:    1,
				PerPage: 10,
			},
			filter: regencyDto.ReqSubdistrictIndexFilter{},
			setupMock: func(m *MockRegencyRepository) {
				subdistricts := []models.Subdistrict{
					{
						ID:        uuid.New(),
						Name:      "Arjuna",
						CreatedAt: time.Now(),
					},
				}
				m.On("GetSubdistrictIndex", ctx, request.PageRequest{Page: 1, PerPage: 10}, regencyDto.ReqSubdistrictIndexFilter{}).Return(subdistricts, 1, nil).Once()
			},
			expectedError: nil,
			expectedResult: []models.Subdistrict{
				{Name: "Arjuna"},
			},
			expectedTotal: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRegencyRepository)
			tt.setupMock(mockRepo)

			usecaseInstance := usecase.NewRegencyUsecase(mockRepo)

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetRequest(req.WithContext(ctx))

			result, total, err := usecaseInstance.GetSubdistrictIndex(c, tt.req, tt.filter)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
				assert.Equal(t, 0, total)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedTotal, total)
				assert.Equal(t, len(tt.expectedResult), len(result))
				if len(tt.expectedResult) > 0 {
					assert.Equal(t, tt.expectedResult[0].Name, result[0].Name)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestExportSubdistrict(t *testing.T) {
	e := echo.New()
	ctx := context.Background()

	tests := []struct {
		name           string
		search         string
		filter         regencyDto.ReqSubdistrictIndexFilter
		setupMock      func(*MockRegencyRepository)
		expectedError  error
		expectedResult []models.Subdistrict
	}{
		{
			name:   "success export subdistricts",
			search: "",
			filter: regencyDto.ReqSubdistrictIndexFilter{},
			setupMock: func(m *MockRegencyRepository) {
				subdistricts := []models.Subdistrict{
					{
						ID:        uuid.New(),
						Name:      "Arjuna",
						CreatedAt: time.Now(),
					},
				}
				m.On("GetAllSubdistrict", ctx, regencyDto.ReqSubdistrictIndexFilter{}).Return(subdistricts, nil).Once()
			},
			expectedError: nil,
			expectedResult: []models.Subdistrict{
				{Name: "Arjuna"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRegencyRepository)
			tt.setupMock(mockRepo)

			usecaseInstance := usecase.NewRegencyUsecase(mockRepo)

			u := &url.URL{Path: "/export"}
			if tt.search != "" {
				q := u.Query()
				q.Set("search", tt.search)
				u.RawQuery = q.Encode()
			}
			req := httptest.NewRequest(http.MethodGet, u.String(), nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetRequest(req.WithContext(ctx))

			result, err := usecaseInstance.ExportSubdistrict(c, tt.filter)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Greater(t, len(result), 0)
				assert.Equal(t, []byte("PK")[0], result[0])
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
