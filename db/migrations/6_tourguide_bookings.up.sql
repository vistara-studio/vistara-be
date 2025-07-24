-- Create tourguide_bookings table for tour guide booking system in Vistara Backend
-- This table manages tour guide bookings and reviews for tourist attractions
CREATE TABLE tourguide_bookings (
    id UUID PRIMARY KEY,
    payment_url VARCHAR NOT NULL,
    star INT CHECK (star >= 1 AND star <= 5),
    content TEXT,
    photo_url TEXT,
    booked_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    status VARCHAR NOT NULL DEFAULT 'pending',
    user_id UUID REFERENCES users (id) ON DELETE CASCADE,
    tourist_attraction_id UUID REFERENCES tourist_attractions (id) ON DELETE CASCADE
);