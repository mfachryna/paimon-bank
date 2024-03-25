CREATE TABLE IF NOT EXISTS balances (
    ID  SERIAL,
    balance BIGINT NOT NULL,
    currency VARCHAR(40) NOT NULL,
    user_id UUID REFERENCES users(id) NOT NULL,
    PRIMARY KEY(currency,user_id)
);
