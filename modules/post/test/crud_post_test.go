package test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/models"
	paramDto "github.com/rendyfutsuy/base-go/modules/parameter/dto"
	postDto "github.com/rendyfutsuy/base-go/modules/post/dto"
	postUsecase "github.com/rendyfutsuy/base-go/modules/post/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPostRepository struct {
	mock.Mock
}

func (m *MockPostRepository) Create(ctx context.Context, createdBy uuid.UUID, data postDto.ToDBPost) (*models.Post, error) {
	args := m.Called(ctx, createdBy, data)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Post), args.Error(1)
}
func (m *MockPostRepository) Update(ctx context.Context, id uuid.UUID, data postDto.ToDBPost) (*models.Post, error) {
	args := m.Called(ctx, id, data)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Post), args.Error(1)
}
func (m *MockPostRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *MockPostRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Post, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Post), args.Error(1)
}
func (m *MockPostRepository) GetIndex(ctx context.Context, req request.PageRequest, filter postDto.ReqPostIndexFilter) ([]models.Post, int, error) {
	args := m.Called(ctx, req, filter)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]models.Post), args.Int(1), args.Error(2)
}
func (m *MockPostRepository) GetAll(ctx context.Context, filter postDto.ReqPostIndexFilter) ([]models.Post, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Post), args.Error(1)
}

type MockParameterRepository struct {
	mock.Mock
}

func (m *MockParameterRepository) Create(ctx context.Context, code, name string, value, typeVal, desc *string) (*models.Parameter, error) {
	panic("not implemented")
}
func (m *MockParameterRepository) Update(ctx context.Context, id uuid.UUID, code, name string, value, typeVal, desc *string) (*models.Parameter, error) {
	panic("not implemented")
}
func (m *MockParameterRepository) SetParent(ctx context.Context, id uuid.UUID, parentID uuid.UUID) error {
	panic("not implemented")
}
func (m *MockParameterRepository) Delete(ctx context.Context, id uuid.UUID) error {
	panic("not implemented")
}
func (m *MockParameterRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Parameter, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Parameter), args.Error(1)
}
func (m *MockParameterRepository) GetIndex(ctx context.Context, req request.PageRequest, filter paramDto.ReqParameterIndexFilter) ([]models.Parameter, int, error) {
	panic("not implemented")
}
func (m *MockParameterRepository) GetAll(ctx context.Context, filter paramDto.ReqParameterIndexFilter) ([]models.Parameter, error) {
	panic("not implemented")
}
func (m *MockParameterRepository) ExistsByCode(ctx context.Context, code string, excludeID uuid.UUID) (bool, error) {
	panic("not implemented")
}
func (m *MockParameterRepository) ExistsByName(ctx context.Context, name string, excludeID uuid.UUID) (bool, error) {
	panic("not implemented")
}
func (m *MockParameterRepository) AssignParametersToModule(ctx context.Context, moduleType string, moduleID uuid.UUID, parameterIDs []uuid.UUID) error {
	args := m.Called(ctx, moduleType, moduleID, parameterIDs)
	return args.Error(0)
}
func (m *MockParameterRepository) GetByModule(ctx context.Context, moduleType string, moduleID uuid.UUID) ([]models.Parameter, error) {
	args := m.Called(ctx, moduleType, moduleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Parameter), args.Error(1)
}
func (m *MockParameterRepository) ClearParametersForModule(ctx context.Context, moduleType string, moduleID uuid.UUID) error {
	args := m.Called(ctx, moduleType, moduleID)
	return args.Error(0)
}

func TestPostUsecase_Create_TypeValidation(t *testing.T) {
	ctx := context.Background()
	mockPostRepo := new(MockPostRepository)
	mockParamRepo := new(MockParameterRepository)
	useCase := postUsecase.NewPostUsecase(mockPostRepo, mockParamRepo)

	langID := uuid.New()
	topicID := uuid.New()

	// invalid level (not post_level)
	req := &postDto.ReqCreatePost{
		Title:            "A",
		Description:      "B",
		ShortDescription: "C",
		Price:            100,
		DiscountRate:     10,
		LangID:           langID,
		TopicIDs:         []uuid.UUID{topicID},
	}
	// return wrong type for lang to trigger validation error
	mockParamRepo.On("GetByID", ctx, langID).Return(&models.Parameter{ID: langID, Type: ptrStr("wrong")}, nil).Once()
	_, err := useCase.Create(ctx, req, "", nil, "")
	assert.Error(t, err)

	// valid types then success
	mockParamRepo.ExpectedCalls = nil
	mockPostRepo.ExpectedCalls = nil

	mockParamRepo.On("GetByID", ctx, langID).Return(&models.Parameter{ID: langID, Type: ptrStr("lang")}, nil).Once()
	mockParamRepo.On("GetByID", ctx, topicID).Return(&models.Parameter{ID: topicID, Type: ptrStr("topic")}, nil).Once()

	cID := uuid.New()
	mockPostRepo.On("Create", ctx, uuid.Nil, mock.Anything).
		Return(&models.Post{ID: cID, Title: "A", Description: "B", ShortDescription: "C", Price: 100, DiscountRate: 10, CreatedAt: time.Now(), UpdatedAt: time.Now()}, nil).Once()
	mockParamRepo.On("AssignParametersToModule", ctx, "post", cID, []uuid.UUID{langID}).Return(nil).Once()
	mockParamRepo.On("AssignParametersToModule", ctx, "post", cID, []uuid.UUID{topicID}).Return(nil).Once()

	res, err := useCase.Create(ctx, &postDto.ReqCreatePost{
		Title:            "A",
		Description:      "B",
		ShortDescription: "C",
		Price:            100,
		DiscountRate:     10,
		LangID:           langID,
		TopicIDs:         []uuid.UUID{topicID},
	}, "", nil, "")
	assert.NoError(t, err)
	assert.NotNil(t, res)
	mockPostRepo.AssertExpectations(t)
	mockParamRepo.AssertExpectations(t)
}

func ptrStr(s string) *string { return &s }
