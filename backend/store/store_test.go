package store

import (
	"database/sql"
	"regexp"
	"testing"
	"time"

	"hydro-habitat/backend/domain"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

// UWAGA: Usunęliśmy linię "type DB = sqlx.DB", ponieważ typ DB istnieje już w pakiecie.

func TestPgTankStore_Create(t *testing.T) {
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Wystąpił błąd podczas tworzenia mocka bazy danych: %s", err)
	}
	defer func() {
		_ = mockDb.Close() // Ignore error in test cleanup
	}()

	sqlxDB := sqlx.NewDb(mockDb, "sqlmock")
	// POPRAWKA: Tworzymy instancję prawdziwej struktury DB z pakietu store
	db := &DB{DB: sqlxDB}
	store := NewTankStore(db)

	dto := domain.CreateTankDTO{
		Name:         "Test Tank",
		VolumeLiters: 100,
		Water:        domain.WaterTypeTap,
	}

	t.Run("Success", func(t *testing.T) {
		now := time.Now()
		newId := uuid.New()

		query := regexp.QuoteMeta(`
			INSERT INTO tanks (name, room, rack_location, volume_liters, inventory_number, water, notes)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING id, name, room, rack_location, volume_liters, inventory_number, water, notes, created_at, updated_at
		`)

		rows := sqlmock.NewRows([]string{"id", "name", "volume_liters", "water", "created_at", "updated_at"}).
			AddRow(newId, dto.Name, dto.VolumeLiters, dto.Water, now, now)

		mock.ExpectQuery(query).
			WithArgs(dto.Name, dto.Room, dto.RackLocation, dto.VolumeLiters, dto.InventoryNumber, dto.Water, dto.Notes).
			WillReturnRows(rows)

		createdTank, err := store.Create(dto)

		assert.NoError(t, err)
		assert.NotNil(t, createdTank)
		assert.Equal(t, newId, createdTank.ID)
		assert.Equal(t, dto.Name, createdTank.Name)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Database Error", func(t *testing.T) {
		query := regexp.QuoteMeta(`INSERT INTO tanks`)

		mock.ExpectQuery(query).
			WithArgs(dto.Name, dto.Room, dto.RackLocation, dto.VolumeLiters, dto.InventoryNumber, dto.Water, dto.Notes).
			WillReturnError(sql.ErrConnDone)

		createdTank, err := store.Create(dto)

		assert.Error(t, err)
		assert.Nil(t, createdTank)
		assert.Equal(t, sql.ErrConnDone, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestPgTankStore_GetAll(t *testing.T) {
	mockDb, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer func() {
		_ = mockDb.Close() // Ignore error in test cleanup
	}()

	sqlxDB := sqlx.NewDb(mockDb, "sqlmock")
	// POPRAWKA:
	db := &DB{DB: sqlxDB}
	store := NewTankStore(db)

	t.Run("Success", func(t *testing.T) {
		query := regexp.QuoteMeta(`SELECT * FROM tanks ORDER BY created_at DESC`)

		rows := sqlmock.NewRows([]string{"id", "name"}).
			AddRow(uuid.New(), "Tank 1").
			AddRow(uuid.New(), "Tank 2")

		mock.ExpectQuery(query).WillReturnRows(rows)

		tanks, err := store.GetAll()

		assert.NoError(t, err)
		assert.Len(t, tanks, 2)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Database Error", func(t *testing.T) {
		query := regexp.QuoteMeta(`SELECT * FROM tanks`)
		mock.ExpectQuery(query).WillReturnError(sql.ErrNoRows)

		tanks, err := store.GetAll()

		assert.Error(t, err)
		assert.Nil(t, tanks)
		assert.Equal(t, sql.ErrNoRows, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestPgTankStore_GetByID(t *testing.T) {
	mockDb, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer func() {
		_ = mockDb.Close() // Ignore error in test cleanup
	}()

	sqlxDB := sqlx.NewDb(mockDb, "sqlmock")
	// POPRAWKA:
	db := &DB{DB: sqlxDB}
	store := NewTankStore(db)

	tankId := uuid.New()

	t.Run("Success", func(t *testing.T) {
		query := regexp.QuoteMeta(`SELECT * FROM tanks WHERE id = $1`)

		rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(tankId, "Found Tank")
		mock.ExpectQuery(query).WithArgs(tankId).WillReturnRows(rows)

		tank, err := store.GetByID(tankId)

		assert.NoError(t, err)
		assert.NotNil(t, tank)
		assert.Equal(t, tankId, tank.ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Not Found", func(t *testing.T) {
		query := regexp.QuoteMeta(`SELECT * FROM tanks WHERE id = $1`)

		mock.ExpectQuery(query).WithArgs(tankId).WillReturnError(sql.ErrNoRows)

		tank, err := store.GetByID(tankId)

		assert.Error(t, err)
		assert.Equal(t, sql.ErrNoRows, err)
		assert.Nil(t, tank)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestPgTankStore_Update(t *testing.T) {
	mockDb, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer func() {
		_ = mockDb.Close() // Ignore error in test cleanup
	}()

	sqlxDB := sqlx.NewDb(mockDb, "sqlmock")
	// POPRAWKA:
	db := &DB{DB: sqlxDB}
	store := NewTankStore(db)

	tankId := uuid.New()
	dto := domain.UpdateTankDTO{
		Name:         "Updated Name",
		VolumeLiters: 150,
		Water:        domain.WaterTypeRO,
	}

	t.Run("Success", func(t *testing.T) {
		query := regexp.QuoteMeta(`
			UPDATE tanks SET name = $1, room = $2, rack_location = $3, volume_liters = $4, inventory_number = $5, water = $6, notes = $7, updated_at = NOW()
			WHERE id = $8
			RETURNING id, name, room, rack_location, volume_liters, inventory_number, water, notes, created_at, updated_at
		`)

		rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(tankId, dto.Name)
		mock.ExpectQuery(query).
			WithArgs(dto.Name, dto.Room, dto.RackLocation, dto.VolumeLiters, dto.InventoryNumber, dto.Water, dto.Notes, tankId).
			WillReturnRows(rows)

		updatedTank, err := store.Update(tankId, dto)

		assert.NoError(t, err)
		assert.NotNil(t, updatedTank)
		assert.Equal(t, tankId, updatedTank.ID)
		assert.Equal(t, "Updated Name", updatedTank.Name)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Database Error", func(t *testing.T) {
		query := regexp.QuoteMeta(`UPDATE tanks SET`)
		mock.ExpectQuery(query).WillReturnError(sql.ErrConnDone)

		updatedTank, err := store.Update(tankId, dto)

		assert.Error(t, err)
		assert.Nil(t, updatedTank)
		assert.Equal(t, sql.ErrConnDone, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestPgTankStore_Delete(t *testing.T) {
	mockDb, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer func() {
		_ = mockDb.Close() // Ignore error in test cleanup
	}()

	sqlxDB := sqlx.NewDb(mockDb, "sqlmock")
	// POPRAWKA:
	db := &DB{DB: sqlxDB}
	store := NewTankStore(db)

	idToDelete := uuid.New()

	t.Run("Success", func(t *testing.T) {
		query := regexp.QuoteMeta(`DELETE FROM tanks WHERE id = $1`)

		mock.ExpectExec(query).
			WithArgs(idToDelete).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := store.Delete(idToDelete)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Database Error", func(t *testing.T) {
		query := regexp.QuoteMeta(`DELETE FROM tanks WHERE id = $1`)

		mock.ExpectExec(query).
			WithArgs(idToDelete).
			WillReturnError(sql.ErrConnDone)

		err := store.Delete(idToDelete)

		assert.Error(t, err)
		assert.Equal(t, sql.ErrConnDone, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
