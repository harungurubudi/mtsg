-- +migrate Down
-- Drop users table and all related objects

-- Drop trigger first
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- Drop trigger function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop indexes
DROP INDEX IF EXISTS idx_users_tenant_status;
DROP INDEX IF EXISTS idx_users_is_deleted;
DROP INDEX IF EXISTS idx_users_status;
DROP INDEX IF EXISTS idx_users_tenant_id;
DROP INDEX IF EXISTS idx_users_email;

-- Drop foreign key constraint (if it was added)
-- ALTER TABLE users DROP CONSTRAINT IF EXISTS fk_users_tenant_id;

-- Drop the users table
DROP TABLE IF EXISTS users; 