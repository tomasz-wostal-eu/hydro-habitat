package api

import (
	"hydro-habitat/backend/domain"
	"hydro-habitat/backend/store"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TankHandler struct {
	store store.TankStore
}

func NewTankHandler(store store.TankStore) *TankHandler {
	return &TankHandler{store: store}
}

// CreateTank godoc
// @Summary Create a new tank
// @Description Adds a new tank to the system inventory.
// @Tags tanks
// @Accept json
// @Produce json
// @Param tank body domain.CreateTankDTO true "Tank to create"
// @Success 201 {object} domain.Tank
// @Failure 400 {object} map[string]string
// @Router /tanks [post]
func (h *TankHandler) CreateTank(c *gin.Context) {
	var dto domain.CreateTankDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tank, err := h.store.Create(dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tank"})
		return
	}
	c.JSON(http.StatusCreated, tank)
}

// GetAllTanks godoc
// @Summary Get all tanks
// @Description Retrieves a list of all tanks.
// @Tags tanks
// @Produce json
// @Success 200 {array} domain.Tank
// @Router /tanks [get]
func (h *TankHandler) GetAllTanks(c *gin.Context) {
	tanks, err := h.store.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tanks"})
		return
	}
	c.JSON(http.StatusOK, tanks)
}

// GetTankByID godoc
// @Summary Get a tank by ID
// @Description Retrieves a single tank by its UUID.
// @Tags tanks
// @Produce json
// @Param id path string true "Tank ID" format(uuid)
// @Success 200 {object} domain.Tank
// @Failure 404 {object} map[string]string
// @Router /tanks/{id} [get]
func (h *TankHandler) GetTankByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}
	tank, err := h.store.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tank not found"})
		return
	}
	c.JSON(http.StatusOK, tank)
}

// UpdateTank godoc
// @Summary Update a tank
// @Description Updates an existing tank.
// @Tags tanks
// @Accept json
// @Produce json
// @Param id path string true "Tank ID" format(uuid)
// @Param tank body domain.UpdateTankDTO true "Tank data to update"
// @Success 200 {object} domain.Tank
// @Failure 400 {object} map[string]string
// @Router /tanks/{id} [put]
func (h *TankHandler) UpdateTank(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}
	var dto domain.UpdateTankDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tank, err := h.store.Update(id, dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update tank"})
		return
	}
	c.JSON(http.StatusOK, tank)
}

// DeleteTank godoc
// @Summary Delete a tank
// @Description Deletes a tank by its UUID.
// @Tags tanks
// @Param id path string true "Tank ID" format(uuid)
// @Success 204 "No Content"
// @Router /tanks/{id} [delete]
func (h *TankHandler) DeleteTank(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}
	if err := h.store.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete tank"})
		return
	}
	c.Status(http.StatusNoContent)
}
