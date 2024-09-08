
CREATE TABLE cashback (
    id UUID PRIMARY KEY NOT NULL,
    merchant_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    percentage FLOAT NOT NULL,
    eligible_products  UUID[] NOT NULL,
    active BOOLEAN NOT NULL,
    expiration TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);


CREATE TABLE cashback_user (
      merchant_id UUID NOT NULL,
    cashback_id UUID NOT NULL,
    user_id UUID NOT NULL,
    points FLOAT NOT NULL,
    PRIMARY KEY (cashback_id, user_id)
    
);
