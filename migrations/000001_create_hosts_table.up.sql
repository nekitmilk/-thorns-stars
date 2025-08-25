CREATE TABLE hosts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE,
    ip VARCHAR(45) NOT NULL UNIQUE,
    priority INTEGER NOT NULL DEFAULT 1 CHECK (
        priority BETWEEN 1 AND 100
    ),
    status VARCHAR(20) NOT NULL DEFAULT 'unknown' CHECK (status IN ('online', 'offline', 'unknown')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_hosts_name ON hosts(name);
CREATE INDEX idx_hosts_ip ON hosts(ip);
CREATE INDEX idx_hosts_status ON hosts(status);
CREATE INDEX idx_hosts_priority ON hosts(priority);
INSERT INTO hosts (
        name,
        ip,
        priority,
        status,
        created_at,
        updated_at
    )
VALUES (
        'master-node-1',
        '192.168.1.100',
        100,
        'online',
        CURRENT_TIMESTAMP - INTERVAL '5 days',
        CURRENT_TIMESTAMP - INTERVAL '1 hour'
    ),
    (
        'agent-node-1',
        '192.168.1.101',
        50,
        'online',
        CURRENT_TIMESTAMP - INTERVAL '4 days',
        CURRENT_TIMESTAMP - INTERVAL '30 minutes'
    ),
    (
        'agent-node-2',
        '192.168.1.102',
        75,
        'offline',
        CURRENT_TIMESTAMP - INTERVAL '3 days',
        CURRENT_TIMESTAMP - INTERVAL '2 hours'
    ),
    (
        'agent-node-3',
        '192.168.1.103',
        25,
        'online',
        CURRENT_TIMESTAMP - INTERVAL '2 days',
        CURRENT_TIMESTAMP - INTERVAL '15 minutes'
    ),
    (
        'agent-node-4',
        '192.168.1.104',
        10,
        'unknown',
        CURRENT_TIMESTAMP - INTERVAL '1 day',
        CURRENT_TIMESTAMP - INTERVAL '5 minutes'
    ),
    (
        'backup-master',
        '192.168.1.105',
        90,
        'online',
        CURRENT_TIMESTAMP - INTERVAL '6 days',
        CURRENT_TIMESTAMP - INTERVAL '45 minutes'
    ),
    (
        'monitoring-node',
        '192.168.1.106',
        30,
        'offline',
        CURRENT_TIMESTAMP - INTERVAL '7 days',
        CURRENT_TIMESTAMP - INTERVAL '3 hours'
    );