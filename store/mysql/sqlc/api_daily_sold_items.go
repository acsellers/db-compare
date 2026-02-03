package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/acsellers/golang-db-compare/store/common"
)

func DailySoldItems(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("debug") != "" {
		logQueries = true
		defer func() { logQueries = false }()
	}
	date := time.Now()
	if r.URL.Query().Get("date") != "" {
		date, _ = time.Parse("2006-01-02", r.URL.Query().Get("date"))
	}
	summaries, err := db.DailySoldItems(r.Context(), date)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid date", http.StatusBadRequest)
		return
	}
	summs := []common.ItemSummary{}
	for _, s := range summaries {
		summs = append(summs, common.ItemSummary{
			Name:          s.Name,
			Category:      s.Category,
			TotalQuantity: s.TotalQuantity,
			TotalSales:    s.TotalSales,
			OrderDate:     date.Truncate(24 * time.Hour),
		})
	}
	json.NewEncoder(w).Encode(summs)
}
