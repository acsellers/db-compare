package main

import "net/http"

func JSONQuery(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
