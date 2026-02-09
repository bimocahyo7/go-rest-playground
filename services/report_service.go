package services

import (
	"go-rest-playground/models"
	"go-rest-playground/repositories"
	"time"
)

type ReportService struct {
	repo *repositories.ReportRepository
}

func NewReportService(repo *repositories.ReportRepository) *ReportService {
	return &ReportService{repo: repo}
}

// Retrieves sales report for today
func (s *ReportService) GetTodayReport() (*models.SalesReport, error) {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	return s.repo.GetSalesReport(&startOfDay, &endOfDay)
}

// Retrieves sales report for a specific date range
func (s *ReportService) GetDateRangeReport(startDate, endDate time.Time) (*models.SalesReport, error) {
	// Add one day to endDate to include the entire end date
	endDate = endDate.Add(24 * time.Hour)
	return s.repo.GetSalesReport(&startDate, &endDate)
}

// Retrieves all-time sales report (no date filter)
func (s *ReportService) GetAllTimeReport() (*models.SalesReport, error) {
	return s.repo.GetSalesReport(nil, nil)
}
