CREATE TABLE IF NOT EXISTS users (
    Email VARCHAR(50) PRIMARY KEY,
    Password VARCHAR(60) NOT NULL,
    Name VARCHAR(50),
    PhoneNumber VARCHAR(20)
);

CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    description TEXT,
    price INTEGER,
    image_url VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS cart (
    id SERIAL PRIMARY KEY,
    user_email VARCHAR(50) REFERENCES users(Email),
    product_id INTEGER REFERENCES products(id),
    quantity INTEGER,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS favorites (
    id SERIAL PRIMARY KEY,
    user_email VARCHAR(50) REFERENCES users(Email),
    product_id INTEGER REFERENCES products(id)
);

CREATE TABLE IF NOT EXISTS usersessions (
    token VARCHAR(64) PRIMARY KEY,
    user_email VARCHAR(50) REFERENCES users(Email),
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW()
);