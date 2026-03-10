CREATE TABLE products (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name varchar(255) NOT NULL,
    category varchar(32) NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_products_category ON products (category);

create table discounts (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name varchar(255) NOT NULL,
    category varchar(32) not null,
    discount_type varchar(12) NOT NULL,
    discount DECIMAL(10,4) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_discounts_category ON discounts (category);


CREATE TABLE locations (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name varchar(255) NOT NULL,
    address varchar(255) NOT NULL,
    city varchar(255) NOT NULL,
    state varchar(32) NOT NULL,
    zip varchar(5) NOT NULL,
    tax_rate DECIMAL(10,4) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
ALTER TABLE locations ADD UNIQUE (name);
CREATE INDEX idx_locations_state_city ON locations (state, city);


create table payment_names (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    payment_type varchar(12) not null,
    name varchar(255) not null,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_payment_names_name ON payment_names (name);

CREATE TABLE customers (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name varchar(255) NOT NULL,
    phone VARCHAR(24),
    email varchar(255),
    marketing_opt_in boolean DEFAULT false,
    external_id varchar(12),
    join_location_id BIGINT,
    last_location_id BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
ALTER TABLE customers ADD UNIQUE (external_id);
CREATE INDEX idx_customers_external_id ON customers (external_id);
ALTER TABLE customers ADD FOREIGN KEY (join_location_id) REFERENCES locations(id);
ALTER TABLE customers ADD FOREIGN KEY (last_location_id) REFERENCES locations(id);

CREATE TABLE orders (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    order_date DATETIME NOT NULL,
    customer_id BIGINT,
    discount_id BIGINT,
    order_type varchar(12) NOT NULL,
    subtotal DECIMAL(20,2) NOT NULL,
    discount_amount DECIMAL(20,2) NOT NULL,
    tax_amount DECIMAL(20,2) NOT NULL,
    total DECIMAL(20,2) NOT NULL,
    location_id BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_orders_order_date ON orders (order_date);
CREATE INDEX idx_orders_customer_id ON orders (customer_id);
CREATE INDEX idx_orders_discount_id ON orders (discount_id);
ALTER TABLE orders ADD FOREIGN KEY (customer_id) REFERENCES customers(id);
ALTER TABLE orders ADD FOREIGN KEY (discount_id) REFERENCES discounts(id);
ALTER TABLE orders ADD FOREIGN KEY (location_id) REFERENCES locations(id);

CREATE TABLE order_items (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    order_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    discount_id BIGINT,
    quantity INT NOT NULL,
    price DECIMAL(20,2) NOT NULL,
    discount_amount DECIMAL(20,2) NOT NULL,
    item_total DECIMAL(20,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
ALTER TABLE order_items ADD FOREIGN KEY (order_id) REFERENCES orders(id);
ALTER TABLE order_items ADD FOREIGN KEY (product_id) REFERENCES products(id);
ALTER TABLE order_items ADD FOREIGN KEY (discount_id) REFERENCES discounts(id);

CREATE TABLE order_payments (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    order_id BIGINT NOT NULL,
    payment_type varchar(1) NOT NULL,
    amount DECIMAL(20,2) NOT NULL,
    payment_info json COMMENT 'payments.Info',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
ALTER TABLE order_payments ADD FOREIGN KEY (order_id) REFERENCES orders(id);

create table reporting_order (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    category varchar(32) not null,
    order_type varchar(12) not null,
    report_order INT not null,
    title varchar(255) not null,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
INSERT INTO reporting_order (id, category, order_type, report_order, title)
VALUES 
(1, 'pc', 'members', 1, 'Member PC Items'),
(2, 'phone', 'members', 2, 'Member Phone & Wearable Items'),
(3, 'audio', 'members', 3, 'Member Audio Items'),
(4, 'taxes', 'members', 4, 'Member Tax Items'),
(5, 'discounts', 'members', 5, 'Member Discounts'),
(6, 'payments', 'members', 6, 'Member Payments'),
(7, 'pc', 'non-members', 7, 'Non-Member PC Items'),
(8, 'phone', 'non-members', 8, 'Non-Member Phone & Wearable Items'),
(9, 'audio', 'non-members', 9, 'Non-Member Audio Items'),
(10, 'taxes', 'non-members', 10, 'Non-Member Tax Items'),
(11, 'discounts', 'non-members', 11, 'Non-Member Discounts'),
(12, 'payments', 'non-members', 12, 'Non-Member Payments'),
(13, 'pc', 'general', 1, 'PC Items'),
(14, 'phone', 'general', 2, 'Phone & Wearable Items'),
(15, 'audio', 'general', 3, 'Audio Items'),
(16, 'taxes', 'general', 4, 'Tax Items'),
(17, 'discounts', 'general', 5, 'Discounts'),
(18, 'payments', 'general', 6, 'Payments');

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
CREATE OR REPLACE VIEW item_summaries as 

select p.id, p.name, p.category, o.order_type, date(o.order_date),
    sum(oi.quantity) as total_quantity,
    sum(oi.quantity * oi.price) as total_sales,
    count(distinct oi.id) as order_count,
    o.location_id
from products p
inner join order_items oi on p.id = oi.product_id
inner join orders o on oi.order_id = o.id
GROUP BY p.id, p.name, p.category, o.order_type, date(o.order_date), o.location_id
UNION 
select d.id, d.name, 'discounts', o.order_type, date(o.order_date),
	count(o.id) as total_quantity,
	-sum(o.discount_amount) as total_sales,
	count(o.id) as order_count,
    o.location_id
from discounts d
inner join orders o on o.discount_id = d.id
group by d.id, d.name, o.order_type, date(o.order_date), o.location_id 
UNION
select d.id, d.name, 'discounts', o.order_type, date(o.order_date),
  count(oi.id) as total_quantity,
  -sum(oi.discount_amount) as total_sales,
  count(distinct o.id) as order_count,
  o.location_id
from discounts d
inner join order_items oi on oi.discount_id = d.id
inner join orders o on o.id = oi.order_id 
group by d.id, d.name, o.order_type, date(o.order_date), o.location_id
UNION
select NULL, 'Sales Tax', 'taxes', o.order_type, date(o.order_date),
  count(o.id) as total_quantity,
  sum(o.tax_amount) as total_sales,
  count(o.id) as order_count,
  o.location_id
from orders o
group by o.order_type, date(o.order_date), o.location_id
UNION
select null, coalesce(pn.name, op.payment_type), 'payments', o.order_type, date(o.order_date),
  count(op.id) as total_quantity,
  -sum(op.amount) as total_sales,
  count(distinct o.id) as order_count,
  o.location_id
from order_payments op
inner join orders o on op.order_id  = o.id 
left join payment_names pn on pn.payment_type = op.payment_type 
group by coalesce(pn.name, op.payment_type), o.order_type, date(o.order_date), o.location_id;
