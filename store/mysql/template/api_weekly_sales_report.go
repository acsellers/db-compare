package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func TypedSales(w http.ResponseWriter, r *http.Request) {
	endDate := time.Now().Truncate(24 * time.Hour)
	if r.URL.Query().Get("end_date") != "" {
		endDate, _ = time.Parse("2006-01-02", r.URL.Query().Get("end_date"))
	}
	startDate := endDate.AddDate(0, 0, -7)
	if r.URL.Query().Get("start_date") != "" {
		startDate, _ = time.Parse("2006-01-02", r.URL.Query().Get("start_date"))
	}

	/*
		select dim_date.year, dim_date.WEEK_OF_YEAR, ro.title, ro.report_order, t.name,
		sum(t.order_count) as order_count,
		sum(t.total_quantity) as quantity,
		sum(t.total_sales) as total_sales
		from item_summaries t
		inner join reporting_order ro on ro.order_type = t.order_type and ro.category = t.category
		inner join dim_date on dim_date.date = t.order_date
		where t.order_date >= ? and t.order_date <= ?
		group by dim_date.year, dim_date.WEEK_OF_YEAR,ro.title, ro.report_order, t.name
		order by dim_date.year, dim_date.WEEK_OF_YEAR, ro.report_order, t.name;
	*/
	// TODO: run the typed sales query
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid date", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(results)

}
