-- Enable pgcrypto for UUID generation
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ENUM types for controlled vocabularies
CREATE TYPE water_type AS ENUM ('tap', 'ro', 'rodi');

-- Tanks table
CREATE TABLE tanks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    room VARCHAR(255),
    rack_location VARCHAR(255),
    volume_liters INT NOT NULL CHECK (volume_liters > 0),
    inventory_number VARCHAR(100) UNIQUE,
    water water_type NOT NULL,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Function to automatically update the updated_at timestamp
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger for the tanks table
CREATE TRIGGER set_timestamp_tanks
BEFORE UPDATE ON tanks
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

-- Insert some sample data
INSERT INTO tanks (name, room, rack_location, volume_liters, inventory_number, water, notes) VALUES
('SPS Reef', 'Fish Room', 'A1', 300, 'T-001', 'rodi', 'Main display tank for SPS corals.'),
('Malawi Community', 'Living Room', 'B2', 240, 'T-002', 'tap', 'Community of Mbuna cichlids.'),
('Shrimp Breeding', 'Fish Room', 'C5', 25, 'T-003', 'ro', 'Crystal Red Shrimp breeding project.');