package repository

import (
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
	"github.com/ners1us/order-service/internal/model"
)

type PVZRepository interface {
	CreatePVZ(pvz *model.PVZ) error
	GetPVZs(page, limit int) ([]model.PVZ, error)
	GetAllPVZs() ([]model.PVZ, error)
	GetPVZByID(id string) (*model.PVZ, error)
}

type pvzRepositoryImpl struct {
	db *sql.DB
}

func NewPVZRepository(db *sql.DB) PVZRepository {
	return &pvzRepositoryImpl{db}
}

func (pr *pvzRepositoryImpl) CreatePVZ(pvz *model.PVZ) error {
	query := "INSERT INTO pvzs (id, registration_date, city) VALUES ($1, $2, $3)"
	_, err := pr.db.Exec(query, pvz.ID, pvz.RegistrationDate, pvz.City)
	return err
}

func (pr *pvzRepositoryImpl) GetPVZs(page, limit int) ([]model.PVZ, error) {
	offset := (page - 1) * limit
	query := "SELECT id, registration_date, city FROM pvzs ORDER BY id LIMIT $1 OFFSET $2"
	rows, err := pr.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var pvzs []model.PVZ
	for rows.Next() {
		var pvz model.PVZ
		if err := rows.Scan(&pvz.ID, &pvz.RegistrationDate, &pvz.City); err != nil {
			return nil, err
		}
		pvzs = append(pvzs, pvz)
	}
	return pvzs, nil
}

func (pr *pvzRepositoryImpl) GetAllPVZs() ([]model.PVZ, error) {
	query := "SELECT id, registration_date, city FROM pvzs"
	rows, err := pr.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var pvzs []model.PVZ
	for rows.Next() {
		var pvz model.PVZ
		if err := rows.Scan(&pvz.ID, &pvz.RegistrationDate, &pvz.City); err != nil {
			return nil, err
		}
		pvzs = append(pvzs, pvz)
	}
	return pvzs, nil
}

func (pr *pvzRepositoryImpl) GetPVZByID(id string) (*model.PVZ, error) {
	var pvz model.PVZ
	query := "SELECT id, registration_date, city FROM pvzs WHERE id = $1"
	err := pr.db.QueryRow(query, id).Scan(&pvz.ID, &pvz.RegistrationDate, &pvz.City)
	if errors.Is(err, sql.ErrNoRows) {
		return &model.PVZ{}, nil
	}
	if err != nil {
		return &model.PVZ{}, err
	}
	return &pvz, nil
}
