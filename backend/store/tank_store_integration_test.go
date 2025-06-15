//go:build integration

package store_test

import (
	"hydro-habitat/backend/config"
	"hydro-habitat/backend/domain"
	"hydro-habitat/backend/store"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

// TankStoreIntegrationSuite is a suite of integration tests for TankStore.
type TankStoreIntegrationSuite struct {
	suite.Suite
	db    *store.DB
	store store.TankStore
}

// SetupSuite is run once before all tests in this suite.
// Its purpose is to establish a connection to the database.
func (s *TankStoreIntegrationSuite) SetupSuite() {
	cfg, err := config.Load()
	s.Require().NoError(err, "Failed to load config")

	// Check if we are in a CI environment
	isCI := os.Getenv("CI") == "true"

	// If running locally and DATABASE_URL is not set, skip the test
	if os.Getenv("DATABASE_URL") == "" && !isCI {
		s.T().Skip("Skipping integration test: DATABASE_URL not set")
	}

	// Log connection attempt details (helpful for debugging)
	s.T().Logf("Attempting to connect to database: %s (masked credentials)", 
		strings.Replace(cfg.DatabaseURL, cfg.DatabaseURL, "[REDACTED]", 1))

	// Try to connect to the database
	db, err := store.NewPostgresStore(cfg)
	if err != nil {
		if isCI {
			// In CI, we want to fail rather than skip
			s.Require().NoError(err, "Failed to connect to database in CI environment: %v", err)
		} else {
			// In local dev, skip with informative message
			s.T().Skipf("Skipping integration test: %v", err)
		}
	}

	s.db = db
	s.store = store.NewTankStore(db)
}

// TearDownSuite is run once after all tests.
// It closes the database connection.
func (s *TankStoreIntegrationSuite) TearDownSuite() {
	s.db.Close()
}

// TestTankStoreIntegration runs the test suite.
func TestTankStoreIntegration(t *testing.T) {
	suite.Run(t, new(TankStoreIntegrationSuite))
}

// TestTankCRUD verifies the complete lifecycle (CRUD) for a tank.
func (s *TankStoreIntegrationSuite) TestTankCRUD() {
	var newTankID uuid.UUID

	// Step 1: Create
	s.Run("1_CreateTank", func() {
		createDTO := domain.CreateTankDTO{
			Name:         "Integration Test Tank",
			VolumeLiters: 120,
			Water:        "ro",
			Room:         &[]string{"Test Room"}[0],
		}

		createdTank, err := s.store.Create(createDTO)
		s.Require().NoError(err)
		s.Require().NotNil(createdTank)

		newTankID = createdTank.ID
		s.Equal(createDTO.Name, createdTank.Name)
		s.Equal(createDTO.VolumeLiters, createdTank.VolumeLiters)
	})

	// Step 2: Read (GetByID)
	s.Run("2_GetTankByID", func() {
		s.Require().NotEqual(uuid.Nil, newTankID, "Tank ID should be set from create step")

		foundTank, err := s.store.GetByID(newTankID)
		s.Require().NoError(err)
		s.Require().NotNil(foundTank)
		s.Equal("Integration Test Tank", foundTank.Name)
	})

	// Step 3: Update
	s.Run("3_UpdateTank", func() {
		s.Require().NotEqual(uuid.Nil, newTankID)
		updatedNotes := "These are updated notes."
		updateDTO := domain.UpdateTankDTO{
			Name:         "Updated Test Tank",
			VolumeLiters: 150,
			Water:        "rodi",
			Notes:        &updatedNotes,
		}

		updatedTank, err := s.store.Update(newTankID, updateDTO)
		s.Require().NoError(err)
		s.Require().NotNil(updatedTank)
		s.Equal("Updated Test Tank", updatedTank.Name)
		s.Equal(150, updatedTank.VolumeLiters)
		s.Equal(domain.WaterTypeRODI, updatedTank.Water)
		s.NotNil(updatedTank.Notes)
		s.Equal(updatedNotes, *updatedTank.Notes)
	})

	// Step 4: Delete
	s.Run("4_DeleteTank", func() {
		s.Require().NotEqual(uuid.Nil, newTankID)

		err := s.store.Delete(newTankID)
		s.Require().NoError(err)
	})

	// Step 5: Verify Deletion
	s.Run("5_VerifyDeletion", func() {
		s.Require().NotEqual(uuid.Nil, newTankID)

		_, err := s.store.GetByID(newTankID)
		s.Require().Error(err, "Expected an error when getting a deleted tank")
	})
}
