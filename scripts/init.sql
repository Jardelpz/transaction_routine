SELECT 'starting execution';

CREATE TABLE IF NOT EXISTS operation_types(
    operation_type_id SMALLINT PRIMARY KEY,
    description VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS accounts(
    account_id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    document_hash VARCHAR(64) NOT NULL UNIQUE,
    document_encrypted TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS transactions(
    transaction_id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    account_id INT NOT NULL REFERENCES accounts(account_id),
    operation_type_id SMALLINT NOT NULL REFERENCES operation_types(operation_type_id),
    amount NUMERIC(15,2) NOT NULL,
    event_date TIMESTAMP
);

CREATE TABLE IF NOT EXISTS audit_logs (
    audit_id BIGSERIAL PRIMARY KEY,
    event_type VARCHAR(50) NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    entity_id VARCHAR(50) NOT NULL,
    payload JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO operation_types(operation_type_id, description) VALUES
    (1, 'Normal Purchase'),
    (2, 'Purchase with installments'),
    (3, 'Withdrawal'),
    (4, 'Credit Voucher')
;

SELECT 'ending execution';



