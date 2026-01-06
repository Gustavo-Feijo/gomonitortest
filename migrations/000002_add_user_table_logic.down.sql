DROP TRIGGER IF EXISTS update_users_updated_at ON users;

DROP FUNCTION IF EXISTS update_updated_at_column ();

DROP INDEX IF EXISTS idx_users_email;

ALTER TABLE users
DROP COLUMN IF EXISTS user_name,
DROP COLUMN IF EXISTS email,
DROP COLUMN IF EXISTS password,
DROP COLUMN IF EXISTS user_role,
DROP COLUMN IF EXISTS created_at,
DROP COLUMN IF EXISTS updated_at;

DROP TYPE IF EXISTS user_role;