package main

import (
	"bytes"
	"encoding/csv"
	"io"
	"net/http"

	"github.com/go-sql-driver/mysql"
)

const (
	tempCustomerTable = `CREATE TEMPORARY TABLE temp_customer_update (
		external_id VARCHAR(255), 
		name VARCHAR(255), 
		email VARCHAR(255), 
		phone VARCHAR(255)
	)`
	loadCustomerTable = `LOAD DATA INFILE 'customer_update' 
	INTO TABLE temp_customer_update 
	FIELDS TERMINATED BY '\t' 
	LINES TERMINATED BY '\n'`
	updateCustomerName = `UPDATE customers 
	INNER JOIN temp_customer_update ON customers.external_id = temp_customer_update.external_id 
	SET name = temp_customer_update.name 
	WHERE temp_customer_update.name IS NOT NULL`
	updateCustomerEmail = `UPDATE customers 
	INNER JOIN temp_customer_update ON customers.external_id = temp_customer_update.external_id 
	SET email = temp_customer_update.email 
	WHERE temp_customer_update.email IS NOT NULL`
	updateCustomerPhone = `UPDATE customers 
	INNER JOIN temp_customer_update ON customers.external_id = temp_customer_update.external_id 
	SET phone = temp_customer_update.phone 
	WHERE temp_customer_update.phone IS NOT NULL`
)

func CustomerUpdate2(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Invalid file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	cr := csv.NewReader(file)
	row, err := cr.Read()
	buf := &bytes.Buffer{}
	cw := csv.NewWriter(buf)
	cw.Comma = '\t'
	cw.Write([]string{"external_id", "name", "email", "phone"})

	for err == nil {
		externalId := row[0]
		name := row[1]
		email := row[2]
		phone := row[3]
		cw.Write([]string{externalId, name, email, phone})

		row, err = cr.Read()
	}
	cw.Flush()

	mysql.RegisterReaderHandler("customer_update", func() io.Reader {
		return bytes.NewReader(buf.Bytes())
	})
	db.ExecContext(r.Context(), tempCustomerTable)
	db.ExecContext(r.Context(), loadCustomerTable)
	db.ExecContext(r.Context(), updateCustomerName)
	db.ExecContext(r.Context(), updateCustomerEmail)
	db.ExecContext(r.Context(), updateCustomerPhone)
	db.ExecContext(r.Context(), "DROP TABLE temp_customer_update")
	mysql.DeregisterReaderHandler("customer_update")

}
