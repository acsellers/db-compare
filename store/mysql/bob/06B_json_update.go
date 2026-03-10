package main

import (
	"encoding/json"
	"net/http"

	"github.com/aarondl/opt/omitnull"
	"github.com/acsellers/golang-db-compare/store/common"
	"github.com/acsellers/golang-db-compare/store/mysql/bob/models"
	"github.com/stephenafamo/bob/types"
)

type paymentUpdateRequest struct {
	OrderID           int64  `json:"order_id"`
	TransactionID     string `json:"transaction_id"`
	TransactionResult string `json:"transaction_result"`
}

func JSONUpdate(w http.ResponseWriter, r *http.Request) {
	var req paymentUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	op, err := models.FindOrderPayment(r.Context(), db, req.OrderID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	pi := common.OrderPaymentInfo{}
	if op.PaymentInfo.IsValue() {
		if err := json.Unmarshal(op.PaymentInfo.MustGet().Val, &pi); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	pi.TransactionID = req.TransactionID
	pi.TransactionResult = req.TransactionResult
	b, _ := json.Marshal(pi)
	ops := &models.OrderPaymentSetter{
		PaymentInfo: omitnull.From(types.JSON[json.RawMessage]{Val: b}),
	}
	if err := op.Update(r.Context(), db, ops); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
