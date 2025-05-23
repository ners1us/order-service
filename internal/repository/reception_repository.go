package repository

import (
	"database/sql"
	"errors"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/ners1us/order-service/internal/model"
	"time"
)

type ReceptionRepository interface {
	CreateReception(reception *model.Reception) error
	GetLastReceptionByPVZID(pvzID string) (*model.Reception, error)
	UpdateReceptionStatus(id string, status string) error
	GetReceptionsByPVZIDsAndDate(pvzIDs []string, startDate, endDate time.Time) ([]model.Reception, error)
}

type receptionRepositoryImpl struct {
	db *sql.DB
}

func NewReceptionRepository(db *sql.DB) ReceptionRepository {
	return &receptionRepositoryImpl{db}
}

func (rr *receptionRepositoryImpl) CreateReception(reception *model.Reception) error {
	query := "INSERT INTO receptions (id, date_time, pvz_id, status) VALUES ($1, $2, $3, $4)"
	_, err := rr.db.Exec(query, reception.ID, reception.DateTime, reception.PVZID, reception.Status)
	return err
}

func (rr *receptionRepositoryImpl) GetLastReceptionByPVZID(pvzID string) (*model.Reception, error) {
	var reception model.Reception
	query := "SELECT id, date_time, pvz_id, status FROM receptions WHERE pvz_id = $1 ORDER BY date_time DESC LIMIT 1"
	err := rr.db.QueryRow(query, pvzID).Scan(&reception.ID, &reception.DateTime, &reception.PVZID, &reception.Status)
	if errors.Is(err, sql.ErrNoRows) {
		return &model.Reception{}, nil
	}
	return &reception, err
}

func (rr *receptionRepositoryImpl) UpdateReceptionStatus(id string, status string) error {
	query := "UPDATE receptions SET status = $1 WHERE id = $2"
	_, err := rr.db.Exec(query, status, id)
	return err
}

func (rr *receptionRepositoryImpl) GetReceptionsByPVZIDsAndDate(pvzIDs []string, startDate, endDate time.Time) ([]model.Reception, error) {
	query := "SELECT id, date_time, pvz_id, status FROM receptions WHERE pvz_id = ANY($1) AND date_time BETWEEN $2 AND $3"
	rows, err := rr.db.Query(query, pq.Array(pvzIDs), startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var receptions []model.Reception
	for rows.Next() {
		var reception model.Reception
		if err := rows.Scan(&reception.ID, &reception.DateTime, &reception.PVZID, &reception.Status); err != nil {
			return nil, err
		}
		receptions = append(receptions, reception)
	}
	return receptions, nil
}
