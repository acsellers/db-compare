-- name: GetDiscount :many
SELECT * FROM discounts WHERE id IN (sqlc.slice('ids'));

-- name: GetProducts :many
SELECT * FROM products WHERE id IN (sqlc.slice('ids'));

-- name: GetSale :one
select orders.*, customers.name
from orders 
left join customers on orders.customer_id = customers.id
where orders.id = ?;

-- name: GetSaleItems :many
select order_items.*, 
  products.name, products.category 
from order_items 
left join products on order_items.product_id = products.id
where order_items.order_id = ?;

-- name: GetSalePayments :many
select * from order_payments where order_id = ?;

-- name: CreateSale :exec
insert into orders (order_date, customer_id, discount_id, order_type, subtotal, discount_amount, tax_amount, total)
values (?, ?, ?, ?, ?, ?, ?, ?);

-- name: CreateSaleItems :exec
insert into order_items (order_id, product_id, discount_id, quantity, price, discount_amount)
values (?, ?, ?, ?, ?, ?);

-- name: CreateSalePayments :exec
insert into order_payments (order_id, payment_type, amount, payment_info)
values (?, ?, ?, ?);

-- name: InsertCustomers :copyfrom
INSERT INTO customers (name, email, phone) VALUES (?, ?, ?);

-- name: DailySoldItems :many
select id, name, category, 
cast(sum(total_quantity) as signed) as total_quantity, 
cast(sum(total_sales) as double) as total_sales
from item_summaries
where order_date = ?
group by id, name, category
order by category, name;

-- name: DailyRevenue :many
select order_type,sum(total) as total_revenue
from orders
where order_date >= sqlc.arg(start_date) and order_date <= sqlc.arg(end_date)
group by order_type;

-- name: CustomerSales :many
select customers.id, customers.name, 
  sum(orders.total) as total_sales, count(*) as total_orders
from customers
join orders on customers.id = orders.customer_id
where orders.order_date >= sqlc.arg(start_date) and orders.order_date <= sqlc.arg(end_date)
group by customers.id, customers.name
order by total_sales desc;

-- name: GeneralSales :many
select ro.title, ro.report_order, t.name, 
sum(t.order_count) as order_count,
sum(t.total_quantity) as quantity,
sum(t.total_sales) as total_sales
from item_summaries t 
inner join reporting_order ro on ro.order_type = 'general' and ro.category = t.category
where t.order_date >= sqlc.arg(start_date) and t.order_date <= sqlc.arg(end_date)
group by ro.title, ro.report_order, t.name 
order by ro.report_order, t.name;

-- name: WeeklyTypedSales :many
select dim_date.year, dim_date.WEEK_OF_YEAR, ro.title, ro.report_order, t.name, 
sum(t.order_count) as order_count,
sum(t.total_quantity) as quantity,
sum(t.total_sales) as total_sales
from item_summaries t 
inner join reporting_order ro on ro.order_type = t.order_type and ro.category = t.category
inner join dim_date on dim_date.date = t.order_date
where t.order_date >= sqlc.arg(start_date) and t.order_date <= sqlc.arg(end_date)
group by dim_date.year, dim_date.WEEK_OF_YEAR,ro.title, ro.report_order, t.name 
order by dim_date.year, dim_date.WEEK_OF_YEAR, ro.report_order, t.name;

-- name: SearchSales :many
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