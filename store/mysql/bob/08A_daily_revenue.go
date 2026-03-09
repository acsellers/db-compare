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

func DailyRevenue(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("debug") != "" {
		logQueries = true
		defer func() { logQueries = false }()
	}
	endDate := time.Now().Truncate(24 * time.Hour)
	if r.URL.Query().Get("end_date") != "" {
		endDate, _ = time.Parse("2006-01-02", r.URL.Query().Get("end_date"))
	}
	startDate := endDate.AddDate(0, 0, -7)
	if r.URL.Query().Get("start_date") != "" {
		startDate, _ = time.Parse("2006-01-02", r.URL.Query().Get("start_date"))
	}

	revenues, err := bob.All(
		r.Context(),
		db,
		models.Orders.Query(
			models.SelectWhere.Orders.OrderDate.GTE(startDate),
			models.SelectWhere.Orders.OrderDate.LTE(endDate),
			sm.Columns(
				models.Orders.Columns.OrderType,
				models.Orders.Columns.OrderDate,
				mysql.Raw("SUM(`orders`.`total`) as total_revenue"),
			),
			sm.GroupBy("orders.order_type, orders.order_date"),
		),
		scan.StructMapper[common.DailyRevenue](),
	)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid date", http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(revenues)
}
