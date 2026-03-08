package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/scan"
)

var (
	db   bob.Executor
	ldb  bob.Executor
	root bob.Executor
	conn *sql.DB
)

func init() {
	var err error
	os.Setenv("MYSQL_DSN", "service_store:gopher@tcp(localhost:3306)/store?parseTime=true")
	conn, err = sql.Open("mysql", os.Getenv("MYSQL_DSN"))
	if err != nil {
		log.Fatal(err)
	}
	db = bob.NewDB(conn)
	root = db
	ldb = QueryLogger(bob.NewDB(conn))
}

func main() {
	http.HandleFunc("GET /01/sales/{id}", GetSale)
	http.HandleFunc("GET /01/sales/", GetSales)

	http.HandleFunc("POST /02/sales/$", CreateSale)

	http.HandleFunc("POST /03/sales/search", SaleSearch)

	http.HandleFunc("POST /04/customers/$", BulkLoadCustomers)

	http.HandleFunc("POST /05/customers/update", CustomerUpdate)
	http.HandleFunc("POST /05/customers/update2", CustomerUpdate2)

	http.HandleFunc("GET /06/customers/payment_cards/", JSONQuery)
	http.HandleFunc("GET /06/webhook/update_payment", JSONUpdate)
	http.HandleFunc("GET /06/locations/payments", JSONReport)

	http.HandleFunc("GET /07/generic_report", WithQuery)

	http.HandleFunc("GET /08/reports/daily-revenue/", DailyRevenue)
	http.HandleFunc("GET /08/reports/customer-sales/", CustomerSales)
	http.HandleFunc("GET /08/reports/customer-sales/alt", CustomerSales2)
	http.HandleFunc("GET /08/reports/daily-sold-items/", DailySoldItems)
	http.HandleFunc("GET /08/reports/daily-sold-items/alt", DailySoldItems2)

	http.HandleFunc("GET /09/reports/general-sales/", GeneralSales)
	http.HandleFunc("GET /09/reports/general-sales/alt", GeneralSales2)
	http.HandleFunc("GET /09/reports/weekly-sales/", TypedSales)
	http.HandleFunc("GET /09reports/weekly-sales/alt", TypedSales2)

	port := "8080"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	fmt.Println("Listening on :" + port)
	if os.Getenv("BENCHMARK") == "true" {
		log.Fatal(http.ListenAndServe(":"+port, nil))
	} else {
		log.Fatal(http.ListenAndServe(":"+port, &QueryLoggingHandler{Next: http.DefaultServeMux}))
	}
}

var logQueries = false

type QueryLoggingHandler struct {
	Next http.Handler
}

func (qlh *QueryLoggingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("debug") != "" {
		logQueries = true
		db = ldb
	}
	qlh.Next.ServeHTTP(w, r)
	if logQueries {
		logQueries = false
		db = root
	}
}

func QueryLogger(db bob.DB) bob.Executor {
	return &QueryLoggerExecutor{db: db}
}

type QueryLoggerExecutor struct {
	db bob.DB
}

func (qle *QueryLoggerExecutor) QueryContext(ctx context.Context, query string, args ...any) (scan.Rows, error) {
	if logQueries {
		log.Println(query)
	}
	return qle.db.QueryContext(ctx, query, args...)
}
func (qle *QueryLoggerExecutor) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	if logQueries {
		log.Println(query)
	}
	return qle.db.ExecContext(ctx, query, args...)
}
