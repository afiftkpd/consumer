package es

import (
	"consumer/models"
	"context"
	"encoding/json"
	"strconv"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/update"
)

type productRepository struct {
	DB *elasticsearch.TypedClient
}

func NewProductRepository(db *elasticsearch.TypedClient) ProductRepository {
	return &productRepository{db}
}

func (p *productRepository) Update(ctx context.Context, product models.Product) error {
	b, err := json.Marshal(product)
	if err != nil {
		return err
	}

	_, err = p.DB.Update("products", strconv.Itoa(int(product.ID))).
		Request(&update.Request{
			Doc: b,
		}).Do(ctx)
	return err
}

func (p *productRepository) Delete(ctx context.Context, id string) error {
	_, err := p.DB.Delete("products", id).Do(ctx)
	return err
}

func (p *productRepository) Store(ctx context.Context, product models.Product) error {
	_, err := p.DB.Index("products").Id(strconv.Itoa(int(product.ID))).Request(product).Do(ctx)
	return err
}
