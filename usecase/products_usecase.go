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

func (pu *productsUsecase) Product(c context.Context, id string) (*web.ProductsResponse, error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, errors.New("invalid id")
	}

	product, err := pu.productsRepository.GetDataById(c, idInt)
	if err != nil {
		return nil, err
	}

	image := pu.imagesRepository.GetDataById(c, product.ImagesId)

	response := &web.ProductsResponse{
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
