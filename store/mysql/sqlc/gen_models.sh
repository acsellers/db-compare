export MYSQL_DSN="service_store:gopher@tcp(localhost:3306)/store"

go run github.com/stephenafamo/bob/gen/bobgen-mysql@latest
