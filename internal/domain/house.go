package domain

import "time"

type House struct {
	ID        int
	Address   string
	Year      int
	Developer string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type HouseFlats struct {
	ID      int
	HouseID int
	FlatID  int
}
