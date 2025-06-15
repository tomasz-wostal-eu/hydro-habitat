package store

import (
	"testing"

	"hydro-habitat/backend/config"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestNewTankStore(t *testing.T) {
	// --- Arrange ---
	mockDb, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer func() {
		_ = mockDb.Close() // Ignore error in test cleanup
	}()

	sqlxDB := sqlx.NewDb(mockDb, "sqlmock")
	db := &DB{DB: sqlxDB}

	// --- Act ---
	store := NewTankStore(db)

	// --- Assert ---
	assert.NotNil(t, store)
	assert.IsType(t, &pgTankStore{}, store)
}

func TestNewPostgresStore_InvalidConnectionString(t *testing.T) {
	// --- Arrange ---
	cfg := &config.Config{
		DatabaseURL: "invalid://connection/string",
	}

	// --- Act ---
	db, err := NewPostgresStore(cfg)

	// --- Assert ---
	assert.Error(t, err)
	assert.Nil(t, db)
}

func TestNewPostgresStore_ValidConnection(t *testing.T) {
	// --- Arrange ---
	// Note: This test requires a real database connection
	// In practice, you might want to skip this test or use environment variables
	cfg := &config.Config{
		DatabaseURL: "postgres://user:password@localhost:5432/testdb?sslmode=disable",
	}

	// --- Act ---
	db, err := NewPostgresStore(cfg)

	// --- Assert ---
	if err != nil {
		// If we can't connect (which is expected in CI/testing), that's okay
		t.Skipf("Skipping database connection test: %v", err)
	} else {
		assert.NotNil(t, db)
		assert.NotNil(t, db.DB)
		if err := db.Close(); err != nil {
			t.Logf("Error closing database: %v", err)
		}
	}
}

func TestDB_Close(t *testing.T) {
	// --- Arrange ---
	mockDb, mock, err := sqlmock.New()
	assert.NoError(t, err)

	// Expect Close to be called
	mock.ExpectClose()

	sqlxDB := sqlx.NewDb(mockDb, "sqlmock")
	db := &DB{DB: sqlxDB}

	// --- Act ---
	err = db.Close()

	// --- Assert ---
	assert.NoError(t, err)

	// Verify all expectations were met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
