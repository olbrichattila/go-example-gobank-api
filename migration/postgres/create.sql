CREATE TABLE IF NOT EXISTS account (
    id SERIAL PRIMARY KEY,
    email VARCHAR(128),
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    encrypted_password VARCHAR(100),
    account_number SERIAL,
    balance SERIAL,
    created_at timestamp
)