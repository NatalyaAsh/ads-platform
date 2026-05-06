package model

type CreateAdRequest struct {
	Title       string
	Description string
	Price       float64
	UserID      uint
	CategoryID  uint
}

type UpdateAdRequest struct {
	ID          uint
	Title       string
	Description string
	Price       float64
	CategoryID  uint
}
