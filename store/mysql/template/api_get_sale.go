package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func GetSale(w http.ResponseWriter, r *http.Request) {
	sid := r.PathValue("id")
	id, err := strconv.ParseInt(sid, 10, 64)
	if err != nil {
		fmt.Println(sid, err)
		http.Error(w, "Invalid sale ID", http.StatusBadRequest)
		return
	}

	// TODO: get the order, order items and order payments
	if err != nil {
		fmt.Println("Query: ", err)
		http.Error(w, "Invalid sale ID", http.StatusBadRequest)
		return
	}
	// TODO: output the sale structs from common

	json.NewEncoder(w).Encode(sale)
}
