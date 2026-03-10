package main

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"net/http"

	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/acsellers/golang-db-compare/store/mysql/bob/models"
)

func CustomerUpdate(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Invalid file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	cr := csv.NewReader(file)
	row, err := cr.Read()
	changed := 0
	for err == nil {
		externalId := row[0]
		name := row[1]
		email := row[2]
		phone := row[3]
		changes := &models.CustomerSetter{}
		if name != "" {
			changes.Name = omit.From(name)
		}
		if email != "" {
			changes.Email = omitnull.From(email)
		}
		if phone != "" {
			changes.Phone = omitnull.From(phone)
		}
		_, err := models.Customers.Update(
			models.UpdateWhere.Customers.ExternalID.EQ(externalId),
			changes.UpdateMod(),
		).Exec(r.Context(), db)
		if err != nil {
			http.Error(w, "Invalid file", http.StatusBadRequest)
			return
		}
		changed++
		row, err = cr.Read()
	}
	if err != io.EOF {
		http.Error(w, "Invalid file", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(changed)
}
