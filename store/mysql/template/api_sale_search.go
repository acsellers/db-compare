package main

import "net/http"

func SaleSearch(w http.ResponseWriter, r *http.Request) {
	/*
		select orders.*, customers.name
		from orders
		left join customers on orders.customer_id = customers.id
		where (customers.name like @name or @name = '')
		and (orders.order_date >= @order_start or @order_start = '')
		and (orders.order_date <= @order_end or @order_end = '')
		and (orders.order_type = @order_type or @order_type = '')
		and (orders.total >= @total_start or @total_start = 0)
		and (orders.total <= @total_end or @total_end = 0)
		order by orders.order_date desc;
	*/
	// TODO: implement sale search
}
