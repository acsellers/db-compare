package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func DailySoldItems(w http.ResponseWriter, r *http.Request) {
	date := time.Now()
	if r.URL.Query().Get("date") != "" {
		date, _ = time.Parse("2006-01-02", r.URL.Query().Get("date"))
	}
	/*
		select id, name, category,
		cast(sum(total_quantity) as signed) as total_quantity,
		cast(sum(total_sales) as double) as total_sales
		from item_summaries
		where order_date = ?
		group by id, name, category
		order by category, name;
	*/
	// TODO: run the daily sold items query
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid date", http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(summs)
}
