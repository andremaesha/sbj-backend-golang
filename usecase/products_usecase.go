package usecase

import (
	"errors"
	"sbj-backend/domain"
	"sbj-backend/domain/web"
	"strconv"
	"time"
)

type productsUsecase struct {
	productsRepository domain.ProductsRepository
	contextTimeout     time.Duration
}

func NewProductsUsecase(productsRepository domain.ProductsRepository, contextTimeout time.Duration) web.ProductsUsecase {
	return &productsUsecase{productsRepository: productsRepository, contextTimeout: contextTimeout}
}

func (pu *productsUsecase) Product(id string) (*domain.Product, error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, errors.New("invalid id")
	}

	product, err := pu.productsRepository.GetDataById(idInt)
	if err != nil {
		return nil, err
	}

	return product, nil
}
