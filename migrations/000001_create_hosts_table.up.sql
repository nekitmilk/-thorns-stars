CREATE TABLE hosts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE,
    ip VARCHAR(45) NOT NULL UNIQUE,
    priority INTEGER NOT NULL DEFAULT 1 CHECK (priority BETWEEN 1 AND 100),
    status VARCHAR(20) NOT NULL DEFAULT 'unknown' CHECK (status IN ('online', 'offline', 'unknown')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_hosts_name ON hosts(name);
CREATE INDEX idx_hosts_ip ON hosts(ip);
CREATE INDEX idx_hosts_status ON hosts(status);
CREATE INDEX idx_hosts_priority ON hosts(priority);