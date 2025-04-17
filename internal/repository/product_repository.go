package repository

import (
	"database/sql"
	"errors"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/ners1us/order-service/internal/model"
)

type ProductRepository interface {
	CreateProduct(product *model.Product) error
	GetLastProductByReceptionID(receptionID string) (*model.Product, error)
	DeleteProduct(id string) error
	GetProductsByReceptionIDs(receptionIDs []string) ([]model.Product, error)
}

type productRepositoryImpl struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepositoryImpl{db}
}

func (pr *productRepositoryImpl) CreateProduct(product *model.Product) error {
	query := "INSERT INTO products (id, date_time, type, reception_id) VALUES ($1, $2, $3, $4)"
	_, err := pr.db.Exec(query, product.ID, product.DateTime, product.Type, product.ReceptionID)
	return err
}

func (pr *productRepositoryImpl) GetLastProductByReceptionID(receptionID string) (*model.Product, error) {
	var product model.Product
	query := "SELECT id, date_time, type, reception_id FROM products WHERE reception_id = $1 ORDER BY date_time DESC LIMIT 1"
	err := pr.db.QueryRow(query, receptionID).Scan(&product.ID, &product.DateTime, &product.Type, &product.ReceptionID)
	if errors.Is(err, sql.ErrNoRows) {
		return &model.Product{}, err
	}
	return &product, nil
}

func (pr *productRepositoryImpl) DeleteProduct(id string) error {
	query := "DELETE FROM products WHERE id = $1"
	_, err := pr.db.Exec(query, id)
	return err
}

func (pr *productRepositoryImpl) GetProductsByReceptionIDs(receptionIDs []string) ([]model.Product, error) {
	query := "SELECT id, date_time, type, reception_id FROM products WHERE reception_id = ANY($1)"
	rows, err := pr.db.Query(query, pq.Array(receptionIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var products []model.Product
	for rows.Next() {
		var product model.Product
		if err := rows.Scan(&product.ID, &product.DateTime, &product.Type, &product.ReceptionID); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}
