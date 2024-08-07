package dtos

type HouseCreateInput struct {
	Address   string `json:"address" binding:"required"`
	Year      int    `json:"year" binding:"required"`
	Developer string `json:"developer,omitempty"`
}

type HouseSubscribeInput struct {
	Email string `json:"email" binding:"required"`
}
