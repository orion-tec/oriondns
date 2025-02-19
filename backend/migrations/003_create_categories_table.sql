CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE domains (
    domain VARCHAR(255) NOT NULL,
    used_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT domains_uq UNIQUE (domain),
    PRIMARY KEY (domain)
);

CREATE TABLE domain_categories (
    domain VARCHAR(255) NOT NULL,
    category_id INT NOT NULL,
    FOREIGN KEY (domain) REFERENCES domains(domain),
    FOREIGN KEY (category_id) REFERENCES categories(id)
);
