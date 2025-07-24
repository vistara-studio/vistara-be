-- Create users table for Vistara Backend application
-- This table stores user information including authentication and premium status
CREATE TABLE users (
    id UUID PRIMARY KEY,
    full_name VARCHAR NOT NULL,
    email VARCHAR UNIQUE NOT NULL,
    password VARCHAR,
    auth_provider VARCHAR NOT NULL,
    photo_url TEXT NOT NULL,
    is_premium BOOL NOT NULL DEFAULT FALSE,
    expired_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
