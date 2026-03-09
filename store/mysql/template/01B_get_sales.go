package main

import "net/http"

func GetSales(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
