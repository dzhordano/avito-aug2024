package dtos

import (
	"errors"
	"github.com/dzhordano/avito-bootcamp2024/internal/domain"
)

type FlatCreateInput struct {
	FlatNumber int `json:"flat_number" binding:"required"` // same there.
	HouseId    int `json:"house_id" binding:"required"`
	Price      int `json:"price" binding:"required"`
	Rooms      int `json:"rooms" binding:"required"`
}

type FlatUpdateInput struct {
	FlatId int           `json:"flat_id" binding:"required"`
	Status domain.Status `json:"status" binding:"required"`
}

func (s *FlatUpdateInput) Validate() error {
	if s.FlatId <= 0 {
		return errors.New("invalid flat_id")
	}

	if !s.Status.Validate() {
		return errors.New("invalid status")
	}

	return nil
}
