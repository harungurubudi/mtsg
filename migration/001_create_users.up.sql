-- +migrate Up
-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    email VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'member',
    ciphertext VARCHAR(255) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE NULL,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    
    -- Check constraints
    CONSTRAINT chk_user_role CHECK (role IN ('admin', 'member')),
    CONSTRAINT chk_user_status CHECK (status IN ('active', 'inactive'))
);

-- Create indexes
CREATE UNIQUE INDEX idx_users_email ON users(email) WHERE NOT is_deleted;
CREATE INDEX idx_users_tenant_id ON users(tenant_id);
CREATE INDEX idx_users_status ON users(status) WHERE status = 'active';
CREATE INDEX idx_users_is_deleted ON users(is_deleted) WHERE is_deleted = false;
CREATE INDEX idx_users_tenant_status ON users(tenant_id, status);

-- Add foreign key constraint (assuming tenants table exists)
-- ALTER TABLE users ADD CONSTRAINT fk_users_tenant_id FOREIGN KEY (tenant_id) REFERENCES tenants(id); 