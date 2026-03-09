package main

import "net/http"

func JSONReport(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
