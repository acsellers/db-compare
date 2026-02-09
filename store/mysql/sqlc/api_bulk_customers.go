package main

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"io"
	"net/http"

	"github.com/acsellers/golang-db-compare/store/mysql/sqlc/models"
)

func BulkLoadCustomers(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Invalid file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	args := []models.InsertCustomersBulkParams{}
	cr := csv.NewReader(file)
	row, err := cr.Read()
	for err == nil {
		args = append(args, models.InsertCustomersBulkParams{
			Name:  row[0],
			Phone: sql.NullString{String: row[1], Valid: row[1] != ""},
			Email: sql.NullString{String: row[2], Valid: row[2] != ""},
		})
		row, err = cr.Read()
	}
	if err != io.EOF {
		http.Error(w, "Invalid file", http.StatusBadRequest)
		return
	}

	cnt, err := db.InsertCustomersBulk(r.Context(), args)
	if err != nil {
		http.Error(w, "Invalid file", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(cnt)
}
