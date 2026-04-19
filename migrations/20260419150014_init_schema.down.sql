//// =========================
//// DROP TABLES (REVERSE ORDER)
//// =========================

DROP TABLE IF EXISTS wishlists;

DROP TABLE IF EXISTS promotions;

DROP TABLE IF EXISTS reviews;

DROP TABLE IF EXISTS loyalty_transactions;
DROP TABLE IF EXISTS loyalty_tiers;

DROP TABLE IF EXISTS payment_logs;
DROP TABLE IF EXISTS payments;

DROP TABLE IF EXISTS booking_items;
DROP TABLE IF EXISTS bookings;

DROP TABLE IF EXISTS room_prices;
DROP TABLE IF EXISTS room_inventory;

DROP TABLE IF EXISTS hotel_facility_maps;
DROP TABLE IF EXISTS hotel_facilities;

DROP TABLE IF EXISTS room_images;
DROP TABLE IF EXISTS rooms;

DROP TABLE IF EXISTS hotel_images;
DROP TABLE IF EXISTS hotels;

DROP TABLE IF EXISTS user_tokens;
DROP TABLE IF EXISTS users;

DROP TABLE IF EXISTS roles;

//// =========================
//// DROP ENUM TYPES
//// =========================

DROP TYPE IF EXISTS payment_status;
DROP TYPE IF EXISTS booking_status;
DROP TYPE IF EXISTS hotel_status;