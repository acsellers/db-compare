package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/acsellers/golang-db-compare/store/common"
	"github.com/acsellers/golang-db-compare/store/mysql/bob/models"
	"github.com/stephenafamo/bob/dialect/mysql"
	"github.com/stephenafamo/bob/dialect/mysql/sm"
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
	summaries, err := models.ItemSummaries.Query(
		models.SelectWhere.ItemSummaries.OrderDate.EQ(date),
		sm.Columns(
			models.ItemSummaries.Columns.Name,
			models.ItemSummaries.Columns.Category,
			mysql.Raw("SUM(`item_summaries`.`total_quantity`) as total_quantity"),
			mysql.Raw("SUM(`item_summaries`.`total_sales`) as total_sales"),
		),
		sm.GroupBy("item_summaries.name, item_summaries.category"),
	).All(r.Context(), db)
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
			TotalQuantity: s.TotalQuantity.GetOrZero().IntPart(),
			TotalSales:    s.TotalSales.GetOrZero().InexactFloat64(),
			OrderDate:     date.Truncate(24 * time.Hour),
		})
	}
	json.NewEncoder(w).Encode(summs)
}
