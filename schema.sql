CREATE TABLE products (
  id INT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  price DECIMAL(10, 2) NOT NULL,
  description TEXT,
  image VARCHAR(255)
);

CREATE TABLE customers (
   id INT AUTO_INCREMENT PRIMARY KEY,
   name VARCHAR(255) NOT NULL,
   email VARCHAR(255) NOT NULL,
   password VARCHAR(255) NOT NULL,
   token VARCHAR(255),
   is_admin BOOLEAN DEFAULT false
);

CREATE TABLE orders (
    id INT AUTO_INCREMENT PRIMARY KEY,
    customer_id INT NOT NULL,
    order_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(50),
    FOREIGN KEY (customer_id) REFERENCES customers (customer_id)
);

CREATE TABLE order_items (
     id INT AUTO_INCREMENT PRIMARY KEY,
     order_id INT NOT NULL,
     product_id INT NOT NULL,
     quantity INT NOT NULL,
     FOREIGN KEY (order_id) REFERENCES orders (order_id),
     FOREIGN KEY (product_id) REFERENCES products (product_id)
);