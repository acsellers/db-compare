package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/acsellers/golang-db-compare/store/common"
	"github.com/acsellers/golang-db-compare/store/mysql/bob/models"
	"github.com/shopspring/decimal"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/mysql/dialect"
)

func SaleSearch(w http.ResponseWriter, r *http.Request) {

	args := common.SaleSearch{}
	json.NewDecoder(r.Body).Decode(&args)

	filters := []bob.Mod[*dialect.SelectQuery]{}
	if args.StartDate != "" {
		sd, err := time.Parse("2006-01-02", args.StartDate)
		if err == nil {
			filters = append(filters, models.SelectWhere.Orders.OrderDate.GTE(sd))
		}
	}
	if args.EndDate != "" {
		ed, err := time.Parse("2006-01-02", args.EndDate)
		if err == nil {
			filters = append(filters, models.SelectWhere.Orders.OrderDate.LTE(ed))
		}
	}
	if args.OrderType != "" {
		filters = append(filters, models.SelectWhere.Orders.OrderType.EQ(args.OrderType))
	}
	if args.MinTotal != 0 {
		filters = append(filters, models.SelectWhere.Orders.Total.GTE(decimal.NewFromFloat(args.MinTotal)))
	}
	if args.MaxTotal != 0 {
		filters = append(filters, models.SelectWhere.Orders.Total.LTE(decimal.NewFromFloat(args.MaxTotal)))
	}

	if args.LocationID != 0 {
		filters = append(filters, models.SelectWhere.Orders.LocationID.EQ(args.LocationID))
	}

	if args.CustomerName != "" {
		filters = append(filters, models.SelectJoins.Orders.InnerJoin.Customer)
		filters = append(filters, models.SelectWhere.Customers.Name.Like(args.CustomerName))
	}

	sales, err := models.Orders.Query(
		filters...,
	).All(r.Context(), db)
	if err != nil {
		http.Error(w, "Invalid sale search", http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(sales)
}
