package store

import (
	"hydro-habitat/backend/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockTankStore is a mock for the TankStore interface.
// We use it to test handlers in isolation from the real database.
type MockTankStore struct {
	mock.Mock
}

// Create is a mocked implementation of the Create method.
func (m *MockTankStore) Create(dto domain.CreateTankDTO) (*domain.Tank, error) {
	args := m.Called(dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Tank), args.Error(1)
}

// GetAll is a mocked implementation of the GetAll method.
func (m *MockTankStore) GetAll() ([]domain.Tank, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Tank), args.Error(1)
}

// GetByID is a mocked implementation of the GetByID method.
func (m *MockTankStore) GetByID(id uuid.UUID) (*domain.Tank, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Tank), args.Error(1)
}

// Update is a mocked implementation of the Update method.
func (m *MockTankStore) Update(id uuid.UUID, dto domain.UpdateTankDTO) (*domain.Tank, error) {
	args := m.Called(id, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Tank), args.Error(1)
}

// Delete is a mocked implementation of the Delete method.
func (m *MockTankStore) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}
