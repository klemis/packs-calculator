package repositories

import (
	_ "database/sql"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/klemis/packs-calculator/models"
	"github.com/stretchr/testify/assert"
)

func TestCreatePackSize(t *testing.T) {
	tests := []struct {
		name        string
		size        uint32
		mockSetup   func(mock sqlmock.Sqlmock)
		expectedErr error
	}{
		{
			name: "successful insert",
			size: 100,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO pack_sizes (size) VALUES ($1) ON CONFLICT DO NOTHING`)).
					WithArgs(100).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedErr: nil,
		},
		{
			name: "insert conflict (no rows affected)",
			size: 100,
			mockSetup: func(mock sqlmock.Sqlmock) {
				// Expect conflict handling, no rows affected.
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO pack_sizes (size) VALUES ($1) ON CONFLICT DO NOTHING`)).
					WithArgs(100).
					WillReturnResult(sqlmock.NewResult(1, 0))
			},
			expectedErr: fmt.Errorf("no rows were affected"),
		},
		{
			name: "insert failure",
			size: 100,
			mockSetup: func(mock sqlmock.Sqlmock) {
				// Simulate insert error.
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO pack_sizes (size) VALUES ($1) ON CONFLICT DO NOTHING`)).
					WithArgs(100).
					WillReturnError(errors.New("insert failed"))
			},
			expectedErr: errors.New("insert failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			repo := NewSQLPackSizeRepository(db)
			tt.mockSetup(mock)

			err = repo.CreatePackSize(tt.size)
			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDeletePackSize(t *testing.T) {
	tests := []struct {
		name        string
		size        uint32
		mockSetup   func(mock sqlmock.Sqlmock)
		expectedErr error
	}{
		{
			name: "successful delete",
			size: 100,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM pack_sizes WHERE size = $1`)).
					WithArgs(100).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedErr: nil,
		},
		{
			name: "delete non-existing size",
			size: 100,
			mockSetup: func(mock sqlmock.Sqlmock) {
				// Simulate deletion for a non-existing size.
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM pack_sizes WHERE size = $1`)).
					WithArgs(100).
					WillReturnResult(sqlmock.NewResult(1, 0))
			},
			expectedErr: fmt.Errorf("no rows affected, pack size not found"),
		},
		{
			name: "delete failure",
			size: 100,
			mockSetup: func(mock sqlmock.Sqlmock) {
				// Simulate deletion error.
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM pack_sizes WHERE size = $1`)).
					WithArgs(100).
					WillReturnError(errors.New("delete failed"))
			},
			expectedErr: errors.New("delete failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			repo := NewSQLPackSizeRepository(db)
			tt.mockSetup(mock)

			err = repo.DeletePackSize(tt.size)
			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetPackSizes(t *testing.T) {
	tests := []struct {
		name          string
		mockSetup     func(mock sqlmock.Sqlmock)
		expected      []models.PackSize
		expectedError string
	}{
		{
			name: "successful retrieval",
			mockSetup: func(mock sqlmock.Sqlmock) {
				// Expect query to return pack sizes.
				rows := sqlmock.NewRows([]string{"id", "size"}).
					AddRow(1, 250).
					AddRow(2, 500)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, size FROM pack_sizes ORDER BY size DESC`)).
					WillReturnRows(rows)
			},
			expected: []models.PackSize{
				{ID: 1, Size: 250},
				{ID: 2, Size: 500},
			},
			expectedError: "",
		},
		{
			name: "query error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				// Simulate query failure.
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, size FROM pack_sizes ORDER BY size DESC`)).
					WillReturnError(errors.New("query error"))
			},
			expected:      nil,
			expectedError: "query error",
		},
		{
			name: "scan error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				// Simulate scan failure with an invalid data type for size.
				rows := sqlmock.NewRows([]string{"id", "size"}).
					AddRow(1, "invalid_data") // Simulate scan error due to type mismatch.

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, size FROM pack_sizes ORDER BY size DESC`)).
					WillReturnRows(rows)
			},
			expected:      nil,
			expectedError: "sql: Scan error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			repo := NewSQLPackSizeRepository(db)
			tt.mockSetup(mock)

			packSizes, err := repo.GetPackSizes()
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, packSizes)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
