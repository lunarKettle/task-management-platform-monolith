CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    description TEXT NOT NULL,
    employee_id INT NOT NULL,
    project_id INT NOT NULL,
    CONSTRAINT fk_employee FOREIGN KEY (employee_id) REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT fk_project FOREIGN KEY (project_id) REFERENCES projects (id) ON DELETE CASCADE
);