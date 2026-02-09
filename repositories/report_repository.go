package repositories

import (
	"database/sql"
	"go-rest-playground/models"
	"time"
)

type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

// Retrieves sales summary for a specific date range (or all-time if dates are nil)
func (repo *ReportRepository) GetSalesReport(startDate, endDate *time.Time) (*models.SalesReport, error) {
	report := &models.SalesReport{}

	// Build query based on whether date filter is provided
	var query string
	var args []interface{}

	if startDate != nil && endDate != nil {
		// With date filter
		query = `
			SELECT 
				COALESCE(SUM(total_amount), 0) as total_revenue,
				COUNT(*) as total_transactions
			FROM transactions
			WHERE created_at >= $1 AND created_at < $2
		`
		args = []interface{}{startDate, endDate}
	} else {
		// All-time (no date filter)
		query = `
			SELECT 
				COALESCE(SUM(total_amount), 0) as total_revenue,
				COUNT(*) as total_transactions
			FROM transactions
		`
		args = []interface{}{}
	}

	err := repo.db.QueryRow(query, args...).Scan(&report.TotalRevenue, &report.TotalTransactions)
	if err != nil {
		return nil, err
	}

	// Get best selling product
	var bestProductQuery string
	if startDate != nil && endDate != nil {
		// With date filter
		bestProductQuery = `
			SELECT 
				p.name,
				SUM(td.quantity) as quantity_sold
			FROM transaction_details td
			JOIN transactions t ON td.transaction_id = t.id
			JOIN products p ON td.product_id = p.id
			WHERE t.created_at >= $1 AND t.created_at < $2
			GROUP BY p.id, p.name
			ORDER BY quantity_sold DESC
			LIMIT 1
		`
	} else {
		// All-time (no date filter)
		bestProductQuery = `
			SELECT 
				p.name,
				SUM(td.quantity) as quantity_sold
			FROM transaction_details td
			JOIN transactions t ON td.transaction_id = t.id
			JOIN products p ON td.product_id = p.id
			GROUP BY p.id, p.name
			ORDER BY quantity_sold DESC
			LIMIT 1
		`
	}

	var productName string
	var quantitySold int
	err = repo.db.QueryRow(bestProductQuery, args...).Scan(&productName, &quantitySold)
	if err == sql.ErrNoRows {
		// No transactions in this period, leave best_selling_product as nil
		return report, nil
	}
	if err != nil {
		return nil, err
	}

	report.BestSellingProduct = &models.BestSellingProduct{
		Name:         productName,
		QuantitySold: quantitySold,
	}

	return report, nil
}
