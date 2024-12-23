CREATE TABLE IF NOT EXIST orders(
    id CHAR(27) PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    account_id TIMESTAMP WITH TIME ZONE NOT NULL,
    total_price MONEY NOT NULL
);

CREATE TABLE IF NOT EXIST order_products(
    order_id CHAR(27) REFERENCES order (id) ON DELETE CASCADE,
    product_id CHAR(27),
    quantity INT NOT NULL,
    PRYMARY KEY (product_id, order_id)
);