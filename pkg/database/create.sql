-- Schema creation script
-- Ensure timezone is set for the target database created by docker-compose
ALTER DATABASE mergemoney SET TIME ZONE 'Asia/Kolkata';

-- Drop tables if they exist
DROP TABLE IF EXISTS user_ref;
DROP TABLE IF EXISTS user_accounts;
DROP TABLE IF EXISTS transfers;

-- Table: user_ref
CREATE TABLE user_ref (
    user_id VARCHAR(255) PRIMARY KEY,
    user_name VARCHAR(255) NOT NULL,
    mobile VARCHAR(15) NOT NULL,
    country_code VARCHAR(10) NOT NULL,
    user_role VARCHAR(50) NOT NULL,
    user_mpin VARCHAR(10) DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Table: user_accounts
CREATE TABLE user_accounts (
    serial_id SERIAL PRIMARY KEY,
    user_id VARCHAR(255) REFERENCES user_ref(user_id),
    type VARCHAR(50) NOT NULL,  /* e.g., 'bank', 'wallet', 'upi' */
    bank_name VARCHAR(255),
    account_number VARCHAR(50),
    ifsc VARCHAR(50),
    linked_via VARCHAR(50) NOT NULL,
    wallet_name VARCHAR(255),
    wallet_id VARCHAR(255),
    upi_id VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    CONSTRAINT fk_user_accounts_user FOREIGN KEY (user_id)
        REFERENCES user_ref(user_id)
        ON UPDATE CASCADE
        ON DELETE CASCADE
);

-- Table: orders
CREATE TABLE orders (
    order_id VARCHAR(64) PRIMARY KEY,
    user_id VARCHAR(255) REFERENCES user_ref(user_id),
    source_sid BIGINT REFERENCES user_accounts(serial_id),
    source_currency VARCHAR(10) NOT NULL,
    source_amount NUMERIC(18,4) NOT NULL,
    destination_currency VARCHAR(10) NOT NULL,
    destination_amount NUMERIC(18,4) NOT NULL,
    conversion_rate NUMERIC(18,6),
    conversion_rate_date DATE,
    order_status VARCHAR(30) DEFAULT 'created',  -- lifecycle: created → inprogress → [failed | completed]
    remarks VARCHAR(255),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Table: order_destinations
CREATE TABLE order_destinations (
    destination_id BIGSERIAL PRIMARY KEY,
    order_id VARCHAR(64) NOT NULL REFERENCES orders(order_id) ON DELETE CASCADE,
    destination_type VARCHAR(50) NOT NULL, -- 'wallet', 'upi', 'netbanking'
    wallet_id VARCHAR(100), -- for wallet transfers
    upi_id VARCHAR(100), -- for UPI
    bank_account_number VARCHAR(50), -- for netbanking
    ifsc_code VARCHAR(20), -- for netbanking
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (order_id) -- ensure one destination per order
);

-- Table: transactions
CREATE TABLE transactions (
    transaction_id VARCHAR(64) PRIMARY KEY, -- also used as idempotency key for retries
    order_id VARCHAR(64) REFERENCES orders(order_id) ON DELETE CASCADE,
    provider VARCHAR(255) NOT NULL,
    provider_id VARCHAR(255) NOT NULL,
    provider_request JSONB,
    provider_response JSONB,
    status VARCHAR(30) DEFAULT 'initiated', -- initiated, pending, inprogress, completed, failed
    error_message TEXT,
    retry_count INT DEFAULT 0,
    last_retry_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (provider, provider_id)
);


-- create indexes
CREATE INDEX IF NOT EXISTS idx_user_accounts_user_id ON user_accounts (user_id);
CREATE INDEX IF NOT EXISTS idx_user_accounts_account_number ON user_accounts (account_number);
CREATE INDEX IF NOT EXISTS idx_user_accounts_wallet_id ON user_accounts (wallet_id);
CREATE INDEX IF NOT EXISTS idx_user_accounts_upi_id ON user_accounts (upi_id);

-- CREATE INDEX IF NOT EXISTS idx_transfers_user_id ON transfers (user_id);
-- CREATE INDEX IF NOT EXISTS idx_transfers_idempotency_key ON transfers (external_transfer_key);

CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(order_status);
CREATE INDEX IF NOT EXISTS idx_orders_created_at ON orders(created_at);

CREATE INDEX IF NOT EXISTS idx_destinations_order_id ON order_destinations(order_id);
CREATE INDEX IF NOT EXISTS idx_destinations_type ON order_destinations(destination_type);

CREATE INDEX IF NOT EXISTS idx_transactions_order_id ON transactions(order_id);
CREATE INDEX IF NOT EXISTS idx_transactions_status ON transactions(status);
CREATE INDEX IF NOT EXISTS idx_transactions_created_at ON transactions(created_at);

-- Function to update the updated_at column
CREATE OR REPLACE FUNCTION on_row_update_fn()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Triggers

-- Trigger for user_ref table
CREATE TRIGGER update_user_ref_updated_at
BEFORE UPDATE ON user_ref
FOR EACH ROW
EXECUTE FUNCTION on_row_update_fn();

-- Trigger for user_accounts table to update updated_at
CREATE TRIGGER update_user_accounts_updated_at
BEFORE UPDATE ON user_accounts
FOR EACH ROW
EXECUTE FUNCTION on_row_update_fn();

-- Trigger for orders table to update updated_at
CREATE TRIGGER update_orders_updated_at
BEFORE UPDATE ON orders
FOR EACH ROW
EXECUTE FUNCTION on_row_update_fn();

-- Trigger for transactions table to update updated_at
CREATE TRIGGER update_transactions_updated_at
BEFORE UPDATE ON transactions
FOR EACH ROW
EXECUTE FUNCTION on_row_update_fn();