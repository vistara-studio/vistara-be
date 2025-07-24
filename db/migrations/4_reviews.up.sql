-- Create reviews table for local business reviews in Vistara Backend
-- This table stores user reviews and ratings for local businesses
CREATE TABLE reviews (
    id UUID PRIMARY KEY,
    star INT NOT NULL CHECK (star >= 1 AND star <= 5),
    content VARCHAR NOT NULL,
    photo_url TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    user_id UUID REFERENCES users (id) ON DELETE CASCADE,
    local_id UUID REFERENCES locals (id) ON DELETE CASCADE
);