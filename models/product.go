package models

type Product struct {
	ID          int64  `json:"id"`
	Description string `json:"description"`
	Name        string `json:"name"`
	Price       int64  `json:"price"`
	Rating      int    `json:"rating"`
	ImageURL    string `json:"image_url"`
	Stock       int    `json:"stock"`
}
