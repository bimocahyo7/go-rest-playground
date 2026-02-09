package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"go-rest-playground/services"
)

type ReportHandler struct {
	service *services.ReportService
}

func NewReportHandler(service *services.ReportService) *ReportHandler {
	return &ReportHandler{service: service}
}

// Handles both /api/report and /api/report/today endpoints
func (h *ReportHandler) HandleReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if the path is /api/report/today
	if strings.HasSuffix(r.URL.Path, "/today") {
		h.GetTodayReport(w, r)
		return
	}

	// Otherwise, handle date range report
	h.GetDateRangeReport(w, r)
}

// Handles GET /api/report/today
func (h *ReportHandler) GetTodayReport(w http.ResponseWriter, r *http.Request) {
	report, err := h.service.GetTodayReport()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}

// Handles GET /api/report?start_date=&end_date=
func (h *ReportHandler) GetDateRangeReport(w http.ResponseWriter, r *http.Request) {
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	// If no date parameters provided, return all-time report
	if startDateStr == "" && endDateStr == "" {
		report, err := h.service.GetAllTimeReport()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(report)
		return
	}

	// Validate that both parameters are provided when one is given
	if startDateStr == "" || endDateStr == "" {
		http.Error(w, "Both start_date and end_date are required when filtering by date", http.StatusBadRequest)
		return
	}

	// Parse dates (expected format: 2026-01-01)
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		http.Error(w, "Invalid start_date format. Use YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		http.Error(w, "Invalid end_date format. Use YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	// Validate date range
	if endDate.Before(startDate) {
		http.Error(w, "end_date must be after start_date", http.StatusBadRequest)
		return
	}

	report, err := h.service.GetDateRangeReport(startDate, endDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}
