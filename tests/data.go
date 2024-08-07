package tests

import (
	"github.com/dzhordano/avito-bootcamp2024/internal/domain"
	"github.com/google/uuid"
	"time"
)

var (
	house = domain.House{
		ID:        1,
		Address:   "test address 1",
		Year:      2001,
		Developer: "good developer",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	houseId = house.ID // HouseID for flats

	flatCreated = domain.Flat{
		ID:         100,
		FlatNumber: 54,
		Price:      10000,
		Rooms:      3,
		Status:     "created",
	}

	flatApproved = domain.Flat{
		ID:         101,
		FlatNumber: 53,
		Price:      20000,
		Rooms:      4,
		Status:     "approved",
	}

	userModerator = domain.User{
		ID:       uuid.New(),
		Email:    "initTester@mail.ru",
		Password: "qwerty",
		UserType: domain.UserTypeModerator,
	}
)
