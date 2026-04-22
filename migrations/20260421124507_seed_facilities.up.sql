-- 🔥 Seed hotel facilities (static UUID)

INSERT INTO hotel_facilities (id, name, icon) VALUES
('11111111-1111-1111-1111-111111111111', 'WiFi', 'wifi'),
('22222222-2222-2222-2222-222222222222', 'Swimming Pool', 'pool'),
('33333333-3333-3333-3333-333333333333', 'Gym', 'gym'),
('44444444-4444-4444-4444-444444444444', 'Parking', 'parking'),
('55555555-5555-5555-5555-555555555555', 'Restaurant', 'restaurant'),
('66666666-6666-6666-6666-666666666666', '24h Front Desk', 'reception'),
('77777777-7777-7777-7777-777777777777', 'Air Conditioning', 'ac'),
('88888888-8888-8888-8888-888888888888', 'Elevator', 'elevator'),
('99999999-9999-9999-9999-999999999999', 'Spa', 'spa'),
('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'Laundry', 'laundry'),
('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'Room Service', 'room-service'),
('cccccccc-cccc-cccc-cccc-cccccccccccc', 'Meeting Room', 'meeting'),
('dddddddd-dddd-dddd-dddd-dddddddddddd', 'Airport Shuttle', 'shuttle'),
('eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee', 'Bar', 'bar'),
('ffffffff-ffff-ffff-ffff-ffffffffffff', 'Breakfast', 'breakfast')
ON CONFLICT (id) DO NOTHING;