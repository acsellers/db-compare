package main

import (
	"encoding/json"
	"net/http"

	"github.com/acsellers/golang-db-compare/store/common"
	"github.com/acsellers/golang-db-compare/store/mysql/bob/models"
	"github.com/stephenafamo/bob/dialect/mysql"
	"github.com/stephenafamo/bob/dialect/mysql/sm"
)

func JSONQuery(w http.ResponseWriter, r *http.Request) {
	cardType := r.URL.Query().Get("card_type")

	op, err := models.OrderPayments.Query(
		sm.Columns(
			models.OrderPayments.Columns.OrderID,
			models.OrderPayments.Columns.Amount,
		),
		sm.Where(mysql.Raw("payment_info->>'$.card_type' = ?", cardType)),
	).All(r.Context(), db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	ret := common.SalesByCardType{}
	for _, v := range op {
		ret.OrderID = append(ret.OrderID, v.OrderID)
		ret.TotalSales += v.Amount.InexactFloat64()
		ret.OrderCount++
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ret)
}
