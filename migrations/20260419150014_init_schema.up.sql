-- =========================
-- ENUMS
-- =========================

CREATE TYPE hotel_status AS ENUM ('draft','published');

CREATE TYPE booking_status AS ENUM (
  'pending','paid','confirmed','cancelled','completed'
);

CREATE TYPE payment_status AS ENUM (
  'pending','settlement','cancel','expire','failure'
);

-- =========================
-- IDENTITY
-- =========================

CREATE TABLE roles (
  id UUID PRIMARY KEY,
  name VARCHAR UNIQUE NOT NULL,
  description TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL
);

CREATE TABLE loyalty_tiers (
  id UUID PRIMARY KEY,
  name VARCHAR NOT NULL,
  min_points INT DEFAULT 0,
  benefits_json JSONB,
  deleted_at TIMESTAMP NULL
);

CREATE TABLE users (
  id UUID PRIMARY KEY,
  name VARCHAR NOT NULL,
  email VARCHAR UNIQUE NOT NULL,
  password VARCHAR NOT NULL,
  phone VARCHAR,
  member_code VARCHAR(20) UNIQUE NOT NULL,

  loyalty_points INT DEFAULT 0,
  tier_id UUID,
  role_id UUID,

  email_verified_at TIMESTAMP,
  last_login_at TIMESTAMP,
  login_attempt INT DEFAULT 0,

  reset_password_token VARCHAR,
  reset_password_expired_at TIMESTAMP,

  is_active BOOLEAN DEFAULT TRUE,

  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL,

  CONSTRAINT fk_users_tier FOREIGN KEY (tier_id) REFERENCES loyalty_tiers(id),
  CONSTRAINT fk_users_role FOREIGN KEY (role_id) REFERENCES roles(id)
);

CREATE TABLE user_tokens (
  id UUID PRIMARY KEY,
  user_id UUID,
  refresh_token TEXT UNIQUE NOT NULL,
  user_agent TEXT,
  ip_address VARCHAR,
  expires_at TIMESTAMP,

  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL,

  CONSTRAINT fk_user_tokens_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_user_tokens_user_id ON user_tokens(user_id);

-- =========================
-- HOTEL
-- =========================

CREATE TABLE hotels (
  id UUID PRIMARY KEY,
  name VARCHAR NOT NULL,
  description TEXT,

  address TEXT NOT NULL,
  city VARCHAR,
  country VARCHAR,

  lat DECIMAL(10,6),
  lng DECIMAL(10,6),

  star_rating INT CHECK (star_rating BETWEEN 1 AND 5),

  status hotel_status NOT NULL DEFAULT 'draft',
  is_active BOOLEAN DEFAULT TRUE,

  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL
);

CREATE INDEX idx_hotels_city ON hotels(city);
CREATE INDEX idx_hotels_status ON hotels(status);
CREATE INDEX idx_hotels_deleted_at ON hotels(deleted_at);

-- =========================
-- HOTEL IMAGES
-- =========================

CREATE TABLE hotel_images (
  id UUID PRIMARY KEY,
  hotel_id UUID,
  image_url TEXT NOT NULL,
  is_primary BOOLEAN DEFAULT FALSE,
  sort_order INT DEFAULT 0,
  deleted_at TIMESTAMP NULL,

  CONSTRAINT fk_hotel_images FOREIGN KEY (hotel_id) REFERENCES hotels(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX unique_primary_image
ON hotel_images (hotel_id)
WHERE is_primary = true;

CREATE INDEX idx_hotel_images_hotel_id ON hotel_images(hotel_id);
CREATE INDEX idx_hotel_images_deleted_at ON hotel_images(deleted_at);

-- =========================
-- ROOMS
-- =========================

CREATE TABLE rooms (
  id UUID PRIMARY KEY,
  hotel_id UUID,
  name VARCHAR NOT NULL,
  description TEXT,
  capacity INT CHECK (capacity > 0),
  base_price DECIMAL(12,2) NOT NULL,
  deleted_at TIMESTAMP NULL,

  CONSTRAINT fk_rooms_hotel FOREIGN KEY (hotel_id) REFERENCES hotels(id) ON DELETE CASCADE
);

CREATE INDEX idx_rooms_hotel_id ON rooms(hotel_id);
CREATE INDEX idx_rooms_deleted_at ON rooms(deleted_at);

CREATE TABLE room_images (
  id UUID PRIMARY KEY,
  room_id UUID,
  image_url TEXT NOT NULL,
  is_primary BOOLEAN DEFAULT FALSE,
  deleted_at TIMESTAMP NULL,

  CONSTRAINT fk_room_images FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX unique_room_primary
ON room_images (room_id)
WHERE is_primary = true;

CREATE INDEX idx_room_images_deleted_at ON room_images(deleted_at);

-- =========================
-- FACILITIES
-- =========================

CREATE TABLE hotel_facilities (
  id UUID PRIMARY KEY,
  name VARCHAR NOT NULL,
  icon VARCHAR NOT NULL,
  deleted_at TIMESTAMP NULL
);

CREATE UNIQUE INDEX unique_facility_name
ON hotel_facilities (LOWER(name));

CREATE INDEX idx_hotel_facilities_deleted_at ON hotel_facilities(deleted_at);

CREATE TABLE hotel_facility_maps (
  hotel_id UUID,
  facility_id UUID,
  deleted_at TIMESTAMP NULL,
  PRIMARY KEY (hotel_id, facility_id),

  CONSTRAINT fk_hfm_hotel FOREIGN KEY (hotel_id) REFERENCES hotels(id) ON DELETE CASCADE,
  CONSTRAINT fk_hfm_facility FOREIGN KEY (facility_id) REFERENCES hotel_facilities(id)
);

CREATE INDEX idx_hfm_deleted_at ON hotel_facility_maps(deleted_at);

-- =========================
-- INVENTORY & PRICING
-- =========================

CREATE TABLE room_inventory (
  id UUID PRIMARY KEY,
  room_id UUID,
  date DATE,
  available_stock INT CHECK (available_stock >= 0),
  deleted_at TIMESTAMP NULL,

  CONSTRAINT fk_inventory_room FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX unique_inventory
ON room_inventory (room_id, date);

CREATE INDEX idx_room_inventory_room_date ON room_inventory(room_id, date);
CREATE INDEX idx_room_inventory_deleted_at ON room_inventory(deleted_at);

CREATE TABLE room_prices (
  id UUID PRIMARY KEY,
  room_id UUID,
  date DATE,
  price DECIMAL(12,2),
  deleted_at TIMESTAMP NULL,

  CONSTRAINT fk_prices_room FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX unique_price
ON room_prices (room_id, date);

CREATE INDEX idx_room_prices_room_date ON room_prices(room_id, date);
CREATE INDEX idx_room_prices_deleted_at ON room_prices(deleted_at);

-- =========================
-- BOOKING
-- =========================

CREATE TABLE bookings (
  id UUID PRIMARY KEY,
  user_id UUID,
  status booking_status,

  check_in DATE,
  check_out DATE,
  total_amount DECIMAL(12,2),

  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL,

  CONSTRAINT fk_booking_user FOREIGN KEY (user_id) REFERENCES users(id),
  CONSTRAINT check_date_valid CHECK (check_out > check_in)
);

CREATE INDEX idx_bookings_user_id ON bookings(user_id);

CREATE TABLE booking_items (
  id UUID PRIMARY KEY,
  booking_id UUID,
  room_id UUID,
  qty INT CHECK (qty > 0),
  price_per_night DECIMAL(12,2),
  total_price DECIMAL(12,2),
  deleted_at TIMESTAMP NULL,

  CONSTRAINT fk_booking_items_booking FOREIGN KEY (booking_id) REFERENCES bookings(id) ON DELETE CASCADE,
  CONSTRAINT fk_booking_items_room FOREIGN KEY (room_id) REFERENCES rooms(id)
);

CREATE INDEX idx_booking_items_deleted_at ON booking_items(deleted_at);

-- =========================
-- PAYMENT
-- =========================

CREATE TABLE payments (
  id UUID PRIMARY KEY,
  booking_id UUID,
  external_id VARCHAR UNIQUE NOT NULL,
  transaction_id VARCHAR,
  amount DECIMAL(12,2),

  payment_type VARCHAR,
  method VARCHAR,

  status payment_status,
  fraud_status VARCHAR,
  raw_response TEXT,

  paid_at TIMESTAMP,
  expired_at TIMESTAMP,

  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL,

  CONSTRAINT fk_payment_booking FOREIGN KEY (booking_id) REFERENCES bookings(id)
);

CREATE INDEX idx_payments_booking_id ON payments(booking_id);

CREATE TABLE payment_logs (
  id UUID PRIMARY KEY,
  payment_id UUID,
  payload TEXT,

  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL,

  CONSTRAINT fk_payment_logs FOREIGN KEY (payment_id) REFERENCES payments(id) ON DELETE CASCADE
);

CREATE INDEX idx_payment_logs_payment_id ON payment_logs(payment_id);

-- =========================
-- LOYALTY
-- =========================

CREATE TABLE loyalty_transactions (
  id UUID PRIMARY KEY,
  user_id UUID,
  points INT,
  type VARCHAR,
  reference_type VARCHAR,
  reference_id UUID,

  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL,

  CONSTRAINT fk_loyalty_user FOREIGN KEY (user_id) REFERENCES users(id)
);

-- =========================
-- REVIEW
-- =========================

CREATE TABLE reviews (
  id UUID PRIMARY KEY,
  user_id UUID,
  hotel_id UUID,

  rating INT CHECK (rating BETWEEN 1 AND 5),
  comment TEXT,

  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL,

  CONSTRAINT fk_reviews_user FOREIGN KEY (user_id) REFERENCES users(id),
  CONSTRAINT fk_reviews_hotel FOREIGN KEY (hotel_id) REFERENCES hotels(id)
);

-- =========================
-- PROMO
-- =========================

CREATE TABLE promotions (
  id UUID PRIMARY KEY,
  name VARCHAR,
  type VARCHAR,
  value DECIMAL(12,2),
  start_date DATE,
  end_date DATE,
  deleted_at TIMESTAMP NULL
);

CREATE INDEX idx_promotions_deleted_at ON promotions(deleted_at);

-- =========================
-- WISHLIST
-- =========================

CREATE TABLE wishlists (
  user_id UUID,
  hotel_id UUID,
  deleted_at TIMESTAMP NULL,

  PRIMARY KEY (user_id, hotel_id),

  CONSTRAINT fk_wishlist_user FOREIGN KEY (user_id) REFERENCES users(id),
  CONSTRAINT fk_wishlist_hotel FOREIGN KEY (hotel_id) REFERENCES hotels(id)
);

CREATE INDEX idx_wishlists_deleted_at ON wishlists(deleted_at);