-- Create tourist_attractions table for Vistara Backend
-- This table stores information about tourist attractions and tour guide services
CREATE TABLE tourist_attractions (
    id UUID PRIMARY KEY,
    name VARCHAR NOT NULL,
    description TEXT NOT NULL,
    address VARCHAR NOT NULL,
    city VARCHAR NOT NULL,
    province VARCHAR NOT NULL,
    longitude DECIMAL(9,6) NOT NULL,
    latitude DECIMAL(9,6) NOT NULL,
    photo_url TEXT NOT NULL,
    tour_guide_price BIGINT NOT NULL,
    tour_guide_count INT NOT NULL,
    tour_guide_discount_percentage NUMERIC(5,2) NOT NULL,
    price BIGINT NOT NULL,
    discount_percentage NUMERIC(5,2) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);