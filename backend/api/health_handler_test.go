package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	// --- Arrange ---
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/health", HealthCheck)

	// --- Act ---
	req, _ := http.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// --- Assert ---
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	assert.JSONEq(t, `{"status":"ok"}`, w.Body.String())
}

func TestHealthCheck_ResponseStructure(t *testing.T) {
	// --- Arrange ---
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/health", HealthCheck)

	// --- Act ---
	req, _ := http.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// --- Assert ---
	assert.Equal(t, http.StatusOK, w.Code)
	
	// Verify the response contains the expected JSON structure
	expectedResponse := `{"status":"ok"}`
	assert.JSONEq(t, expectedResponse, w.Body.String())
	
	// Verify Content-Type header
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
}
