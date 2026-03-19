CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE stats (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    total_items INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
