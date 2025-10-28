package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"sbj-backend/bootstrap"
	"sbj-backend/domain"
	"sbj-backend/domain/web"
	"sbj-backend/internal/curl"
	"sbj-backend/internal/helpers"
	"strconv"
	"time"

	"github.com/google/uuid"
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

func (pu *productsUsecase) ProductCreate(c context.Context, sessionId string, env *bootstrap.Env, request *web.ProductRequest) error {
	session, err := helpers.DecryptAES(sessionId, env.Key)
	if err != nil {
		panic(err)
	}

	user, err := pu.userRepository.GetSession(c, session)
	if err != nil {
		return err
	}

	_, err = pu.userRepository.GetByEmail(c, user.Email)
	if err != nil {
		return err
	}

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
		CreatedBy:   user.Id,
		IsActive:    true,
		CreatedAt:   &now,
	}

	// Create the product
	err = pu.productsRepository.Create(ctx, product)
	if err != nil {
		panic(err)
	}

	if len(request.Images) == 0 {
		image := &domain.Images{
			ProductsId: product.Id,
			Url:        env.CloudinaryDefaultProduct,
			CreatedBy:  user.Id,
			CreatedAt:  &sql.NullTime{Time: now, Valid: true},
		}

		err = pu.imagesRepository.Create(ctx, image)
		if err != nil {
			panic(err)
		}
	} else {
		for _, item := range request.Images {
			image := &domain.Images{
				ProductsId: product.Id,
				AssetId:    item.AssetId,
				PublicId:   item.PublicId,
				Url:        item.Url,
				CreatedBy:  user.Id,
				CreatedAt:  &sql.NullTime{Time: now, Valid: true},
			}

			err = pu.imagesRepository.Create(ctx, image)
			if err != nil {
				panic(err)
			}
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

func (pu *productsUsecase) UploadImages(env *bootstrap.Env, fileHeader []*multipart.FileHeader) (*web.ProductsImagesResponse, error) {
	var images []*web.ImagesUrl

	timeStamp := fmt.Sprintf("%d", time.Now().Unix())
	folder := "products"

	for i, file := range fileHeader {
		// Generate unique public ID untuk setiap file
		publicId := fmt.Sprintf("%s_%d", uuid.New().String(), i)

		// Generate signature untuk Cloudinary
		formula := fmt.Sprintf("folder=%s&public_id=%s&timestamp=%s%s", folder, publicId, timeStamp, env.CloudinaryApiSecret)
		signature := helpers.GenerateSH1(formula)

		// Save file sementara
		fileUpload, err := helpers.SaveTempFile(file, "uploads", "img_", "_backup")
		if err != nil {
			return nil, fmt.Errorf("failed to save temp file: %v", err)
		}

		url := env.CloudinaryUrl
		cloudinaryResponse := new(domain.ResponseCloudinary)

		err = curl.Curl[*domain.ResponseCloudinary]("POST", url, nil, nil, map[string]string{
			"file":      fileUpload,
			"api_key":   env.CloudinaryApiKey,
			"timestamp": timeStamp,
			"signature": signature,
			"public_id": publicId,
			"folder":    folder,
		}, cloudinaryResponse)

		if removeErr := os.Remove(fileUpload); removeErr != nil {
			panic(fmt.Sprintf("Warning: Failed to delete temp file %s: %v\n", fileUpload, removeErr))
		}

		if err != nil {
			return nil, fmt.Errorf("failed to upload file %d to cloudinary: %v", i+1, err)
		}

		images = append(images, &web.ImagesUrl{
			AssetId:  cloudinaryResponse.AssetId,
			PublicId: cloudinaryResponse.PublicId,
			Url:      cloudinaryResponse.SecureUrl,
		})
	}

	return &web.ProductsImagesResponse{
		ImagesUrl:       images,
		ResponseMessage: "success",
	}, nil
}
