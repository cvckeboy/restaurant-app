package models

import (
	"github.com/google/uuid"
)

type Product struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       int       `json:"price"`
	CategoryID  uuid.UUID `json:"category_id"`
	ImageID     uuid.UUID `json:"image_id"`
}

type Category struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type Image struct {
	ID  uuid.UUID `json:"id"`
	URL string    `json:"url"`
}

type CreateCategoryRequest struct {
	Name string `json:"name" binding:"required"`
}

type CreateImageRequest struct {
	URL string `json:"url" binding:"required"`
}

type CreateProductRequest struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       int       `json:"price"`
	CategoryID  uuid.UUID `json:"category_id"`
	ImageID     uuid.UUID `json:"image_id"`
}

type UpdateProductRequest struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       int       `json:"price"`
	CategoryID  uuid.UUID `json:"category_id"`
	ImageID     uuid.UUID `json:"image_id"`
}
