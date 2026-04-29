DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP INDEX IF EXISTS idx_users_role;
DROP INDEX IF EXISTS idx_users_is_blocked;
DROP TABLE IF EXISTS users CASCADE;
