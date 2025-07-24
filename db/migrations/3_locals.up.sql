-- Create locals table for local businesses in Vistara Backend
-- This table stores information about local businesses and places
CREATE TABLE locals (
    id UUID PRIMARY KEY,
    name VARCHAR NOT NULL,
    description VARCHAR NOT NULL,
    address VARCHAR NOT NULL,
    city VARCHAR NOT NULL,
    province VARCHAR NOT NULL,
    longitude DECIMAL(9,6) NOT NULL,
    latitude DECIMAL(9,6) NOT NULL,
    label VARCHAR NOT NULL,
    opened_time VARCHAR NOT NULL,
    photo_url TEXT NOT NULL,
    is_business BOOL NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);