CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE domains (
    name VARCHAR(255) NOT NULL,
    used_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (name)
);

CREATE TABLE domain_categories (
    domain_name VARCHAR(255) NOT NULL,
    category_id INT NOT NULL,
    FOREIGN KEY (domain_name) REFERENCES domains(name),
    FOREIGN KEY (category_id) REFERENCES categories(id)
);
