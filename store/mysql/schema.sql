create user 'service_store'@'%' identified by 'gopher';
create database store;
grant all on store.* to 'service_store'@'%';
flush privileges;

CREATE TABLE products (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name varchar(255) NOT NULL,
    category varchar(32) NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
create table discounts (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name varchar(255) NOT NULL,
    category varchar(32) not null,
    discount_type varchar(12) NOT NULL,
    discount DECIMAL(10,4) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE customers (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name varchar(255) NOT NULL,
    phone VARCHAR(24),
    email varchar(255),
    marketing_opt_in boolean DEFAULT false,
    external_id varchar(12),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
ALTER TABLE customers ADD UNIQUE (external_id);

CREATE TABLE orders (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    order_date DATE NOT NULL,
    customer_id BIGINT,
    discount_id BIGINT,
    order_type varchar(12) NOT NULL,
    subtotal DECIMAL(10,2) NOT NULL,
    discount_amount DECIMAL(10,2) NOT NULL,
    tax_amount DECIMAL(10,2) NOT NULL,
    total DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE order_items (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    order_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    discount_id BIGINT,
    quantity INT NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    discount_amount DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE order_payments (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    order_id BIGINT NOT NULL,
    payment_type varchar(12) NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    payment_info json COMMENT 'payments.Info',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

create table reporting_order (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    category varchar(32) not null,
    order_type varchar(12) not null,
    report_order INT not null,
    title varchar(255) not null,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

create table payment_names (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    payment_type varchar(12) not null,
    name varchar(255) not null,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

create table dim_date (
    date DATE PRIMARY KEY,
    month int NOT NULL,
    year int NOT NULL,
    quarter int NOT NULL,
    day_of_week int NOT NULL,
    day_of_month int NOT NULL,
    day_of_year int NOT NULL,
    week_of_year int NOT NULL,
    week_of_month int NOT NULL
);

CREATE VIEW item_summaries as 
select p.id, p.name, p.category, o.order_type, o.order_date,
    sum(oi.quantity) as total_quantity,
    sum(oi.quantity * oi.price) as total_sales,
    count(distinct oi.id) as order_count
from products p
inner join order_items oi on p.id = oi.product_id
inner join orders o on oi.order_id = o.id
GROUP BY p.id, p.name, p.category, o.order_type, o.order_date
UNION 
select d.id, d.name, 'discounts', o.order_type, o.order_date,
	count(o.id) as total_quantity,
	-sum(o.discount_amount) as total_sales,
	count(o.id) as order_count
from discounts d
inner join orders o on o.discount_id = d.id
group by d.id, d.name, o.order_type, o.order_date 
UNION
select d.id, d.name, 'discounts', o.order_type, o.order_date,
  count(oi.id) as total_quantity,
  -sum(oi.discount_amount) as total_sales,
  count(distinct o.id) as order_count
from discounts d
inner join order_items oi on oi.discount_id = d.id
inner join orders o on o.id = oi.order_id 
group by d.id, d.name, o.order_type, o.order_date
UNION
select NULL, 'Sales Tax', 'taxes', o.order_type, o.order_date,
  count(o.id) as total_quantity,
  sum(o.tax_amount) as total_sales,
  count(o.id) as order_count
from orders o
group by o.order_type, o.order_date
UNION
select null, coalesce(pn.name, op.payment_type), 'payments', o.order_type, o.order_date,
  count(op.id) as total_quantity,
  -sum(op.amount) as total_sales,
  count(distinct o.id) as order_count
from order_payments op
inner join orders o on op.order_id  = o.id 
left join payment_names pn on pn.payment_type = op.payment_type 
group by coalesce(pn.name, op.payment_type), o.order_type, o.order_date;

ALTER TABLE orders ADD FOREIGN KEY (customer_id) REFERENCES customers(id);
ALTER TABLE orders ADD FOREIGN KEY (discount_id) REFERENCES discounts(id);
ALTER TABLE order_items ADD FOREIGN KEY (order_id) REFERENCES orders(id);
ALTER TABLE order_items ADD FOREIGN KEY (product_id) REFERENCES products(id);
ALTER TABLE order_items ADD FOREIGN KEY (discount_id) REFERENCES discounts(id);
ALTER TABLE order_payments ADD FOREIGN KEY (order_id) REFERENCES orders(id);