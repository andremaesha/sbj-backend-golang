package usecase

import (
	"context"
	"database/sql"
	"errors"
	"sbj-backend/bootstrap"
	"sbj-backend/domain"
	"sbj-backend/domain/web"
	"sbj-backend/internal/helpers"
	"strconv"
	"time"
)

type productsUsecase struct {
	productsRepository domain.ProductsRepository
	userRepository     domain.UserRepository
	imagesRepository   domain.ImagesRepository
	contextTimeout     time.Duration
}

func NewProductsUsecase(productsRepository domain.ProductsRepository, userRepository domain.UserRepository, imagesRepository domain.ImagesRepository, contextTimeout time.Duration) web.ProductsUsecase {
	return &productsUsecase{productsRepository: productsRepository, userRepository: userRepository, imagesRepository: imagesRepository, contextTimeout: contextTimeout}
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

	images := pu.imagesRepository.GetDataByProductsId(c, product.Id)

	response := &web.ProductResponse{
		Id:              strconv.Itoa(product.Id),
		ProductName:     product.Name,
		Price:           product.Price,
		Description:     product.Description,
		Rating:          product.Ratings,
		Category:        product.Category,
		Stock:           product.Stock,
		NumberOfReviews: product.NumOfReviews,
		ResponseMessage: "success",
	}

	for _, item := range images {
		response.ImageUrl = append(response.ImageUrl, &web.ImagesUrl{Url: item.Url})
	}

	return response, nil
}

func (pu *productsUsecase) Products(c context.Context) (*web.ProductsResponse, error) {
	datas := pu.productsRepository.Datas(c)

	products := make([]*web.ProductResponse, 0, len(datas))

	for _, item := range datas {
		images := pu.imagesRepository.GetDataByProductsId(c, item.Id)

		product := &web.ProductResponse{
			Id:              strconv.Itoa(item.Id),
			ProductName:     item.Name,
			Price:           item.Price,
			Description:     item.Description,
			Rating:          item.Ratings,
			Category:        item.Category,
			Stock:           item.Stock,
			NumberOfReviews: item.NumOfReviews,
			ImageUrl:        make([]*web.ImagesUrl, 0, len(images)),
		}

		for _, image := range images {
			product.ImageUrl = append(product.ImageUrl, &web.ImagesUrl{Url: image.Url})
		}

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

func (pu *productsUsecase) ProductCreate(c context.Context, env *bootstrap.Env, request *web.ProductRequest) error {
	// Create a new context with timeout
	ctx, cancel := context.WithTimeout(c, pu.contextTimeout)
	defer cancel()

	// Create a new product
	now := time.Now()
	product := &domain.Product{
		Name:        request.Name,
		Price:       request.Price,
		Description: request.Description,
		Category:    request.Category,
		Stock:       request.Stock,
		IsActive:    true,
		CreatedAt:   &now,
	}

	// Create the product
	err := pu.productsRepository.Create(ctx, product)
	if err != nil {
		panic(err)
	}

	image := &domain.Images{
		ProductsId: product.Id,
		Url:        env.CloudinaryDefaultProduct,
		CreatedAt:  &sql.NullTime{Time: now, Valid: true},
	}

	if request.ImageUrl == "" {
		err = pu.imagesRepository.Create(ctx, image)
		if err != nil {
			panic(err)
		}
	}

	return nil
}

func (pu *productsUsecase) ValidatePermission(c context.Context, key, data string) error {
	if data == "" {
		return errors.New("unauthorized access, please login first")
	}

	content, err := helpers.DecryptAES(data, key)
	if err != nil {
		panic(err)
	}

	user, err := pu.userRepository.GetSession(c, content)
	if err != nil {
		return err
	}

	if user.Role != "admin" {
		return errors.New("unauthorized access")
	}

	return nil
}
