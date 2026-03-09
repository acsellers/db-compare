package main

import (
	"encoding/json"
	"net/http"
	"time"
)

func CustomerSales(w http.ResponseWriter, r *http.Request) {
	endDate := time.Now().Truncate(24 * time.Hour)
	if r.URL.Query().Get("end_date") != "" {
		endDate, _ = time.Parse("2006-01-02", r.URL.Query().Get("end_date"))
	}
	startDate := endDate.AddDate(0, 0, -7)
	if r.URL.Query().Get("start_date") != "" {
		startDate, _ = time.Parse("2006-01-02", r.URL.Query().Get("start_date"))
	}

	/*
		select customers.id, customers.name,
		sum(orders.total) as total_sales, count(*) as total_orders
		from customers
		join orders on customers.id = orders.customer_id
		where orders.order_date >= ? and orders.order_date <= ?
		group by customers.id, customers.name
		order by total_sales desc;
	*/
	// TODO: run the customer sales query

	if err != nil {
		http.Error(w, "Invalid date", http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(totals)
}
