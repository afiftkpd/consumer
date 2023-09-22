package usecase

import (
	"consumer/models"
	"consumer/repository/es"
	"context"
)

type productUsecase struct {
	ElasticRepo es.ProductRepository
}

func NewProductUsecase(esRepo es.ProductRepository) ProductUsecase {
	return &productUsecase{esRepo}
}

func (p *productUsecase) Update(ctx context.Context, product models.Product) error {
	return p.ElasticRepo.Update(ctx, product)
}

func (p *productUsecase) Delete(ctx context.Context, id string) error {
	return p.ElasticRepo.Delete(ctx, id)
}

func (p *productUsecase) Store(ctx context.Context, product models.Product) error {
	return p.ElasticRepo.Store(ctx, product)
}
