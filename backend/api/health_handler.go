package api

import (
"net/http"

"github.com/gin-gonic/gin"
)

// HealthResponse represents the health check response
type HealthResponse struct {
	Status string `json:"status" example:"ok"`
}

// HealthCheck godoc
// @Summary Health check
// @Description Health check endpoint to verify the API is running
// @Tags health
// @Produce json
// @Success 200 {object} HealthResponse "Service is healthy"
// @Router /health [get]
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, HealthResponse{
Status: "ok",
})
}
