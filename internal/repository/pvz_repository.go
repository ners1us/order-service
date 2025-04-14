package repository

import (
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
	"github.com/ners1us/order-service/internal/models"
)

type PVZRepository interface {
	CreatePVZ(pvz *models.PVZ) error
	GetPVZs(page, limit int) ([]models.PVZ, error)
	GetAllPVZs() ([]models.PVZ, error)
	GetPVZByID(id string) (*models.PVZ, error)
}

type pvzRepositoryImpl struct {
	db *sql.DB
}

func NewPVZRepository(db *sql.DB) PVZRepository {
	return &pvzRepositoryImpl{db}
}

func (pr *pvzRepositoryImpl) CreatePVZ(pvz *models.PVZ) error {
	_, err := pr.db.Exec("INSERT INTO pvzs (id, registration_date, city) VALUES ($1, $2, $3)", pvz.ID, pvz.RegistrationDate, pvz.City)
	return err
}

func (pr *pvzRepositoryImpl) GetPVZs(page, limit int) ([]models.PVZ, error) {
	offset := (page - 1) * limit
	rows, err := pr.db.Query("SELECT id, registration_date, city FROM pvzs ORDER BY id LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var pvzs []models.PVZ
	for rows.Next() {
		var pvz models.PVZ
		if err := rows.Scan(&pvz.ID, &pvz.RegistrationDate, &pvz.City); err != nil {
			return nil, err
		}
		pvzs = append(pvzs, pvz)
	}
	return pvzs, nil
}

func (pr *pvzRepositoryImpl) GetAllPVZs() ([]models.PVZ, error) {
	rows, err := pr.db.Query("SELECT id, registration_date, city FROM pvzs")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var pvzs []models.PVZ
	for rows.Next() {
		var pvz models.PVZ
		if err := rows.Scan(&pvz.ID, &pvz.RegistrationDate, &pvz.City); err != nil {
			return nil, err
		}
		pvzs = append(pvzs, pvz)
	}
	return pvzs, nil
}

func (pr *pvzRepositoryImpl) GetPVZByID(id string) (*models.PVZ, error) {
	var pvz models.PVZ
	err := pr.db.QueryRow("SELECT id, registration_date, city FROM pvzs WHERE id = $1", id).
		Scan(&pvz.ID, &pvz.RegistrationDate, &pvz.City)
	if errors.Is(err, sql.ErrNoRows) {
		return &models.PVZ{}, nil
	}
	if err != nil {
		return &models.PVZ{}, err
	}
	return &pvz, nil
}
