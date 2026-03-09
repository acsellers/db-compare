package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/acsellers/golang-db-compare/store/common"
	"github.com/acsellers/golang-db-compare/store/mysql/bob/models"
	"github.com/stephenafamo/bob/dialect/mysql"
)

func GetSales(w http.ResponseWriter, r *http.Request) {
	locationId, err := strconv.Atoi(r.URL.Query().Get("location"))
	if err != nil {
		http.Error(w, "Invalid location ID", http.StatusBadRequest)
		return
	}
	startDate, err := time.Parse("2006-01-02", r.URL.Query().Get("start_date"))
	if err != nil {
		http.Error(w, "Invalid Start Date", http.StatusBadRequest)
		return
	}

	// end date is inclusive, so move forward by 1 day and use less than in the query
	endDate, err := time.Parse("2006-01-02", r.URL.Query().Get("end_date"))
	if err != nil {
		http.Error(w, "Invalid End Date", http.StatusBadRequest)
		return
	}
	endDate.AddDate(0, 0, 1)

	orders, err := models.Orders.Query(
		models.SelectWhere.Orders.LocationID.EQ(int64(locationId)),
		models.SelectWhere.Orders.OrderDate.GTE(startDate),
		models.SelectWhere.Orders.OrderDate.LT(endDate),
		models.Preload.Order.Location(
			mysql.PreloadOnly("name"),
		),
		models.Preload.Order.Customer(
			mysql.PreloadOnly("name"),
		),
		models.SelectThenLoad.Order.OrderItems(
			models.Preload.OrderItem.Product(
				mysql.PreloadOnly("name", "category"),
			),
		),
		models.SelectThenLoad.Order.OrderPayments(),
	).All(r.Context(), db)
	if err != nil {
		http.Error(w, "Invalid sale ID", http.StatusBadRequest)
		return
	}

	sales := make([]common.Sale, 0, len(orders))
	for _, order := range orders {
		sales = append(sales, toSale(order))
	}

	json.NewEncoder(w).Encode(sales)
}
