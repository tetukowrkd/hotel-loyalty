-- HOTEL STATUS
CREATE TYPE hotel_status AS ENUM (
  'active',
  'inactive',
  'draft',
  'blocked'
);

-- BOOKING STATUS
CREATE TYPE booking_status AS ENUM (
  'pending',
  'paid',
  'confirmed',
  'cancelled',
  'completed'
);

-- PAYMENT STATUS
CREATE TYPE payment_status AS ENUM (
  'pending',
  'settlement',
  'cancel',
  'expire',
  'failure'
);


-- //// =========================
-- //// 1. IDENTITY
-- //// =========================

CREATE TABLE roles (
  id UUID PRIMARY KEY,
  name VARCHAR UNIQUE NOT NULL,
  description TEXT,
  created_at TIMESTAMP,
  updated_at TIMESTAMP
);

CREATE TABLE loyalty_tiers (
  id UUID PRIMARY KEY,
  name VARCHAR,
  min_points INT,
  benefits_json TEXT
);

CREATE TABLE users (
  id UUID PRIMARY KEY,
  name VARCHAR,
  email VARCHAR UNIQUE,
  password VARCHAR,
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
  is_delete BOOLEAN DEFAULT FALSE,
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

-- //// =========================
-- //// 2. HOTEL
-- //// =========================

CREATE TABLE hotels (
  id UUID PRIMARY KEY,
  name VARCHAR,
  description TEXT,
  address TEXT,
  city VARCHAR,
  country VARCHAR,
  lat DECIMAL,
  lng DECIMAL,
  star_rating INT,
  status hotel_status,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL
);

CREATE TABLE hotel_images (
  id UUID PRIMARY KEY,
  hotel_id UUID,
  image_url TEXT,
  is_primary BOOLEAN,
  sort_order INT,

  CONSTRAINT fk_hotel_images FOREIGN KEY (hotel_id) REFERENCES hotels(id) ON DELETE CASCADE
);

CREATE TABLE rooms (
  id UUID PRIMARY KEY,
  hotel_id UUID,
  name VARCHAR,
  description TEXT,
  capacity INT,
  base_price DECIMAL,

  CONSTRAINT fk_rooms_hotel FOREIGN KEY (hotel_id) REFERENCES hotels(id) ON DELETE CASCADE
);

CREATE TABLE room_images (
  id UUID PRIMARY KEY,
  room_id UUID,
  image_url TEXT,
  is_primary BOOLEAN,

  CONSTRAINT fk_room_images FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE
);

CREATE TABLE hotel_facilities (
  id UUID PRIMARY KEY,
  name VARCHAR,
  icon VARCHAR
);

CREATE TABLE hotel_facility_maps (
  hotel_id UUID,
  facility_id UUID,
  PRIMARY KEY (hotel_id, facility_id),

  CONSTRAINT fk_hfm_hotel FOREIGN KEY (hotel_id) REFERENCES hotels(id) ON DELETE CASCADE,
  CONSTRAINT fk_hfm_facility FOREIGN KEY (facility_id) REFERENCES hotel_facilities(id)
);

-- //// =========================
-- //// 3. INVENTORY
-- //// =========================

CREATE TABLE room_inventory (
  id UUID PRIMARY KEY,
  room_id UUID,
  date DATE,
  available_stock INT,

  CONSTRAINT fk_inventory_room FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE
);

CREATE TABLE room_prices (
  id UUID PRIMARY KEY,
  room_id UUID,
  date DATE,
  price DECIMAL,

  CONSTRAINT fk_prices_room FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE
);

-- //// =========================
-- //// 4. BOOKING
-- //// =========================

CREATE TABLE bookings (
  id UUID PRIMARY KEY,
  user_id UUID,
  status booking_status,
  check_in DATE,
  check_out DATE,
  total_amount DECIMAL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL,

  CONSTRAINT fk_booking_user FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE booking_items (
  id UUID PRIMARY KEY,
  booking_id UUID,
  room_id UUID,
  qty INT,
  price_per_night DECIMAL,

  CONSTRAINT fk_booking_items_booking FOREIGN KEY (booking_id) REFERENCES bookings(id) ON DELETE CASCADE,
  CONSTRAINT fk_booking_items_room FOREIGN KEY (room_id) REFERENCES rooms(id)
);

-- //// =========================
-- //// 5. PAYMENT
-- //// =========================

CREATE TABLE payments (
  id UUID PRIMARY KEY,
  booking_id UUID,

  external_id VARCHAR UNIQUE NOT NULL,
  transaction_id VARCHAR,

  amount DECIMAL,

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

CREATE TABLE payment_logs (
  id UUID PRIMARY KEY,
  payment_id UUID,
  payload TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL,

  CONSTRAINT fk_payment_logs FOREIGN KEY (payment_id) REFERENCES payments(id) ON DELETE CASCADE
);

-- //// =========================
-- //// 6. LOYALTY
-- //// =========================

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

-- //// =========================
-- //// 7. REVIEW
-- //// =========================

CREATE TABLE reviews (
  id UUID PRIMARY KEY,
  user_id UUID,
  hotel_id UUID,
  rating INT,
  comment TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL,

  CONSTRAINT fk_reviews_user FOREIGN KEY (user_id) REFERENCES users(id),
  CONSTRAINT fk_reviews_hotel FOREIGN KEY (hotel_id) REFERENCES hotels(id)
);

-- //// =========================
-- //// 8. PROMO
-- //// =========================

CREATE TABLE promotions (
  id UUID PRIMARY KEY,
  name VARCHAR,
  type VARCHAR,
  value DECIMAL,
  start_date DATE,
  end_date DATE
);

-- //// =========================
-- //// 9. WISHLIST
-- //// =========================

CREATE TABLE wishlists (
  user_id UUID,
  hotel_id UUID,
  PRIMARY KEY (user_id, hotel_id),

  CONSTRAINT fk_wishlist_user FOREIGN KEY (user_id) REFERENCES users(id),
  CONSTRAINT fk_wishlist_hotel FOREIGN KEY (hotel_id) REFERENCES hotels(id)
);

-- //// =========================
-- //// INDEXES
-- //// =========================

CREATE INDEX idx_bookings_user_id ON bookings(user_id);

CREATE INDEX idx_hotel_images_hotel_id ON hotel_images(hotel_id);
CREATE INDEX idx_rooms_hotel_id ON rooms(hotel_id);

CREATE INDEX idx_room_inventory_room_date ON room_inventory(room_id, date);
CREATE INDEX idx_room_prices_room_date ON room_prices(room_id, date);

CREATE INDEX idx_payments_booking_id ON payments(booking_id);

CREATE INDEX idx_payment_logs_payment_id ON payment_logs(payment_id);