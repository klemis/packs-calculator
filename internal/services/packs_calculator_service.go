package services

import (
	"github.com/klemis/packs-calculator/internal/repositories"
)

// PacksCalculator defines the interface for creating and deleting and calculating packs.
type PacksCalculator interface {
	AddPackSize(size uint32) error
	DeletePackSize(size uint32) error
	CalculatePacks(orderQty uint32) (map[uint32]uint32, error)
}

// PacksCalculatorService is an implementation of PacksCalculatorService
type PacksCalculatorService struct {
	repo repositories.PackSizeRepository
}

// NewPacksCalculatorService creates a new instance of PacksCalculator with injected repository.
func NewPacksCalculatorService(packSizeRepo repositories.PackSizeRepository) PacksCalculator {
	return &PacksCalculatorService{
		repo: packSizeRepo,
	}
}

// AddPackSize inserts a new pack size into the database.
func (s *PacksCalculatorService) AddPackSize(size uint32) error {
	return s.repo.CreatePackSize(size)
}

// DeletePackSize removes a pack size from the database by size.
func (s *PacksCalculatorService) DeletePackSize(size uint32) error {
	return s.repo.DeletePackSize(size)
}

// CalculatePacks calculates the optimal pack sizes for a given order quantity.
func (s *PacksCalculatorService) CalculatePacks(orderQty uint32) (map[uint32]uint32, error) {
	// Get packSizes ordered in desc order.
	packSizes, err := s.repo.GetPackSizes()
	if err != nil {
		return nil, err
	}

	result := make(map[uint32]uint32)
	remainingQty := orderQty

	for _, pack := range packSizes {
		if remainingQty >= pack.Size {
			numPacks := remainingQty / pack.Size
			remainingQty = remainingQty % pack.Size
			result[pack.Size] = numPacks
		}
	}

	// If there's still a remaining quantity, use the smallest pack.
	if remainingQty > 0 && len(packSizes) > 0 {
		result[packSizes[len(packSizes)-1].Size]++
	}

	return result, nil
}
