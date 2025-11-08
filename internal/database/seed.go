package database

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	_ "modernc.org/sqlite"
)

func (db *DB) CheckEmpty() error {
	tables, err := db.GetTables()
	if err != nil {
		return err
	}

	if len(tables) > 0 {
		prompt, err := promptOverwrite()
		if err != nil {
			return err
		}
		if !prompt {
			return nil
		}
	}

	fmt.Println("Seeding database with test data...")

	err = db.execSeed()
	if err != nil {
		return err
	}

	fmt.Println("Database successfully seeded.")

	return nil
}

func promptOverwrite() (bool, error) {
	warning := `
Warning: This database already contains tables.
Seeding will overwrite certain tables and data.
Do you want to continue? (y/N): `

	fmt.Print(warning)

	reader := bufio.NewReader(os.Stdin)
	res, err := reader.ReadString('\n')
	if err != nil {
		return false, fmt.Errorf("Failed to read input: %w", err)
	}

	res = strings.TrimSpace(strings.ToLower(res))
	if res != "yes" && res != "y" {
		fmt.Println("Seeding cancelled.")
		return false, nil
	}

	return true, nil
}

func (db *DB) execSeed() error {
	query := `PRAGMA foreign_keys = OFF;

DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS reviews;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS users;

PRAGMA foreign_keys = ON;

CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    full_name TEXT NOT NULL,
    age INTEGER,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN DEFAULT 1
);

INSERT INTO users (username, email, full_name, age, is_active) VALUES
    ('jdoe', 'john.doe@example.com', 'John Doe', 28, 1),
    ('asmith', 'alice.smith@example.com', 'Alice Smith', 34, 1),
    ('bwilliams', 'bob.williams@example.com', 'Bob Williams', 45, 1),
    ('cjohnson', 'carol.johnson@example.com', 'Carol Johnson', 29, 0),
    ('dmiller', 'david.miller@example.com', 'David Miller', 52, 1),
    ('ebrown', 'emma.brown@example.com', 'Emma Brown', 23, 1),
    ('fdavis', 'frank.davis@example.com', 'Frank Davis', 41, 1),
    ('gwilson', 'grace.wilson@example.com', 'Grace Wilson', 36, 0),
    ('hmoore', 'henry.moore@example.com', 'Henry Moore', 31, 1),
    ('itaylor', 'iris.taylor@example.com', 'Iris Taylor', 27, 1);

CREATE TABLE products (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    category TEXT NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    stock_quantity INTEGER NOT NULL,
    description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO products (name, category, price, stock_quantity, description) VALUES
    ('Laptop Pro 15', 'Electronics', 1299.99, 45, 'High-performance laptop with 16GB RAM'),
    ('Wireless Mouse', 'Electronics', 29.99, 230, 'Ergonomic wireless mouse with USB receiver'),
    ('Office Chair', 'Furniture', 249.50, 67, 'Comfortable ergonomic office chair'),
    ('Standing Desk', 'Furniture', 499.99, 23, 'Adjustable height standing desk'),
    ('USB-C Cable', 'Electronics', 12.99, 450, 'Durable USB-C charging cable'),
    ('Desk Lamp', 'Furniture', 39.99, 120, 'LED desk lamp with adjustable brightness'),
    ('Mechanical Keyboard', 'Electronics', 89.99, 88, 'RGB mechanical keyboard with blue switches'),
    ('Monitor 27"', 'Electronics', 349.99, 34, '4K UHD monitor with HDR support'),
    ('Bookshelf', 'Furniture', 129.99, 41, 'Five-tier wooden bookshelf'),
    ('Webcam HD', 'Electronics', 79.99, 95, '1080p webcam with built-in microphone'),
    ('Coffee Maker', 'Appliances', 69.99, 52, 'Programmable coffee maker with thermal carafe'),
    ('Blender', 'Appliances', 54.99, 38, 'High-speed blender for smoothies'),
    ('Microwave', 'Appliances', 129.99, 28, 'Compact microwave with 900W power'),
    ('Headphones', 'Electronics', 159.99, 102, 'Noise-cancelling wireless headphones'),
    ('Backpack', 'Accessories', 49.99, 145, 'Water-resistant laptop backpack');

CREATE TABLE orders (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    order_date DATETIME DEFAULT CURRENT_TIMESTAMP,
    total_amount DECIMAL(10, 2) NOT NULL,
    status TEXT NOT NULL CHECK(status IN ('pending', 'processing', 'shipped', 'delivered', 'cancelled')),
    shipping_address TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

INSERT INTO orders (user_id, order_date, total_amount, status, shipping_address) VALUES
    (1, '2024-10-15 10:23:45', 1329.98, 'delivered', '123 Main St, Springfield, IL 62701'),
    (2, '2024-10-18 14:32:11', 249.50, 'delivered', '456 Oak Ave, Portland, OR 97201'),
    (3, '2024-10-20 09:15:33', 89.99, 'shipped', '789 Pine Rd, Austin, TX 78701'),
    (1, '2024-10-22 16:44:22', 79.99, 'processing', '123 Main St, Springfield, IL 62701'),
    (4, '2024-10-25 11:05:44', 499.99, 'pending', '321 Elm St, Seattle, WA 98101'),
    (5, '2024-10-27 13:28:56', 159.99, 'delivered', '654 Maple Dr, Boston, MA 02101'),
    (6, '2024-10-29 08:42:18', 549.98, 'shipped', '987 Cedar Ln, Denver, CO 80201'),
    (2, '2024-11-01 10:11:33', 349.99, 'processing', '456 Oak Ave, Portland, OR 97201'),
    (7, '2024-11-02 15:55:27', 129.99, 'pending', '147 Birch St, Miami, FL 33101'),
    (9, '2024-11-04 12:33:44', 219.97, 'delivered', '258 Willow Way, Phoenix, AZ 85001');

CREATE TABLE order_items (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    order_id INTEGER NOT NULL,
    product_id INTEGER NOT NULL,
    quantity INTEGER NOT NULL,
    price_at_purchase DECIMAL(10, 2) NOT NULL,
    FOREIGN KEY (order_id) REFERENCES orders(id),
    FOREIGN KEY (product_id) REFERENCES products(id)
);

INSERT INTO order_items (order_id, product_id, quantity, price_at_purchase) VALUES
    (1, 1, 1, 1299.99),
    (1, 2, 1, 29.99),
    (2, 3, 1, 249.50),
    (3, 7, 1, 89.99),
    (4, 10, 1, 79.99),
    (5, 4, 1, 499.99),
    (6, 14, 1, 159.99),
    (7, 8, 1, 349.99),
    (7, 5, 1, 12.99),
    (7, 6, 5, 39.99),
    (8, 8, 1, 349.99),
    (9, 13, 1, 129.99),
    (10, 15, 2, 49.99),
    (10, 11, 2, 69.99);

CREATE TABLE reviews (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    product_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    rating INTEGER NOT NULL CHECK(rating >= 1 AND rating <= 5),
    comment TEXT,
    review_date DATETIME DEFAULT CURRENT_TIMESTAMP,
    helpful_count INTEGER DEFAULT 0,
    FOREIGN KEY (product_id) REFERENCES products(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

INSERT INTO reviews (product_id, user_id, rating, comment, review_date, helpful_count) VALUES
    (1, 1, 5, 'Excellent laptop! Fast and reliable.', '2024-10-20 14:22:33', 12),
    (1, 5, 4, 'Great performance but a bit pricey.', '2024-10-28 09:45:11', 8),
    (2, 2, 5, 'Perfect mouse, very comfortable.', '2024-10-19 16:33:22', 5),
    (3, 2, 4, 'Comfortable chair, good back support.', '2024-10-23 11:18:44', 15),
    (7, 3, 5, 'Love this keyboard! Keys feel amazing.', '2024-10-25 13:42:55', 22),
    (8, 2, 5, 'Crystal clear display, worth every penny.', '2024-11-03 10:25:18', 9),
    (10, 1, 3, 'Decent webcam but average mic quality.', '2024-10-26 15:11:33', 4),
    (14, 5, 5, 'Best headphones I have ever owned!', '2024-11-01 12:44:22', 18),
    (15, 9, 4, 'Good backpack with lots of compartments.', '2024-11-05 09:33:11', 6);

CREATE TABLE categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    parent_category_id INTEGER,
    FOREIGN KEY (parent_category_id) REFERENCES categories(id)
);

INSERT INTO categories (name, description, parent_category_id) VALUES
    ('Electronics', 'Electronic devices and accessories', NULL),
    ('Furniture', 'Home and office furniture', NULL),
    ('Appliances', 'Home appliances', NULL),
    ('Accessories', 'Various accessories', NULL),
    ('Computers', 'Computer hardware', 1),
    ('Audio', 'Audio equipment', 1),
    ('Office', 'Office furniture', 2);
`

	if _, err := db.Exec(query); err != nil {
		return err
	}

	return nil
}
