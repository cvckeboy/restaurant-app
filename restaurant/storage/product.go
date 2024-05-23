package storage

import (
	"context"
	"errors"
	"github.com/cvckeboy/restaurant-app/restaurant/models"
	"github.com/cvckeboy/restaurant-app/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductStorage struct {
	pool   *pgxpool.Pool
	logger *utils.Logger
}

func NewProductStorage(pool *pgxpool.Pool, logger *utils.Logger) *ProductStorage {
	return &ProductStorage{pool: pool, logger: logger}
}

func (s *ProductStorage) CreateCategory(ctx context.Context, req *models.CreateCategoryRequest) (uuid.UUID, error) {
	s.logger.Info("Inserting category into database", "name", req.Name)

	var id uuid.UUID
	query := `INSERT INTO restaurant.public.categories (name) VALUES ($1) RETURNING id`
	err := s.pool.QueryRow(ctx, query, req.Name).Scan(&id)
	if err != nil {
		s.logger.Error("Error inserting category", "error", err)
		return uuid.Nil, err
	}

	s.logger.Info("Category inserted into database", "id", id)
	return id, nil
}

func (s *ProductStorage) CreateImage(ctx context.Context, req *models.CreateImageRequest) (uuid.UUID, error) {
	s.logger.Info("Inserting image into database", "url", req.URL)

	var id uuid.UUID
	query := `INSERT INTO images (url) VALUES ($1) RETURNING id`
	err := s.pool.QueryRow(ctx, query, req.URL).Scan(&id)
	if err != nil {
		s.logger.Error("Error inserting image", "error", err)
		return uuid.Nil, err
	}

	s.logger.Info("Image inserted into database", "id", id)
	return id, nil
}

func (s *ProductStorage) CreateProduct(ctx context.Context, req *models.CreateProductRequest) (uuid.UUID, error) {
	s.logger.Info("Inserting req into database", "name", req.Name)

	// Проверяем, существует ли category_id
	var categoryExists bool
	err := s.pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM restaurant.public.categories WHERE id=$1)", req.CategoryID).Scan(&categoryExists)
	if err != nil {
		s.logger.Error("Error checking category existence", "error", err)
		return uuid.Nil, err
	}
	if !categoryExists {
		s.logger.Error("Category does not exist", "category_id", req.CategoryID)
		return uuid.Nil, errors.New("category does not exist")
	}

	// Проверяем, существует ли image_id
	var imageExists bool
	err = s.pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM restaurant.public.images WHERE id=$1)", req.ImageID).Scan(&imageExists)
	if err != nil {
		s.logger.Error("Error checking image existence", "error", err)
		return uuid.Nil, err
	}
	if !imageExists {
		s.logger.Error("Image does not exist", "image_id", req.ImageID)
		return uuid.Nil, errors.New("image does not exist")
	}

	// Вставляем продукт
	var id uuid.UUID
	query := `
        INSERT INTO restaurant.public.products (name, description, price, category_id, image_id)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `
	err = s.pool.QueryRow(ctx, query, req.Name, req.Description, req.Price, req.CategoryID, req.ImageID).Scan(&id)
	if err != nil {
		s.logger.Error("Error inserting req", "error", err)
		return uuid.Nil, err
	}

	s.logger.Info("Product inserted into database", "id", id)
	return id, nil
}

func (s *ProductStorage) GetAllProducts(ctx context.Context) ([]models.Product, error) {
	rows, err := s.pool.Query(ctx, "SELECT ID, Name, Description, Price, Image_id, Category_id FROM restaurant.public.products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var req models.Product
		if err := rows.Scan(&req.ID, &req.Name, &req.Description, &req.Price, &req.ImageID, &req.CategoryID); err != nil {
			return nil, err
		}
		products = append(products, req)
	}

	return products, nil
}

func (s *ProductStorage) UpdateProduct(ctx context.Context, req *models.UpdateProductRequest) error {
	_, err := s.pool.Exec(ctx, "UPDATE restaurant.public.products SET name = ?, description = ?, price = ?, image_id = ?, category_id = ? WHERE id = ?",
		req.Name, req.Description, req.Price, req.ImageID, req.CategoryID, req.ID)
	return err
}

func (s *ProductStorage) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	_, err := s.pool.Exec(ctx, "DELETE FROM products WHERE id = ?", id)
	return err
}

func (s *ProductStorage) GetProductsByCategory(ctx context.Context, categoryID uuid.UUID) ([]models.Product, error) {
	rows, err := s.pool.Query(ctx, "SELECT id, name, description, price, image_id, category_id FROM products WHERE category_id = ?", categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var req models.Product
		if err := rows.Scan(&req.ID, &req.Name, &req.Description, &req.Price, &req.ImageID, &req.CategoryID); err != nil {
			return nil, err
		}
		products = append(products, req)
	}

	return products, nil
}

func (s *ProductStorage) GetProductsSortedByPrice(ctx context.Context, asc bool) ([]models.Product, error) {
	var query string
	if asc {
		query = "SELECT id, name, description, price, image_id, category_id FROM products ORDER BY price ASC"
	} else {
		query = "SELECT id, name, description, price, image_id, category_id FROM products ORDER BY price DESC"
	}

	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var req models.Product
		if err := rows.Scan(&req.ID, &req.Name, &req.Description, &req.Price, &req.ImageID, &req.CategoryID); err != nil {
			return nil, err
		}
		products = append(products, req)
	}

	return products, nil
}
