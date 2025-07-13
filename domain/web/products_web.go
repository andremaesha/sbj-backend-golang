package web

import (
	"context"
)

type ProductsRequest struct {
	Id string `json:"id"`
}

type ProductsResponse struct {
	Message  string             `json:"message,omitempty"`
	Products []*ProductResponse `json:"products,omitempty"`
}

type ProductResponse struct {
	Id              string  `json:"id,omitempty"`
	ProductName     string  `json:"product_name,omitempty"`
	ImageUrl        string  `json:"image_url,omitempty"`
	Price           float64 `json:"price,omitempty"`
	Description     string  `json:"description,omitempty"`
	Rating          float64 `json:"rating,omitempty"`
	Category        string  `json:"category,omitempty"`
	Stock           int     `json:"stock,omitempty"`
	NumberOfReviews int     `json:"number_of_reviews,omitempty"`
	ResponseMessage string  `json:"response_message,omitempty"`
}

type ProductsUsecase interface {
	Product(c context.Context, id string) (*ProductResponse, error)
	Products(c context.Context) (*ProductsResponse, error)
	ValidateProductId(id string) error
}
