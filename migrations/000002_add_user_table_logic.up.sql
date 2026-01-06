CREATE TYPE user_role AS ENUM('admin', 'user', 'moderator');

ALTER TABLE users
ADD COLUMN user_name VARCHAR,
ADD COLUMN email VARCHAR(254) NOT NULL,
ADD COLUMN password CHAR(60) NOT NULL,
ADD COLUMN role user_role NOT NULL DEFAULT 'user',
ADD COLUMN created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
ADD COLUMN updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW();

CREATE UNIQUE INDEX idx_users_email ON users (email);

-- Trigger to set the updated_at.
CREATE
OR REPLACE FUNCTION update_updated_at_column () RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Add trigger to users table.
CREATE TRIGGER update_users_updated_at BEFORE
UPDATE ON users FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column ();