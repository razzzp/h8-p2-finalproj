-- created at, updated at an deleted at columns will be added
--  on gorm auto migrate
CREATE TABLE cars (
    id SERIAL PRIMARY KEY,
    wheel_drive INT NOT NULL,
    type VARCHAR(50) NOT NULL,
    seats INT NOT NULL,
    transmission VARCHAR(20) NOT NULL,
    manufacturer VARCHAR(255) NOT NULL,
    car_model VARCHAR(255) NOT NULL,
    year INT NOT NULL,
    stock INT NOT NULL,
    rate_per_day DECIMAL NOT NULL,
    UNIQUE(manufacturer, car_model)
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'user',
    password VARCHAR(255) NOT NULL,
    deposit DECIMAL NOT NULL DEFAULT 0
);

CREATE TABLE rentals (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) NOT NULL,
    car_id INT REFERENCES cars(id) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    total_price DECIMAL NOT NULL
);

CREATE TABLE top_ups (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) NOT NULL,
    amount DECIMAL NOT NULL
);

CREATE TABLE payments (
    id SERIAL PRIMARY KEY,
    -- polymorphic FK to either rentals or top_ups
    purchase_id INT NOT NULL,
    purchase_type TEXT NOT NULL,
    payment_url VARCHAR(255) NOT NULL,
    status VARCHAR(100) NOT NULL,
    payment_method VARCHAR(100),
    total_payment DECIMAL
);

-- Insert dummy data into cars table
INSERT INTO cars (wheel_drive, type, seats, transmission, manufacturer, car_model, year, stock, rate_per_day) VALUES
(4, 'SUV', 5, 'Automatic', 'Toyota', 'RAV4', 2022, 10, 500000),
(2, 'Sedan', 5, 'Manual', 'Honda', 'Civic', 2021, 8, 400000),
(4, 'Truck', 2, 'Automatic', 'Ford', 'F-150', 2023, 5, 700000),
(4, 'SUV', 7, 'Automatic', 'Chevrolet', 'Tahoe', 2022, 6, 600000),
(2, 'Coupe', 4, 'Manual', 'BMW', 'M4', 2021, 3, 80.00);

-- Insert dummy data into users table
INSERT INTO users (name, email, password, deposit) VALUES
('John Doe', 'john@example.com', 'password123', 200000),
('Jane Smith', 'jane@example.com', 'securepass', 150000),
('Mike Johnson', 'mike@example.com', 'mikepass', 300000),
('Emily Davis', 'emily@example.com', 'emilysecure', 100000),
('David Wilson', 'david@example.com', 'david123', 250000);

-- Insert dummy data into rentals table
INSERT INTO rentals (user_id, car_id, start_date, end_date, total_price) VALUES
(1, 1, '2024-09-01', '2024-09-05', 2000000),
(2, 2, '2024-09-10', '2024-09-12', 800000),
(3, 3, '2024-09-15', '2024-09-20', 3500000),
(4, 4, '2024-09-05', '2024-09-08', 1800000),
(5, 5, '2024-09-20', '2024-09-25', 4000000);

-- Insert dummy data into top_ups table
INSERT INTO top_ups (user_id, amount) VALUES
(1, 50000),
(2, 30000),
(3, 70000),
(4, 100000),
(5, 25000);

-- Insert dummy data into payments table
INSERT INTO payments (purchase_id, purchase_type, payment_url, status, payment_method, total_payment) VALUES
(1, 'rentals','http://payment.com/1', 'Completed', 'Credit Card', 200000),
(2, 'rentals', 'http://payment.com/2', 'Completed', 'PayPal', 80000),
(3, 'rentals', 'http://payment.com/3', 'Pending', 'Credit Card', 350000),
(4, 'rentals', 'http://payment.com/4', 'Completed', 'Debit Card', 180000),
(5, 'rentals', 'http://payment.com/5', 'Failed', 'Credit Card', 400000);