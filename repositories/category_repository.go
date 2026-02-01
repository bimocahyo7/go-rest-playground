package repositories

import (
	"database/sql"
	"errors"
	"go-rest-playground/models"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) GetAll() ([]models.Category, error) {
	rows, err := r.db.Query("SELECT id, name, description FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var category models.Category
		if err := rows.Scan(&category.ID, &category.Name, &category.Description); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func (r *CategoryRepository) GetByID(id int) (*models.Category, error) {
	var category models.Category
	err := r.db.QueryRow("SELECT id, name, description FROM categories WHERE id = $1", id).
		Scan(&category.ID, &category.Name, &category.Description)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *CategoryRepository) Create(category *models.Category) error {
	return r.db.QueryRow(
		"INSERT INTO categories (name, description) VALUES ($1, $2) RETURNING id",
		category.Name, category.Description,
	).Scan(&category.ID)
}

func (r *CategoryRepository) Update(category *models.Category) error {
	result, err := r.db.Exec(
		"UPDATE categories SET name = $1, description = $2 WHERE id = $3",
		category.Name, category.Description, category.ID,
	)
	if err != nil {
		return err
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("Category not found")
	}

	return nil
}

func (r *CategoryRepository) Delete(id int) error {
	result, err := r.db.Exec("DELETE FROM categories WHERE id = $1", id)
	if err != nil {
		return err
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("Category not found")
	}

	return nil
}
