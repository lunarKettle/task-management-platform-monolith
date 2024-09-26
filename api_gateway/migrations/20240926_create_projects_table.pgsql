CREATE TABLE projects (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    start_date TIMESTAMP DEFAULT NOW(),
    planned_end_date TIMESTAMP,
    actual_end_date TIMESTAMP,
    status VARCHAR(15),
    priority SMALLINT,
    manager_id INT,
    budget REAL
);