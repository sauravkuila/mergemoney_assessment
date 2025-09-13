-- Schema creation script
ALTER DATABASE your_database SET TIME ZONE 'Asia/Kolkata';

-- Drop tables if they exist
DROP TABLE IF EXISTS user_ref;

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

-- create indexes

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