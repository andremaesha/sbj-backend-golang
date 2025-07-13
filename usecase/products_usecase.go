package usecase

import (
	"context"
	"errors"
	"sbj-backend/domain"
	"sbj-backend/domain/web"
	"strconv"
	"time"
)

type productsUsecase struct {
	productsRepository domain.ProductsRepository
	imagesRepository   domain.ImagesRepository
	contextTimeout     time.Duration
}

func NewProductsUsecase(productsRepository domain.ProductsRepository, imagesRepository domain.ImagesRepository, contextTimeout time.Duration) web.ProductsUsecase {
	return &productsUsecase{productsRepository: productsRepository, imagesRepository: imagesRepository, contextTimeout: contextTimeout}
}

func (pu *productsUsecase) Product(c context.Context, id string) (*web.ProductResponse, error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, errors.New("invalid id")
	}

	product, err := pu.productsRepository.GetDataById(c, idInt)
	if err != nil {
		return nil, err
	}

	image := pu.imagesRepository.GetDataById(c, product.ImagesId)

	response := &web.ProductResponse{
		Id:              strconv.Itoa(product.Id),
		ProductName:     product.Name,
		ImageUrl:        image.Url,
		Price:           product.Price,
		Description:     product.Description,
		Rating:          product.Ratings,
		Category:        product.Category,
		Stock:           product.Stock,
		NumberOfReviews: product.NumOfReviews,
		ResponseMessage: "success",
	}

	return response, nil
}

func (pu *productsUsecase) Products(c context.Context) (*web.ProductsResponse, error) {
	datas := pu.productsRepository.Datas(c)

	var products []*web.ProductResponse
	product := new(web.ProductResponse)

	for _, item := range datas {
		image := pu.imagesRepository.GetDataById(c, item.ImagesId)

		product.Id = strconv.Itoa(item.Id)
		product.ProductName = item.Name
		product.ImageUrl = image.Url
		product.Price = item.Price
		product.Description = item.Description
		product.Rating = item.Ratings
		product.Category = item.Category
		product.Stock = item.Stock
		product.NumberOfReviews = item.NumOfReviews

		products = append(products, product)
	}

	return &web.ProductsResponse{
		Message:  "success",
		Products: products,
	}, nil
}

func (pu *productsUsecase) ValidateProductId(id string) error {
	if id == "" {
		return errors.New("id parameter is required")
	}
	return nil
}
