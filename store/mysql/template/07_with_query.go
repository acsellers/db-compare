package main

import "net/http"

func WithQuery(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
