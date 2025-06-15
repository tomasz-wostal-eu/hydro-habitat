package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"hydro-habitat/backend/domain"
	"hydro-habitat/backend/store"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewRouter(t *testing.T) {
	// --- Arrange ---
	// Create a mock database connection
	mockDb, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer func() {
		_ = mockDb.Close() // Ignore error in test cleanup
	}()

	sqlxDB := sqlx.NewDb(mockDb, "sqlmock")
	db := &store.DB{DB: sqlxDB}

	// --- Act ---
	router := NewRouter(db)

	// --- Assert ---
	assert.NotNil(t, router)
}

func TestRouter_CORSMiddleware(t *testing.T) {
	// --- Arrange ---
	mockDb, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer func() {
		_ = mockDb.Close() // Ignore error in test cleanup
	}()

	sqlxDB := sqlx.NewDb(mockDb, "sqlmock")
	db := &store.DB{DB: sqlxDB}
	router := NewRouter(db)

	// --- Act ---
	req, _ := http.NewRequest(http.MethodOptions, "/api/v1/tanks", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// --- Assert ---
	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
	assert.Contains(t, w.Header().Get("Access-Control-Allow-Headers"), "Content-Type")
	assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "POST")
}

func TestRouter_HealthEndpoint(t *testing.T) {
	// --- Arrange ---
	mockDb, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer func() {
		_ = mockDb.Close() // Ignore error in test cleanup
	}()

	sqlxDB := sqlx.NewDb(mockDb, "sqlmock")
	db := &store.DB{DB: sqlxDB}
	router := NewRouter(db)

	// --- Act ---
	req, _ := http.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// --- Assert ---
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"status":"ok"`)
}

func TestRouter_SwaggerEndpoint(t *testing.T) {
	// --- Arrange ---
	mockDb, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer func() {
		_ = mockDb.Close() // Ignore error in test cleanup
	}()

	sqlxDB := sqlx.NewDb(mockDb, "sqlmock")
	db := &store.DB{DB: sqlxDB}
	router := NewRouter(db)

	// --- Act & Assert ---
	// Test that the swagger wildcard route is registered
	// In test environment, swagger docs may not be available, but the route should be registered
	// We check that the route pattern exists by looking at gin's registered routes
	routes := router.Routes()

	swaggerRouteFound := false
	for _, route := range routes {
		if route.Path == "/swagger/*any" && route.Method == "GET" {
			swaggerRouteFound = true
			break
		}
	}

	assert.True(t, swaggerRouteFound, "Swagger route should be registered")
}

func TestRouter_TankRoutes(t *testing.T) {
	// --- Arrange ---
	// Create a simple test store that doesn't need database mocking
	mockStore := &store.MockTankStore{}

	// Set up basic mock responses for route testing
	mockStore.On("GetAll").Return([]domain.Tank{}, nil)
	mockStore.On("GetByID", mock.AnythingOfType("uuid.UUID")).Return(
		&domain.Tank{
			ID:   uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
			Name: "Test Tank",
		}, nil)
	mockStore.On("Delete", mock.AnythingOfType("uuid.UUID")).Return(nil)

	// Create router with mock store instead of real database
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	tankHandler := NewTankHandler(mockStore)

	// Apply CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Register tank routes
	v1 := router.Group("/api/v1")
	tanks := v1.Group("/tanks")
	{
		tanks.POST("", tankHandler.CreateTank)
		tanks.GET("", tankHandler.GetAllTanks)
		tanks.GET("/:id", tankHandler.GetTankByID)
		tanks.PUT("/:id", tankHandler.UpdateTank)
		tanks.DELETE("/:id", tankHandler.DeleteTank)
	}

	tests := []struct {
		name           string
		method         string
		path           string
		expectNotFound bool
	}{
		{"GET all tanks", "GET", "/api/v1/tanks", false},
		{"POST tank", "POST", "/api/v1/tanks", false}, // Will fail validation, but route exists
		{"GET tank by ID", "GET", "/api/v1/tanks/550e8400-e29b-41d4-a716-446655440000", false},
		{"PUT tank", "PUT", "/api/v1/tanks/550e8400-e29b-41d4-a716-446655440000", false}, // Will fail validation, but route exists
		{"DELETE tank", "DELETE", "/api/v1/tanks/550e8400-e29b-41d4-a716-446655440000", false},
		{"Non-existent route", "GET", "/api/v1/nonexistent", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// --- Act ---
			req, _ := http.NewRequest(tt.method, tt.path, nil)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// --- Assert ---
			if tt.expectNotFound {
				assert.Equal(t, http.StatusNotFound, w.Code)
			} else {
				// The route should exist (though it might return errors due to invalid data/mocks)
				// We're testing route registration, not business logic
				assert.NotEqual(t, http.StatusNotFound, w.Code, "Route should be registered")
			}
		})
	}
}
