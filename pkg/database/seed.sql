-- Seed data for development/testing
INSERT INTO user_ref (user_id, user_name, mobile, country_code, user_role, user_mpin, created_at, updated_at)
VALUES ('ST2VNFCU', 'Test User', '9999999999', '+91', 'customer', '1222', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (user_id) DO NOTHING;
