CREATE TABLE customer_tier (            
    merchant_id UUID NOT NULL,       
    customer_id UUID NOT NULL,      
    tier_name VARCHAR(255) NOT NULL,        
    points INT NOT NULL,                  
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),  
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(), 
    PRIMARY KEY (merchant_id, customer_id)  
);