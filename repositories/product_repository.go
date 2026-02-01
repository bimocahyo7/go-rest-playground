package repositories

import (
	"database/sql"
	"errors"
	"go-rest-playground/models"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) GetAll() ([]models.Product, error) {
	query := `
        SELECT p.id, p.name, p.price, p.stock, p.category_id,
               c.id, c.name, c.description
        FROM products p
        LEFT JOIN categories c ON p.category_id = c.id
    `
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		var category models.Category
		var categoryID sql.NullInt64
		var catID sql.NullInt64
		var catName sql.NullString
		var catDesc sql.NullString

		err := rows.Scan(
			&product.ID, &product.Name, &product.Price, &product.Stock, &categoryID,
			&catID, &catName, &catDesc,
		)
		if err != nil {
			return nil, err
		}

		if categoryID.Valid {
			id := int(categoryID.Int64)
			product.CategoryID = &id
		}

		if catID.Valid {
			category.ID = int(catID.Int64)
			category.Name = catName.String
			category.Description = catDesc.String
			product.Category = &category
		}

		products = append(products, product)
	}
	return products, nil
}

func (r *ProductRepository) GetByID(id int) (*models.Product, error) {
	query := `
        SELECT p.id, p.name, p.price, p.stock, p.category_id,
               c.id, c.name, c.description
        FROM products p
        LEFT JOIN categories c ON p.category_id = c.id
        WHERE p.id = $1
    `

	var product models.Product
	var category models.Category
	var categoryID sql.NullInt64
	var catID sql.NullInt64
	var catName sql.NullString
	var catDesc sql.NullString

	err := r.db.QueryRow(query, id).Scan(
		&product.ID, &product.Name, &product.Price, &product.Stock, &categoryID,
		&catID, &catName, &catDesc,
	)
	if err != nil {
		return nil, err
	}

	if categoryID.Valid {
		id := int(categoryID.Int64)
		product.CategoryID = &id
	}

	if catID.Valid {
		category.ID = int(catID.Int64)
		category.Name = catName.String
		category.Description = catDesc.String
		product.Category = &category
	}

	return &product, nil
}

func (r *ProductRepository) Create(product *models.Product) error {
	return r.db.QueryRow(
		"INSERT INTO products (name, price, stock, category_id) VALUES ($1, $2, $3, $4) RETURNING id",
		product.Name, product.Price, product.Stock, product.CategoryID,
	).Scan(&product.ID)
}

func (r *ProductRepository) Update(product *models.Product) error {
	result, err := r.db.Exec(
		"UPDATE products SET name = $1, price = $2, stock = $3, category_id = $4 WHERE id = $5",
		product.Name, product.Price, product.Stock, product.CategoryID, product.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("Product not found")
	}

	return nil
}

func (r *ProductRepository) Delete(id int) error {
	result, err := r.db.Exec("DELETE FROM products WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("Product not found")
	}

	return nil
}
