package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/acsellers/golang-db-compare/store/mysql/sqlc/models"
)

func CustomerSales(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("debug") != "" {
		fmt.Println("debug")
		logQueries = true
		defer func() {
			logQueries = false
		}()
	}

	endDate := time.Now().Truncate(24 * time.Hour)
	if r.URL.Query().Get("end_date") != "" {
		endDate, _ = time.Parse("2006-01-02", r.URL.Query().Get("end_date"))
	}
	startDate := endDate.AddDate(0, 0, -7)
	if r.URL.Query().Get("start_date") != "" {
		startDate, _ = time.Parse("2006-01-02", r.URL.Query().Get("start_date"))
	}

	args := models.CustomerSalesParams{
		StartDate: startDate,
		EndDate:   endDate,
	}
	totals, err := db.CustomerSales(r.Context(), args)

	if err != nil {
		http.Error(w, "Invalid date", http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(totals)
}
