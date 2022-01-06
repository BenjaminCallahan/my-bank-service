CREATE TABLE IF NOT EXISTS currencies (
    id INTEGER PRIMARY KEY autoincrement,
    name VARCHAR(10) NOT NULL
);

CREATE TABLE IF NOT EXISTS currency_rate (
    source_currency_id INTEGER,
    target_currency_id INTEGER,
    exchange_rate DECIMAL(6, 4) NOT NULL,
    FOREIGN KEY(source_currency_id) REFERENCES currencies(id),
    FOREIGN key(target_currency_id) REFERENCES currencies(id)
);

CREATE TABLE IF NOT EXISTS accounts (
    id INTEGER PRIMARY KEY autoincrement,
    balance DECIMAL(32, 2) NOT NULL,
    currency_id INTEGER,
    FOREIGN KEY(currency_id) REFERENCES currencies(id)
);

CREATE TABLE IF NOT EXISTS transfers (
    id INTEGER PRIMARY KEY autoincrement,
    to_account_id INTEGER,
    amount DECIMAL(32, 2),
    is_processed INTEGER NOT NULL DEFAULT 0,
    FOREIGN KEY(to_account_id) REFERENCES accounts(id)
);

INSERT INTO
    currencies (name)
VALUES
    ('SBP'),
    ('RUB');

INSERT INTO
    currency_rate (
        source_currency_id,
        target_currency_id,
        exchange_rate
    )
VALUES
    (1, 2, 0.7523);

INSERT INTO
    accounts (balance, currency_id)
VALUES
    (0, 1);