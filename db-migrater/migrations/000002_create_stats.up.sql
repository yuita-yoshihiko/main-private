CREATE TABLE stats (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    total_items INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
