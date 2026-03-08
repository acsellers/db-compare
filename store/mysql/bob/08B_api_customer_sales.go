package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/acsellers/golang-db-compare/store/common"
	"github.com/acsellers/golang-db-compare/store/mysql/bob/models"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/mysql"
	"github.com/stephenafamo/bob/dialect/mysql/sm"
	"github.com/stephenafamo/scan"
)

func CustomerSales(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("debug") != "" {
		fmt.Println("debug")
		logQueries = true
		defer func() {
			logQueries = false
		}()
	}

	endDate := time.Now().Truncate(24 * time.Hour)
	if r.URL.Query().Get("end_date") != "" {
		endDate, _ = time.Parse("2006-01-02", r.URL.Query().Get("end_date"))
	}
	startDate := endDate.AddDate(0, 0, -7)
	if r.URL.Query().Get("start_date") != "" {
		startDate, _ = time.Parse("2006-01-02", r.URL.Query().Get("start_date"))
	}

	query := models.Customers.Query(
		sm.Columns(
			models.Customers.Columns.ID,
			models.Customers.Columns.Name,
			mysql.Raw("SUM(`orders`.`total`) as total_sales"),
			mysql.Raw("COUNT(`orders`.`id`) as total_orders"),
		),
		models.SelectJoins.Customers.InnerJoin.Orders,
		models.SelectWhere.Orders.OrderDate.GTE(startDate),
		models.SelectWhere.Orders.OrderDate.LTE(endDate),
		sm.GroupBy("customers.id, customers.name"),
	)
	totals, err := bob.All(
		r.Context(),
		db,
		query,
		scan.StructMapper[common.CustomerTotals](),
	)

	if err != nil {
		http.Error(w, "Invalid date", http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(totals)
}
