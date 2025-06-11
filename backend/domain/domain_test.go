package domain

import (
	"encoding/json"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// newValidator creates a new validator instance.
func newValidator() *validator.Validate {
	return validator.New()
}

func TestCreateTankDTO_Validation(t *testing.T) {
	validate := newValidator()

	// Define a valid DTO as a base for tests
	validDTO := CreateTankDTO{
		Name:         "Akwarium Główne",
		VolumeLiters: 120,
		Water:        WaterTypeTap,
	}

	t.Run("Valid case", func(t *testing.T) {
		dto := validDTO // Copy the valid DTO
		err := validate.Struct(dto)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
	})

	t.Run("Missing Name", func(t *testing.T) {
		dto := validDTO
		dto.Name = "" // Required field
		err := validate.Struct(dto)
		if err == nil {
			t.Error("Expected validation error for missing Name field, but got none")
		}
	})

	t.Run("Invalid volume", func(t *testing.T) {
		dto := validDTO
		dto.VolumeLiters = 0 // Required > 0
		err := validate.Struct(dto)
		if err == nil {
			t.Error("Expected validation error for VolumeLiters <= 0, but got none")
		}
	})

	t.Run("Invalid water type", func(t *testing.T) {
		dto := validDTO
		dto.Water = "saltwater" // Value outside allowed ones
		err := validate.Struct(dto)
		if err == nil {
			t.Error("Expected validation error for invalid Water value, but got none")
		}
	})

	t.Run("Missing water type", func(t *testing.T) {
		dto := validDTO
		dto.Water = "" // Required field
		err := validate.Struct(dto)
		if err == nil {
			t.Error("Expected validation error for missing Water field, but got none")
		}
	})
}

// Tests for UpdateTankDTO are analogous, as it has the same validation rules.
func TestUpdateTankDTO_Validation(t *testing.T) {
	validate := newValidator()

	validDTO := UpdateTankDTO{
		Name:         "Akwarium Zaktualizowane",
		VolumeLiters: 240,
		Water:        WaterTypeRO,
	}

	t.Run("Valid case", func(t *testing.T) {
		err := validate.Struct(validDTO)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
	})

	t.Run("Missing Name", func(t *testing.T) {
		dto := validDTO
		dto.Name = ""
		err := validate.Struct(dto)
		if err == nil {
			t.Error("Expected validation error for missing Name field, but got none")
		}
	})
}

// TestWaterType_Constants tests the WaterType constants
func TestWaterType_Constants(t *testing.T) {
	tests := []struct {
		name     string
		value    WaterType
		expected string
	}{
		{"Tap water", WaterTypeTap, "tap"},
		{"RO water", WaterTypeRO, "ro"},
		{"RODI water", WaterTypeRODI, "rodi"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.value) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, string(tt.value))
			}
		})
	}
}

// TestTank_JSONSerialization tests JSON marshaling/unmarshaling
func TestTank_JSONSerialization(t *testing.T) {
	tank := Tank{
		ID:           uuid.New(),
		Name:         "Test Tank",
		VolumeLiters: 100,
		Water:        WaterTypeTap,
	}

	// Test marshaling
	jsonData, err := json.Marshal(tank)
	if err != nil {
		t.Fatalf("Failed to marshal tank to JSON: %v", err)
	}

	// Test unmarshaling
	var unmarshaledTank Tank
	err = json.Unmarshal(jsonData, &unmarshaledTank)
	if err != nil {
		t.Fatalf("Failed to unmarshal tank from JSON: %v", err)
	}

	// Verify fields
	if unmarshaledTank.ID != tank.ID {
		t.Errorf("Expected ID %s, got %s", tank.ID, unmarshaledTank.ID)
	}
	if unmarshaledTank.Name != tank.Name {
		t.Errorf("Expected Name %s, got %s", tank.Name, unmarshaledTank.Name)
	}
	if unmarshaledTank.VolumeLiters != tank.VolumeLiters {
		t.Errorf("Expected VolumeLiters %d, got %d", tank.VolumeLiters, unmarshaledTank.VolumeLiters)
	}
	if unmarshaledTank.Water != tank.Water {
		t.Errorf("Expected Water %s, got %s", tank.Water, unmarshaledTank.Water)
	}
}

// TestCreateTankDTO_EdgeCases tests edge cases for CreateTankDTO validation
func TestCreateTankDTO_EdgeCases(t *testing.T) {
	validate := newValidator()

	t.Run("Minimum valid volume", func(t *testing.T) {
		dto := CreateTankDTO{
			Name:         "Small Tank",
			VolumeLiters: 1, // Minimum valid value
			Water:        WaterTypeTap,
		}
		err := validate.Struct(dto)
		if err != nil {
			t.Errorf("Expected no error for minimum valid volume, got: %v", err)
		}
	})

	t.Run("Very large volume", func(t *testing.T) {
		dto := CreateTankDTO{
			Name:         "Huge Tank",
			VolumeLiters: 100000, // Very large value
			Water:        WaterTypeTap,
		}
		err := validate.Struct(dto)
		if err != nil {
			t.Errorf("Expected no error for large volume, got: %v", err)
		}
	})

	t.Run("All optional fields provided", func(t *testing.T) {
		room := "Fish Room"
		rackLocation := "A1"
		inventoryNumber := "T-001"
		notes := "Test notes"

		dto := CreateTankDTO{
			Name:            "Complete Tank",
			Room:            &room,
			RackLocation:    &rackLocation,
			VolumeLiters:    100,
			InventoryNumber: &inventoryNumber,
			Water:           WaterTypeRO,
			Notes:           &notes,
		}
		err := validate.Struct(dto)
		if err != nil {
			t.Errorf("Expected no error with all fields provided, got: %v", err)
		}
	})
}
