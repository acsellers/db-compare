package main

import "net/http"

func SaleSearch(w http.ResponseWriter, r *http.Request) {
	/*
		select orders.*, customers.name
		from orders
		left join customers on orders.customer_id = customers.id
		where (customers.name like ?)
		and (orders.order_date >= ?)
		and (orders.order_date <= ?)
		and (orders.order_type = ?)
		and (orders.total >= ?)
		and (orders.total <= ?)
		order by orders.order_date desc;
	*/
	// TODO: implement sale search
}
