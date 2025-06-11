package store

import (
	"hydro-habitat/backend/domain"

	"github.com/google/uuid"
)

type TankStore interface {
	Create(dto domain.CreateTankDTO) (*domain.Tank, error)
	GetAll() ([]domain.Tank, error)
	GetByID(id uuid.UUID) (*domain.Tank, error)
	Update(id uuid.UUID, dto domain.UpdateTankDTO) (*domain.Tank, error)
	Delete(id uuid.UUID) error
}

type pgTankStore struct {
	db *DB
}

func NewTankStore(db *DB) TankStore {
	return &pgTankStore{db: db}
}

func (s *pgTankStore) Create(dto domain.CreateTankDTO) (*domain.Tank, error) {
	tank := &domain.Tank{}
	query := `INSERT INTO tanks (name, room, rack_location, volume_liters, inventory_number, water, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, name, room, rack_location, volume_liters, inventory_number, water, notes, created_at, updated_at`

	err := s.db.Get(tank, query, dto.Name, dto.Room, dto.RackLocation, dto.VolumeLiters, dto.InventoryNumber, dto.Water, dto.Notes)
	if err != nil {
		return nil, err // POPRAWKA: Jeśli jest błąd, zwróć nil dla tanka
	}
	return tank, nil
}

func (s *pgTankStore) GetAll() ([]domain.Tank, error) {
	var tanks []domain.Tank
	err := s.db.Select(&tanks, `SELECT * FROM tanks ORDER BY created_at DESC`)
	return tanks, err
}

func (s *pgTankStore) GetByID(id uuid.UUID) (*domain.Tank, error) {
	tank := &domain.Tank{}
	err := s.db.Get(tank, `SELECT * FROM tanks WHERE id = $1`, id)
	if err != nil {
		return nil, err // POPRAWKA: Jeśli jest błąd, zwróć nil dla tanka
	}
	return tank, nil
}

func (s *pgTankStore) Update(id uuid.UUID, dto domain.UpdateTankDTO) (*domain.Tank, error) {
	tank := &domain.Tank{}
	query := `UPDATE tanks SET name = $1, room = $2, rack_location = $3, volume_liters = $4, inventory_number = $5, water = $6, notes = $7, updated_at = NOW()
		WHERE id = $8
		RETURNING id, name, room, rack_location, volume_liters, inventory_number, water, notes, created_at, updated_at`

	err := s.db.Get(tank, query, dto.Name, dto.Room, dto.RackLocation, dto.VolumeLiters, dto.InventoryNumber, dto.Water, dto.Notes, id)
	if err != nil {
		return nil, err // POPRAWKA: Jeśli jest błąd, zwróć nil dla tanka
	}
	return tank, nil
}

func (s *pgTankStore) Delete(id uuid.UUID) error {
	_, err := s.db.Exec(`DELETE FROM tanks WHERE id = $1`, id)
	return err
}
