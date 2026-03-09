package main

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"net/http"

	"github.com/acsellers/golang-db-compare/store/mysql/bob/models"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/mysql"
	"github.com/stephenafamo/bob/dialect/mysql/im"
)

func BulkLoadCustomers(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Invalid file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	args := []bob.Expression{}
	cr := csv.NewReader(file)
	row, err := cr.Read()
	for err == nil {
		args = append(args, mysql.Arg(row[0], row[1], row[2]))
		row, err = cr.Read()
	}
	if err != io.EOF {
		http.Error(w, "Invalid file", http.StatusBadRequest)
		return
	}

	result, err := mysql.Insert(
		im.Into(
			models.Customers.Name().String(),
			models.Customers.Columns.Name.String(),
			models.Customers.Columns.Email.String(),
			models.Customers.Columns.Phone.String(),
		),
		im.Values(args...),
	).Exec(r.Context(), db)
	if err != nil {
		http.Error(w, "Invalid file", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	rows, _ := result.RowsAffected()
	json.NewEncoder(w).Encode(rows)
}
