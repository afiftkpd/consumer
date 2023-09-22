package usecase

import (
	"consumer/models"
	"context"
)

type ProductUsecase interface {
	Update(ctx context.Context, product models.Product) error
	Delete(ctx context.Context, id string) error
	Store(ctx context.Context, product models.Product) error
}
