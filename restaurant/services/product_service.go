package services

import (
	"context"
	"github.com/cvckeboy/restaurant-app/restaurant/models"
	"github.com/cvckeboy/restaurant-app/restaurant/storage"
	"github.com/cvckeboy/restaurant-app/utils"
	"github.com/google/uuid"
	"time"
)

const (
	defaultTimeout = 5 * time.Second
)

type ProductStore interface {
	GetAllProducts(ctx context.Context) ([]models.Product, error)
	CreateProduct(ctx context.Context, req *models.CreateProductRequest) (uuid.UUID, error)
	UpdateProduct(ctx context.Context, req *models.UpdateProductRequest) error
	DeleteProduct(ctx context.Context, id uuid.UUID) error
	GetProductsByCategory(ctx context.Context, categoryID uuid.UUID) ([]models.Product, error)
	GetProductsSortedByPrice(ctx context.Context, asc bool) ([]models.Product, error)
	CreateCategory(ctx context.Context, req *models.CreateCategoryRequest) (uuid.UUID, error)
	CreateImage(ctx context.Context, req *models.CreateImageRequest) (uuid.UUID, error)
}

type ProductService struct {
	store  *storage.ProductStorage
	logger *utils.Logger
}

func NewProductService(store *storage.ProductStorage, logger *utils.Logger) *ProductService {
	return &ProductService{store: store, logger: logger}
}

func (s *ProductService) CreateCategory(ctx context.Context, req *models.CreateCategoryRequest) (uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	category, err := s.store.CreateCategory(ctx, req)
	if err != nil {
		return uuid.Nil, err
	}
	return category, err
}

func (s *ProductService) CreateImage(ctx context.Context, req *models.CreateImageRequest) (uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	s.logger.Info("Creating image", "name", req.URL)
	imageID, err := s.store.CreateImage(ctx, req)
	if err != nil {
		s.logger.Error("Er")
		return uuid.Nil, err
	}
	return imageID, err
}

func (s *ProductService) GetAllProducts(ctx context.Context) ([]models.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	return s.store.GetAllProducts(ctx)
}

func (s *ProductService) CreateProduct(ctx context.Context, req *models.CreateProductRequest) (uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	id, err := s.store.CreateProduct(ctx, req)
	s.logger.Info("Creating product", "name", req.Name)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (s *ProductService) UpdateProduct(ctx context.Context, req *models.UpdateProductRequest) error {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	return s.store.UpdateProduct(ctx, req)
}

func (s *ProductService) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	return s.store.DeleteProduct(ctx, id)
}

func (s *ProductService) GetProductsByCategory(ctx context.Context, categoryID uuid.UUID) ([]models.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	return s.store.GetProductsByCategory(ctx, categoryID)
}

func (s *ProductService) GetProductsSortedByPrice(ctx context.Context, asc bool) ([]models.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	return s.store.GetProductsSortedByPrice(ctx, asc)
}
