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

func GeneralSales(w http.ResponseWriter, r *http.Request) {
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

	lines, err := bob.All(
		r.Context(),
		db,

		models.ItemSummaries.Query(
			sm.Columns(
				models.ReportingOrders.Columns.Title,
				models.ReportingOrders.Columns.ReportOrder,
				models.ItemSummaries.Columns.Name.As("item_name"),
				mysql.Raw("sum(item_summaries.order_count) as order_count"),
				mysql.Raw("sum(item_summaries.total_quantity) as quantity"),
				mysql.Raw("sum(item_summaries.total_sales) as total_sales"),
			),
			sm.InnerJoin("reporting_order").On(
				mysql.Raw("reporting_order.order_type = 'general'"),
				mysql.Raw("reporting_order.category = item_summaries.category"),
			),
			models.SelectWhere.ItemSummaries.OrderDate.GTE(startDate),
			models.SelectWhere.ItemSummaries.OrderDate.LTE(endDate),
			sm.GroupBy("reporting_order.title, reporting_order.report_order, item_summaries.name"),
			sm.OrderBy("reporting_order.report_order, item_summaries.name"),
		),
		scan.StructMapper[common.GeneralSalesReport](),
	)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid date", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(lines)
}
