CREATE TABLE teams (
    id SERIAL PRIMARY KEY,
    name VARCHAR(80) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP
);