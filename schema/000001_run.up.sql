CREATE TABLE IF NOT EXISTS companies (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS departments (
    id SERIAL PRIMARY KEY,
    company_id INT REFERENCES companies(id) ON DELETE CASCADE NOT NULL,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(50)
);

CREATE TABLE IF NOT EXISTS employees (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    surname VARCHAR(255) NOT NULL,
    phone VARCHAR(50),
    company_id INT REFERENCES companies(id) ON DELETE CASCADE NOT NULL,
    department_id INT REFERENCES departments(id) ON DELETE CASCADE NOT NULL
);

CREATE TABLE IF NOT EXISTS passports (
    id SERIAL PRIMARY KEY,
    employee_id INT REFERENCES employees(id) ON DELETE CASCADE NOT NULL,
    type VARCHAR(50),
    number VARCHAR(50),
    UNIQUE(type, number) 
);