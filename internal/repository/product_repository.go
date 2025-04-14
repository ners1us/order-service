package repository

import (
	"database/sql"
	"errors"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/ners1us/order-service/internal/models"
)

type ProductRepository interface {
	CreateProduct(product *models.Product) error
	GetLastProductByReceptionID(receptionID string) (*models.Product, error)
	DeleteProduct(id string) error
	GetProductsByReceptionIDs(receptionIDs []string) ([]models.Product, error)
}

type productRepositoryImpl struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepositoryImpl{db}
}

func (pr *productRepositoryImpl) CreateProduct(product *models.Product) error {
	_, err := pr.db.Exec("INSERT INTO products (id, date_time, type, reception_id) VALUES ($1, $2, $3, $4)", product.ID, product.DateTime, product.Type, product.ReceptionID)
	return err
}

func (pr *productRepositoryImpl) GetLastProductByReceptionID(receptionID string) (*models.Product, error) {
	var product models.Product
	err := pr.db.QueryRow("SELECT id, date_time, type, reception_id FROM products WHERE reception_id = $1 ORDER BY date_time DESC LIMIT 1", receptionID).Scan(&product.ID, &product.DateTime, &product.Type, &product.ReceptionID)
	if errors.Is(err, sql.ErrNoRows) {
		return &models.Product{}, err
	}
	return &product, nil
}

func (pr *productRepositoryImpl) DeleteProduct(id string) error {
	_, err := pr.db.Exec("DELETE FROM products WHERE id = $1", id)
	return err
}

func (pr *productRepositoryImpl) GetProductsByReceptionIDs(receptionIDs []string) ([]models.Product, error) {
	query := "SELECT id, date_time, type, reception_id FROM products WHERE reception_id = ANY($1)"
	rows, err := pr.db.Query(query, pq.Array(receptionIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var products []models.Product
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.ID, &product.DateTime, &product.Type, &product.ReceptionID); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}
