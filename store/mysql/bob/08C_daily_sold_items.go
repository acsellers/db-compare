package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/acsellers/golang-db-compare/store/common"
	"github.com/acsellers/golang-db-compare/store/mysql/bob/queries"
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
	summaries, err := queries.DailySoldItems(date.Format("2006-01-02")).All(r.Context(), db)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid date", http.StatusBadRequest)
		return
	}
	summs := []common.ItemSummary{}
	for _, s := range summaries {
		tq, _ := strconv.Atoi(s.TotalQuantity)
		ts, _ := strconv.ParseFloat(s.TotalSales, 64)
		summs = append(summs, common.ItemSummary{
			Name:          s.Name,
			Category:      s.Category,
			TotalQuantity: int64(tq),
			TotalSales:    ts,
			OrderDate:     date.Truncate(24 * time.Hour),
		})
	}
	json.NewEncoder(w).Encode(summs)
}
