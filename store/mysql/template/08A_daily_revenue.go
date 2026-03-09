package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func DailyRevenue(w http.ResponseWriter, r *http.Request) {
	endDate := time.Now().Truncate(24 * time.Hour)
	if r.URL.Query().Get("end_date") != "" {
		endDate, _ = time.Parse("2006-01-02", r.URL.Query().Get("end_date"))
	}
	startDate := endDate.AddDate(0, 0, -7)
	if r.URL.Query().Get("start_date") != "" {
		startDate, _ = time.Parse("2006-01-02", r.URL.Query().Get("start_date"))
	}

	/*
		select order_type,sum(total) as total_revenue
		from orders
		where order_date >= ? and order_date <= ?
		group by order_type;
	*/
	// TODO: run the daily revenue query

	if err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid date", http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(revenues)
}
