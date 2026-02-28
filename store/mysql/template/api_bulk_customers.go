package main

import (
	"encoding/csv"
	"encoding/json"
	"net/http"
)

func BulkLoadCustomers(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Invalid file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	cr := csv.NewReader(file)
	row, err := cr.Read()

	/*
		TODO:Read each row and insert it into the database.

	*/

	w.WriteHeader(http.StatusCreated)
	rows := 0 // return the number of rows inserted
	json.NewEncoder(w).Encode(rows)
}
