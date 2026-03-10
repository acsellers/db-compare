package main

import (
	"net/http"
)

func JSONQuery(w http.ResponseWriter, r *http.Request) {
	cardType := r.URL.Query().Get("card_type")

	w.WriteHeader(http.StatusOK)
}
