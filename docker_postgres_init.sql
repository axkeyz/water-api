-- Enable PostGIS (as of 3.0 contains just geometry/geography)
-- CREATE EXTENSION postgis;

SET TIME ZONE 'Pacific/Auckland';
SET SESSION TIME ZONE 'NZDT';

-- Auto update updated_at timestamp
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create outage table
CREATE TABLE "outage" (
  id SERIAL PRIMARY KEY,
  outage_id INT NOT NULL UNIQUE,
  street VARCHAR(256),
  suburb VARCHAR(256),
  location POINT NOT NULL,
  start_date TIMESTAMP WITHOUT TIME ZONE,
  end_date TIMESTAMP WITHOUT TIME ZONE,
  outage_type VARCHAR(50),
  created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Auto-update
CREATE TRIGGER set_timestamp BEFORE UPDATE ON outage
FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();