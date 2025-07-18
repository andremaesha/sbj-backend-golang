package web

import (
	"context"
	"sbj-backend/bootstrap"
)

type ProductsRequest struct {
	Id string `json:"id" validate:"required"`
}

type ProductRequest struct {
	Name        string  `json:"name" validate:"required,min=3,max=100"`
	Description string  `json:"description" validate:"required,min=10,max=1000"`
	Price       float64 `json:"price" validate:"required,min=0.01"`
	Category    string  `json:"category" validate:"required"`
	Stock       int     `json:"stock" validate:"required,min=0"`
	ImageUrl    string  `json:"image_url,omitempty"`
}

type ProductsResponse struct {
	Message  string             `json:"message,omitempty"`
	Products []*ProductResponse `json:"products,omitempty"`
}

type ProductResponse struct {
	Id              string       `json:"id,omitempty"`
	ProductName     string       `json:"product_name,omitempty"`
	ImageUrl        []*ImagesUrl `json:"images_url,omitempty"`
	Price           float64      `json:"price,omitempty"`
	Description     string       `json:"description,omitempty"`
	Rating          float64      `json:"rating,omitempty"`
	Category        string       `json:"category,omitempty"`
	Stock           int          `json:"stock,omitempty"`
	NumberOfReviews int          `json:"number_of_reviews,omitempty"`
	ResponseMessage string       `json:"response_message,omitempty"`
}

type ImagesUrl struct {
	Url string `json:"url"`
}

type ProductsUsecase interface {
	Product(c context.Context, id string) (*ProductResponse, error)
	Products(c context.Context) (*ProductsResponse, error)
	ProductCreate(c context.Context, env *bootstrap.Env, request *ProductRequest) error
	ValidateProductId(id string) error
	ValidatePermission(c context.Context, key, data string) error
}
