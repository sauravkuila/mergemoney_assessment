-- Schema creation script
-- Ensure timezone is set for the target database created by docker-compose
ALTER DATABASE mergemoney SET TIME ZONE 'Asia/Kolkata';

-- Drop tables if they exist
DROP TABLE IF EXISTS user_ref;
DROP TABLE IF EXISTS user_accounts;

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
    user_id VARCHAR(255) NOT NULL,
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

-- create indexes
CREATE INDEX IF NOT EXISTS idx_user_accounts_user_id ON user_accounts (user_id);
CREATE INDEX IF NOT EXISTS idx_user_accounts_account_number ON user_accounts (account_number);
CREATE INDEX IF NOT EXISTS idx_user_accounts_wallet_id ON user_accounts (wallet_id);
CREATE INDEX IF NOT EXISTS idx_user_accounts_upi_id ON user_accounts (upi_id);

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