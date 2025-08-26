package configurator

import (
	"embed"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockICFGBuilder struct {
	mock.Mock
}

func (m *MockICFGBuilder) SetMountPoint(mountPoint string) ICFGBuilder {
	m.Called(mountPoint)
	return m
}

func (m *MockICFGBuilder) SetISOPaths(isoPaths []string) ICFGBuilder {
	m.Called(isoPaths)
	return m
}

func (m *MockICFGBuilder) SetTemplatesFS(grubTemplatesFS embed.FS) ICFGBuilder {
	m.Called(grubTemplatesFS)
	return m
}

func (m *MockICFGBuilder) createGrubCfgFile() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockICFGBuilder) insertHeaderTemplate() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockICFGBuilder) insertIsoSpecificTemplates() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockICFGBuilder) GetResult() error {
	args := m.Called()
	return args.Error(0)
}

func TestGrubConfigurator_Construct_Success(t *testing.T) {
	mockBuilder := new(MockICFGBuilder)
	mockBuilder.On("createGrubCfgFile").Return(nil).Once()
	mockBuilder.On("insertHeaderTemplate").Return(nil).Once()
	mockBuilder.On("insertIsoSpecificTemplates").Return(nil).Once()
	mockBuilder.On("GetResult").Return(nil).Once()

	gc := NewGrubConfigurator(mockBuilder)
	err := gc.Construct()

	assert.NoError(t, err)
	mockBuilder.AssertExpectations(t)
}

func TestGrubConfigurator_Construct_CreateGrubCfgFileError(t *testing.T) {
	mockBuilder := new(MockICFGBuilder)
	expectedErr := errors.New("failed to create grub cfg file")
	mockBuilder.On("createGrubCfgFile").Return(expectedErr).Once()

	gc := NewGrubConfigurator(mockBuilder)
	err := gc.Construct()

	assert.ErrorIs(t, err, expectedErr)
	mockBuilder.AssertExpectations(t)
}

func TestGrubConfigurator_Construct_InsertHeaderTemplateError(t *testing.T) {
	mockBuilder := new(MockICFGBuilder)
	expectedErr := errors.New("failed to insert header template")
	mockBuilder.On("createGrubCfgFile").Return(nil).Once()
	mockBuilder.On("insertHeaderTemplate").Return(expectedErr).Once()

	gc := NewGrubConfigurator(mockBuilder)
	err := gc.Construct()

	assert.ErrorIs(t, err, expectedErr)
	mockBuilder.AssertExpectations(t)
}

func TestGrubConfigurator_Construct_InsertIsoSpecificTemplatesError(t *testing.T) {
	mockBuilder := new(MockICFGBuilder)
	expectedErr := errors.New("failed to insert iso specific templates")
	mockBuilder.On("createGrubCfgFile").Return(nil).Once()
	mockBuilder.On("insertHeaderTemplate").Return(nil).Once()
	mockBuilder.On("insertIsoSpecificTemplates").Return(expectedErr).Once()

	gc := NewGrubConfigurator(mockBuilder)
	err := gc.Construct()

	assert.ErrorIs(t, err, expectedErr)
	mockBuilder.AssertExpectations(t)
}

func TestGrubConfigurator_Construct_GetResultError(t *testing.T) {
	mockBuilder := new(MockICFGBuilder)
	expectedErr := errors.New("failed to get result")
	mockBuilder.On("createGrubCfgFile").Return(nil).Once()
	mockBuilder.On("insertHeaderTemplate").Return(nil).Once()
	mockBuilder.On("insertIsoSpecificTemplates").Return(nil).Once()
	mockBuilder.On("GetResult").Return(expectedErr).Once()

	gc := NewGrubConfigurator(mockBuilder)
	err := gc.Construct()

	assert.ErrorIs(t, err, expectedErr)
	mockBuilder.AssertExpectations(t)
}
