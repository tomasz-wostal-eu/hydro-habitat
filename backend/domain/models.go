package domain

import (
	"time"

	"github.com/google/uuid"
)

type WaterType string

const (
	WaterTypeTap  WaterType = "tap"
	WaterTypeRO   WaterType = "ro"
	WaterTypeRODI WaterType = "rodi"
)

type Tank struct {
	ID              uuid.UUID `json:"id" db:"id"`
	Name            string    `json:"name" db:"name"`
	Room            *string   `json:"room,omitempty" db:"room"`
	RackLocation    *string   `json:"rack_location,omitempty" db:"rack_location"`
	VolumeLiters    int       `json:"volume_liters" db:"volume_liters"`
	InventoryNumber *string   `json:"inventory_number,omitempty" db:"inventory_number"`
	Water           WaterType `json:"water" db:"water"`
	Notes           *string   `json:"notes,omitempty" db:"notes"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

type CreateTankDTO struct {
	Name            string    `json:"name" binding:"required" validate:"required"`
	Room            *string   `json:"room"`
	RackLocation    *string   `json:"rack_location"`
	VolumeLiters    int       `json:"volume_liters" binding:"required,gt=0" validate:"required,gt=0"`
	InventoryNumber *string   `json:"inventory_number"`
	Water           WaterType `json:"water" binding:"required,oneof=tap ro rodi" validate:"required,oneof=tap ro rodi"`
	Notes           *string   `json:"notes"`
}

type UpdateTankDTO struct {
	Name            string    `json:"name" binding:"required" validate:"required"`
	Room            *string   `json:"room"`
	RackLocation    *string   `json:"rack_location"`
	VolumeLiters    int       `json:"volume_liters" binding:"required,gt=0" validate:"required,gt=0"`
	InventoryNumber *string   `json:"inventory_number"`
	Water           WaterType `json:"water" binding:"required,oneof=tap ro rodi" validate:"required,oneof=tap ro rodi"`
	Notes           *string   `json:"notes"`
}
