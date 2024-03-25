CREATE TABLE IF NOT EXISTS balance_histories (
    id UUID PRIMARY KEY NOT NULL,
    balance BIGINT NOT NULL,
    bank_number VARCHAR(40) NOT NULL,
    bank_name VARCHAR(40) NOT NULL,
    currency VARCHAR(40) NOT NULL,
    receipt VARCHAR NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    user_id UUID REFERENCES users(id) NOT NULL
);
