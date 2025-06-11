package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"hydro-habitat/backend/domain"
	"hydro-habitat/backend/store"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestNewTankHandler tests the TankHandler constructor
func TestNewTankHandler(t *testing.T) {
	// --- Arrange ---
	mockStore := new(store.MockTankStore)

	// --- Act ---
	handler := NewTankHandler(mockStore)

	// --- Assert ---
	assert.NotNil(t, handler)
	assert.Equal(t, mockStore, handler.store)
}

func setupTestRouter(mockStore *store.MockTankStore) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	handler := NewTankHandler(mockStore)

	// Register routes the same way as in the main application
	v1 := router.Group("/api/v1")
	tanks := v1.Group("/tanks")
	{
		tanks.POST("", handler.CreateTank)
		tanks.GET("", handler.GetAllTanks)
		tanks.GET("/:id", handler.GetTankByID)
		tanks.PUT("/:id", handler.UpdateTank)
		tanks.DELETE("/:id", handler.DeleteTank)
	}
	return router
}

func TestGetAllTanks_Success(t *testing.T) {
	// --- Arrange ---
	mockStore := new(store.MockTankStore)
	expectedTanks := []domain.Tank{
		{ID: uuid.New(), Name: "Test Tank 1", VolumeLiters: 100, Water: "tap"},
		{ID: uuid.New(), Name: "Test Tank 2", VolumeLiters: 200, Water: "ro"},
	}

	// Expect that the GetAll method will be called and return our test tanks
	mockStore.On("GetAll").Return(expectedTanks, nil)

	router := setupTestRouter(mockStore)
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/tanks", nil)
	w := httptest.NewRecorder()

	// --- Act ---
	router.ServeHTTP(w, req)

	// --- Assert ---
	assert.Equal(t, http.StatusOK, w.Code)

	var actualTanks []domain.Tank
	err := json.Unmarshal(w.Body.Bytes(), &actualTanks)
	assert.NoError(t, err)
	assert.Equal(t, expectedTanks, actualTanks)

	// Verify that the mock was called as expected
	mockStore.AssertExpectations(t)
}

func TestGetAllTanks_Error(t *testing.T) {
	// --- Arrange ---
	mockStore := new(store.MockTankStore)
	expectedError := errors.New("database error")

	// Expect that GetAll will return an error
	mockStore.On("GetAll").Return(nil, expectedError)

	router := setupTestRouter(mockStore)
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/tanks", nil)
	w := httptest.NewRecorder()

	// --- Act ---
	router.ServeHTTP(w, req)

	// --- Assert ---
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockStore.AssertExpectations(t)
}

func TestCreateTank_Success(t *testing.T) {
	// --- Arrange ---
	mockStore := new(store.MockTankStore)
	createDTO := domain.CreateTankDTO{
		Name:         "New Reef",
		VolumeLiters: 500,
		Water:        "rodi",
	}
	expectedTank := &domain.Tank{
		ID:           uuid.New(),
		Name:         createDTO.Name,
		VolumeLiters: createDTO.VolumeLiters,
		Water:        createDTO.Water,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Expect that the Create method will be called with our DTO
	// and return the newly created tank
	mockStore.On("Create", createDTO).Return(expectedTank, nil)

	router := setupTestRouter(mockStore)
	body, _ := json.Marshal(createDTO)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/tanks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// --- Act ---
	router.ServeHTTP(w, req)

	// --- Assert ---
	assert.Equal(t, http.StatusCreated, w.Code)

	var actualTank domain.Tank
	err := json.Unmarshal(w.Body.Bytes(), &actualTank)
	assert.NoError(t, err)
	assert.Equal(t, expectedTank.ID, actualTank.ID)
	assert.Equal(t, expectedTank.Name, actualTank.Name)

	mockStore.AssertExpectations(t)
}

func TestCreateTank_InvalidData(t *testing.T) {
	// --- Arrange ---
	mockStore := new(store.MockTankStore) // No expectations for calls
	router := setupTestRouter(mockStore)

	// Send invalid data (e.g., missing name)
	invalidDTO := map[string]interface{}{"volume_liters": 100}
	body, _ := json.Marshal(invalidDTO)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/tanks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// --- Act ---
	router.ServeHTTP(w, req)

	// --- Assert ---
	assert.Equal(t, http.StatusBadRequest, w.Code)
	// Ensure that no method on the mock was called
	mockStore.AssertNotCalled(t, "Create", mock.Anything)
}

func TestDeleteTank_Success(t *testing.T) {
	// --- Arrange ---
	mockStore := new(store.MockTankStore)
	testID := uuid.New()

	// Expect that Delete will be called with the correct ID and return no error
	mockStore.On("Delete", testID).Return(nil)

	router := setupTestRouter(mockStore)
	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/tanks/"+testID.String(), nil)
	w := httptest.NewRecorder()

	// --- Act ---
	router.ServeHTTP(w, req)

	// --- Assert ---
	assert.Equal(t, http.StatusNoContent, w.Code)
	mockStore.AssertExpectations(t)
}

func TestDeleteTank_Error(t *testing.T) {
	// --- Arrange ---
	mockStore := new(store.MockTankStore)
	testID := uuid.New()
	expectedError := errors.New("cannot delete tank")

	// Expect that Delete will return an error
	mockStore.On("Delete", testID).Return(expectedError)

	router := setupTestRouter(mockStore)
	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/tanks/"+testID.String(), nil)
	w := httptest.NewRecorder()

	// --- Act ---
	router.ServeHTTP(w, req)

	// --- Assert ---
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockStore.AssertExpectations(t)
}

func TestGetTankByID_Success(t *testing.T) {
	// --- Arrange ---
	mockStore := new(store.MockTankStore)
	testID := uuid.New()
	expectedTank := &domain.Tank{
		ID:           testID,
		Name:         "Test Tank",
		VolumeLiters: 100,
		Water:        "tap",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Expect that GetByID will be called with the correct ID and return the tank
	mockStore.On("GetByID", testID).Return(expectedTank, nil)

	router := setupTestRouter(mockStore)
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/tanks/"+testID.String(), nil)
	w := httptest.NewRecorder()

	// --- Act ---
	router.ServeHTTP(w, req)

	// --- Assert ---
	assert.Equal(t, http.StatusOK, w.Code)

	var actualTank domain.Tank
	err := json.Unmarshal(w.Body.Bytes(), &actualTank)
	assert.NoError(t, err)
	assert.Equal(t, expectedTank.ID, actualTank.ID)
	assert.Equal(t, expectedTank.Name, actualTank.Name)

	mockStore.AssertExpectations(t)
}

func TestGetTankByID_NotFound(t *testing.T) {
	// --- Arrange ---
	mockStore := new(store.MockTankStore)
	testID := uuid.New()
	expectedError := errors.New("tank not found")

	// Expect that GetByID will return an error
	mockStore.On("GetByID", testID).Return(nil, expectedError)

	router := setupTestRouter(mockStore)
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/tanks/"+testID.String(), nil)
	w := httptest.NewRecorder()

	// --- Act ---
	router.ServeHTTP(w, req)

	// --- Assert ---
	assert.Equal(t, http.StatusNotFound, w.Code)
	mockStore.AssertExpectations(t)
}

func TestGetTankByID_InvalidUUID(t *testing.T) {
	// --- Arrange ---
	mockStore := new(store.MockTankStore)
	// No expectations since the invalid UUID should be caught before store interaction

	router := setupTestRouter(mockStore)
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/tanks/invalid-uuid", nil)
	w := httptest.NewRecorder()

	// --- Act ---
	router.ServeHTTP(w, req)

	// --- Assert ---
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "Invalid UUID format")

	// Ensure no store method was called
	mockStore.AssertNotCalled(t, "GetByID", mock.Anything)
}

func TestUpdateTank_Success(t *testing.T) {
	// --- Arrange ---
	mockStore := new(store.MockTankStore)
	testID := uuid.New()
	updateDTO := domain.UpdateTankDTO{
		Name:         "Updated Tank",
		VolumeLiters: 150,
		Water:        "ro",
	}
	expectedTank := &domain.Tank{
		ID:           testID,
		Name:         updateDTO.Name,
		VolumeLiters: updateDTO.VolumeLiters,
		Water:        updateDTO.Water,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Expect that Update will be called with the correct ID and DTO
	mockStore.On("Update", testID, updateDTO).Return(expectedTank, nil)

	router := setupTestRouter(mockStore)
	body, _ := json.Marshal(updateDTO)
	req, _ := http.NewRequest(http.MethodPut, "/api/v1/tanks/"+testID.String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// --- Act ---
	router.ServeHTTP(w, req)

	// --- Assert ---
	assert.Equal(t, http.StatusOK, w.Code)

	var actualTank domain.Tank
	err := json.Unmarshal(w.Body.Bytes(), &actualTank)
	assert.NoError(t, err)
	assert.Equal(t, expectedTank.ID, actualTank.ID)
	assert.Equal(t, expectedTank.Name, actualTank.Name)
	assert.Equal(t, expectedTank.VolumeLiters, actualTank.VolumeLiters)

	mockStore.AssertExpectations(t)
}

func TestUpdateTank_InvalidUUID(t *testing.T) {
	// --- Arrange ---
	mockStore := new(store.MockTankStore)
	updateDTO := domain.UpdateTankDTO{
		Name:         "Updated Tank",
		VolumeLiters: 150,
		Water:        "ro",
	}

	router := setupTestRouter(mockStore)
	body, _ := json.Marshal(updateDTO)
	req, _ := http.NewRequest(http.MethodPut, "/api/v1/tanks/invalid-uuid", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// --- Act ---
	router.ServeHTTP(w, req)

	// --- Assert ---
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "Invalid UUID format")

	// Ensure no store method was called
	mockStore.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

func TestUpdateTank_InvalidData(t *testing.T) {
	// --- Arrange ---
	mockStore := new(store.MockTankStore)
	testID := uuid.New()

	router := setupTestRouter(mockStore)
	// Send invalid data (missing required name field)
	invalidDTO := map[string]interface{}{"volume_liters": 100}
	body, _ := json.Marshal(invalidDTO)
	req, _ := http.NewRequest(http.MethodPut, "/api/v1/tanks/"+testID.String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// --- Act ---
	router.ServeHTTP(w, req)

	// --- Assert ---
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Ensure no store method was called
	mockStore.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

func TestUpdateTank_StoreError(t *testing.T) {
	// --- Arrange ---
	mockStore := new(store.MockTankStore)
	testID := uuid.New()
	updateDTO := domain.UpdateTankDTO{
		Name:         "Updated Tank",
		VolumeLiters: 150,
		Water:        "ro",
	}
	expectedError := errors.New("database error")

	// Expect that Update will return an error
	mockStore.On("Update", testID, updateDTO).Return(nil, expectedError)

	router := setupTestRouter(mockStore)
	body, _ := json.Marshal(updateDTO)
	req, _ := http.NewRequest(http.MethodPut, "/api/v1/tanks/"+testID.String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// --- Act ---
	router.ServeHTTP(w, req)

	// --- Assert ---
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockStore.AssertExpectations(t)
}
