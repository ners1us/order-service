package rest

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/ners1us/order-service/internal/models"
	"github.com/ners1us/order-service/internal/repository"
	"github.com/ners1us/order-service/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var container testcontainers.Container
var db *sql.DB

func TestMain(m *testing.M) {
	ctx := context.Background()

	request := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpassword",
			"POSTGRES_DB":       "test-db",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}

	var err error
	container, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: request,
		Started:          true,
	})
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		log.Fatalf("failed to get container host: %s", err)
	}
	mappedPort, err := container.MappedPort(ctx, "5432")
	if err != nil {
		log.Fatalf("failed to get mapped port: %s", err)
	}

	connStr := fmt.Sprintf("host=%s port=%s user=testuser password=testpassword dbname=test-db sslmode=disable",
		host, mappedPort.Port())

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to connect to database: %s", err)
	}

	err = applyMigrations(connStr)
	if err != nil {
		log.Fatalf("failed to apply migrations: %s", err)
	}

	code := m.Run()

	if err := db.Close(); err != nil {
		log.Fatalf("failed to close database connection: %s", err)
	}

	if err := container.Terminate(ctx); err != nil {
		log.Fatalf("failed to terminate container: %s", err)
	}

	os.Exit(code)
}

func applyMigrations(connStr string) error {
	migrationDB, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to open database connection for migrations: %w", err)
	}
	defer migrationDB.Close()

	driver, err := postgres.WithInstance(migrationDB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create postgres driver: %w", err)
	}

	migrationsPath, err := filepath.Abs("../../../migrations")
	if err != nil {
		return fmt.Errorf("failed to get migrations path: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	return nil
}

func TestPVZReceptionProductFlow_Integration(t *testing.T) {
	pvzRepo := repository.NewPVZRepository(db)
	receptionRepo := repository.NewReceptionRepository(db)
	productRepo := repository.NewProductRepository(db)

	pvzService := service.NewPVZService(pvzRepo, receptionRepo, productRepo)
	receptionService := service.NewReceptionService(receptionRepo, pvzRepo)
	productService := service.NewProductService(receptionRepo, productRepo)

	moderatorRole := "moderator"
	employeeRole := "employee"
	pvz := &models.PVZ{
		ID:               uuid.New().String(),
		RegistrationDate: time.Now(),
		City:             "Москва",
	}

	createdPVZ, err := pvzService.CreatePVZ(pvz, moderatorRole)
	if err != nil {
		t.Fatalf("failed to create pvz: %v", err)
	}
	assert.Equal(t, pvz.ID, createdPVZ.ID)
	assert.Equal(t, pvz.City, createdPVZ.City)

	reception, err := receptionService.CreateReception(pvz.ID, employeeRole)
	if err != nil {
		t.Fatalf("failed to create reception: %v", err)
	}
	assert.NotEmpty(t, reception.ID)
	assert.Equal(t, "in_progress", reception.Status)
	assert.Equal(t, pvz.ID, reception.PVZID)

	for i := 0; i < 50; i++ {
		product := &models.Product{
			Type: "электроника",
		}
		createdProduct, err := productService.AddProduct(product, pvz.ID, employeeRole)
		if err != nil {
			t.Fatalf("failed to add product #%d: %v", i+1, err)
		}
		assert.NotEmpty(t, createdProduct.ID, i+1)
		assert.Equal(t, reception.ID, createdProduct.ReceptionID)
		assert.Equal(t, "электроника", createdProduct.Type)
	}

	rows, err := db.Query("SELECT COUNT(*) FROM products WHERE reception_id = $1", reception.ID)
	if err != nil {
		t.Fatalf("failed to fetch products: %v", err)
	}
	defer rows.Close()
	var count int
	if rows.Next() {
		if err := rows.Scan(&count); err != nil {
			t.Fatalf("failed to read products: %v", err)
		}
	}
	assert.Equal(t, 50, count)

	closedReception, err := receptionService.CloseLastReception(pvz.ID, employeeRole)
	if err != nil {
		t.Fatalf("failed to close reception: %v", err)
	}
	assert.Equal(t, reception.ID, closedReception.ID)
	assert.Equal(t, "closed", closedReception.Status)
}
