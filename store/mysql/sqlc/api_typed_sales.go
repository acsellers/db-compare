package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/acsellers/golang-db-compare/store/mysql/sqlc/models"
)

func TypedSales(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("debug") != "" {
		logQueries = true
		defer func() { logQueries = false }()
	}
	endDate := time.Now().Truncate(24 * time.Hour)
	if r.URL.Query().Get("end_date") != "" {
		endDate, _ = time.Parse("2006-01-02", r.URL.Query().Get("end_date"))
	}
	startDate := endDate.AddDate(0, 0, -7)
	if r.URL.Query().Get("start_date") != "" {
		startDate, _ = time.Parse("2006-01-02", r.URL.Query().Get("start_date"))
	}
	/*
			-- typed sales report
		select ro.title, ro.report_order, t.name,
		sum(t.order_count) as order_count,
		sum(t.total_quantity) as quantity,
		sum(t.total_sales) as total_sales
		from item_summaries t
		inner join reporting_order ro on ro.order_type = t.order_type and ro.category = t.category
		group by ro.title, ro.report_order, t.name
		order by ro.report_order, t.name;
	*/
	args := models.WeeklyTypedSalesParams{
		StartDate: startDate,
		EndDate:   endDate,
	}
	lines, err := db.WeeklyTypedSales(r.Context(), args)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid date", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(lines)

}
