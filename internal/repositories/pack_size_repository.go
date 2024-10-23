package repositories

import (
	"database/sql"
	"fmt"

	"github.com/klemis/packs-calculator/models"
)

// PackSizeRepository defines the interface for creating and deleting pack sizes.
type PackSizeRepository interface {
	CreatePackSize(size uint32) error
	DeletePackSize(size uint32) error
	GetPackSizes() ([]models.PackSize, error)
}

// SQLPackSizeRepository is the struct that implements PackSizeRepository interface for SQL database.
type SQLPackSizeRepository struct {
	db *sql.DB
}

// NewSQLPackSizeRepository initializes a new SQL-based repository.
func NewSQLPackSizeRepository(db *sql.DB) PackSizeRepository {
	return &SQLPackSizeRepository{db: db}
}

// GetPackSizes retrieves all pack sizes from the database.
func (r *SQLPackSizeRepository) GetPackSizes() ([]models.PackSize, error) {
	rows, err := r.db.Query(`SELECT id, size FROM pack_sizes ORDER BY size DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var packSizes []models.PackSize
	for rows.Next() {
		var pack models.PackSize
		if err := rows.Scan(&pack.ID, &pack.Size); err != nil {
			return nil, err
		}
		packSizes = append(packSizes, pack)
	}

	return packSizes, nil
}

// CreatePackSize inserts a new pack size into the database.
func (r *SQLPackSizeRepository) CreatePackSize(size uint32) error {
	query := `INSERT INTO pack_sizes (size) VALUES ($1) ON CONFLICT DO NOTHING`

	result, err := r.db.Exec(query, size)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %s", err)
	}

	// Determine if the size was added or already exists.
	if rowsAffected == 0 {
		return fmt.Errorf("no rows were affected")
	}

	return nil
}

// DeletePackSize deletes an existing pack size from the database.
func (r *SQLPackSizeRepository) DeletePackSize(size uint32) error {
	query := `DELETE FROM pack_sizes WHERE size = $1`

	result, err := r.db.Exec(query, size)
	if err != nil {
		return err
	}

	// Check if a row was deleted.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no rows affected, pack size not found")
	}

	return nil
}
