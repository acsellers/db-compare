package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/acsellers/golang-db-compare/store/common"
	"github.com/acsellers/golang-db-compare/store/mysql/bob/models"
)

func JSONReport(w http.ResponseWriter, r *http.Request) {
	locationID := r.URL.Query().Get("location")
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")
	locationIDInt, _ := strconv.Atoi(locationID)
	startDateObj, _ := time.Parse("2006-01-02", startDate)
	endDateObj, _ := time.Parse("2006-01-02", endDate)
	if locationIDInt == 0 || startDateObj.IsZero() || endDateObj.IsZero() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	orders, err := models.Orders.Query(
		models.SelectThenLoad.Order.OrderPayments(),
		models.SelectWhere.Orders.OrderDate.GT(startDateObj),
		models.SelectWhere.Orders.OrderDate.LT(endDateObj),
		models.SelectWhere.Orders.LocationID.EQ(int64(locationIDInt)),
	).All(r.Context(), db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	ret := make([]common.LocationPaymentInfo, 0, len(orders))

}
