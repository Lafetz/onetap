CREATE TABLE tier_level (
    tier_id UUID PRIMARY KEY NOT NULL,  
    merchant_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    min_points INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE (merchant_id, name)  
);
