package catalog

import (
	"context"
	"encoding/json"
	"errors"
	"gopkg.in/olivere/elastic.v5"
	"log"
)

var (
	ErrNotFound = errors.New("entity not found")
)

type Repository interface {
	Close()
	PutProduct(ctx context.Context, product Product) error
	GetProductByID(ctx context.Context, productID string) (Product, error)
	ListProducts(ctx context.Context, skip uint64, take uint64) (*[]Product, error)
	ListProductsWithIDs(ctx context.Context, productIDs []string) (*[]Product, error)
	SearchProducts(ctx context.Context, query string, skip uint64, take uint64) (*[]Product, error)
}

type productDocument struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type elasticRepository struct {
	client *elastic.Client
}

func (r *elasticRepository) Close() {
	r.client.Stop()
}

func (r *elasticRepository) PutProduct(ctx context.Context, product Product) error {
	_, err := r.client.Index().Index("catalog").Type("product").
		Id(product.Id).
		BodyJson(productDocument{
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price}).
		Do(ctx)
	return err
}

func (r *elasticRepository) GetProductByID(ctx context.Context, productID string) (Product, error) {
	result, err := r.client.Get().Index("catalog").Type("product").Id(productID).Do(ctx)
	if err != nil {
		return Product{}, err
	}
	if !result.Found {
		return Product{}, ErrNotFound
	}
	p := productDocument{}
	if err = json.Unmarshal(*result.Source, &p); err != nil {
		return Product{}, err
	}
	return Product{
		Id:          productID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
	}, nil
}

func (r *elasticRepository) ListProducts(ctx context.Context, skip uint64, take uint64) (*[]Product, error) {
	result, err := r.client.Search().
		Index("catalog").Type("product").
		Query(elastic.NewMatchAllQuery()).From(int(skip)).Size(int(take)).Do(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	products := []Product{}
	for _, hit := range result.Hits.Hits {
		p := productDocument{}
		if err = json.Unmarshal(*hit.Source, &p); err == nil {
			products = append(products, Product{
				Id:          hit.Id,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
			})
		}
	}
	return &products, err
}

func (r *elasticRepository) ListProductsWithIDs(ctx context.Context, productIDs []string) (*[]Product, error) {
	items := []*elastic.MultiGetItem{}
	for _, id := range productIDs {
		items = append(
			items,
			elastic.NewMultiGetItem().Index("catalog").Type("product").Id(id),
		)
	}
	res, err := r.client.MultiGet().
		Add(items...).Do(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var products []Product
	for _, doc := range res.Docs {
		p := productDocument{}
		if err = json.Unmarshal(*doc.Source, &p); err == nil {
			products = append(products, Product{
				Id:          doc.Id,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
			})
		}
	}
	return &products, err
}

func (r *elasticRepository) SearchProducts(ctx context.Context, query string, skip uint64, take uint64) (*[]Product, error) {
	result, err := r.client.Search().
		Index("catalog").Type("product").
		Query(elastic.NewMultiMatchQuery(query, "name", "description", "price")).
		From(int(skip)).Size(int(take)).Do(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var products []Product
	for _, hit := range result.Hits.Hits {
		p := productDocument{}
		if err = json.Unmarshal(*hit.Source, &p); err == nil {
			products = append(products, Product{
				Id:          hit.Id,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
			})
		}
	}
	return &products, err
}

func NewElasticRepository(url string) (Repository, error) {
	client, err := elastic.NewClient(
		elastic.SetURL(url),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false),
	)
	if err != nil {
		return nil, err
	}
	return &elasticRepository{client}, nil
}
