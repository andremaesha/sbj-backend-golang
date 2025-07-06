package web

import "sbj-backend/domain"

type ProductsRequest struct {
	Id string `json:"id"`
}

type ProductsResponse struct {
	Id              string `json:"id,omitempty"`
	ProductName     string `json:"product_name,omitempty"`
	ImageUrl        string `json:"image_url,omitempty"`
	Price           string `json:"price,omitempty"`
	Description     string `json:"description,omitempty"`
	Rating          string `json:"rating,omitempty"`
	Category        string `json:"category,omitempty"`
	Stock           string `json:"stock,omitempty"`
	NumberOfReviews string `json:"number_of_reviews,omitempty"`
}

type ProductsUsecase interface {
	Product(id string) (*domain.Product, error)
}
