package services

import (
	"go-rest-playground/models"
	"go-rest-playground/repositories"
)

type ProductService struct {
	repo *repositories.ProductRepository
}

func NewProductService(repo *repositories.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) GetAllProducts() ([]models.Product, error) {
	return s.repo.GetAll()
}

func (s *ProductService) GetProductByID(id int) (*models.Product, error) {
	return s.repo.GetByID(id)
}

func (s *ProductService) CreateProduct(product *models.Product) error {
	return s.repo.Create(product)
}

func (s *ProductService) UpdateProduct(product *models.Product) error {
	return s.repo.Update(product)
}

func (s *ProductService) DeleteProduct(id int) error {
	return s.repo.Delete(id)
}
