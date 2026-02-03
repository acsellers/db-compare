package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/acsellers/golang-db-compare/store/mysql/sqlc/models"
	_ "github.com/go-sql-driver/mysql"
)

var (
	db   *models.Queries
	conn *sql.DB
)

func init() {
	var err error
	os.Setenv("MYSQL_DSN", "service_store:gopher@tcp(localhost:3306)/store?parseTime=true")
	conn, err = sql.Open("mysql", os.Getenv("MYSQL_DSN"))
	if err != nil {
		log.Fatal(err)
	}
	db = models.New(QueryLogger(conn))
}

func main() {
	http.HandleFunc("GET /sales/{id}", GetSale)
	http.HandleFunc("POST /sales/$", CreateSale)
	http.HandleFunc("POST /customers/$", BulkLoadCustomers)
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
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

var logQueries = false

func QueryLogger(db *sql.DB) models.DBTX {
	return &QueryLoggerExecutor{db: db}
}

type QueryLoggerExecutor struct {
	db *sql.DB
}

func (qle *QueryLoggerExecutor) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	if logQueries {
		log.Println(query)
	}
	return qle.db.QueryContext(ctx, query, args...)
}
func (qle *QueryLoggerExecutor) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	if logQueries {
		log.Println(query)
	}
	return qle.db.QueryRowContext(ctx, query, args...)
}
func (qle *QueryLoggerExecutor) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	if logQueries {
		log.Println(query)
	}
	return qle.db.ExecContext(ctx, query, args...)
}
func (qle *QueryLoggerExecutor) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	if logQueries {
		log.Println(query)
	}
	return qle.db.PrepareContext(ctx, query)
}
