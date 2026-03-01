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
	http.HandleFunc("GET /sales/{id}", GetSale)
	http.HandleFunc("POST /sales/$", CreateSale)
	http.HandleFunc("POST /customers/$", BulkLoadCustomers)
	http.HandleFunc("POST /customers/update", CustomerUpdate)
	http.HandleFunc("POST /customers/update2", CustomerUpdate2)
	http.HandleFunc("POST /sales/search", SaleSearch)
	http.HandleFunc("GET /reports/daily-sold-items/", DailySoldItems)
	http.HandleFunc("GET /reports/daily-revenue/", DailyRevenue)
	http.HandleFunc("GET /reports/customer-sales/", CustomerSales)
	http.HandleFunc("GET /reports/general-sales/", GeneralSales)
	http.HandleFunc("GET /reports/typed-sales/", TypedSales)

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
