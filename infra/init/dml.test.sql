-- dml.test.sql
USE `campfinderdb`;

-- Userデータのセットアップ
INSERT INTO User (id, name, email, password) VALUES
('5fe0e237-6b49-11ee-b686-0242c0a87001', 'test', 'test@gmail.com', 'password123');

-- Spotデータのセットアップ
INSERT INTO Spot (id, category, name, address, lat, lng, period, phone, price, description, iconpath) VALUES
('5c5323e9-c78f-4dac-94ef-d34ab5ea8fed', 'campsite', '旭川市21世紀の森ふれあい広場', '北海道旭川市東旭川町瑞穂4288', 43.7172721, 142.6674615, '2022年5月1日(日)〜11月30日(水)', '0166-76-2108', '有料。ログハウス大人290円〜750円、高校生以下180〜460円', '旭川市21世紀の森ふれあい広場は、ペーパンダムの周辺に整備された多目的公園、旭川市21世紀の森に隣接するキャンプ場です。', '/static/img/campsiteflag.jpeg');

-- Imageデータのセットアップ
INSERT INTO Image (id, spot_id, user_id, url, created) VALUES
('31894386-3e60-45a8-bc67-f46b72b42554', '5c5323e9-c78f-4dac-94ef-d34ab5ea8fed', '5fe0e237-6b49-11ee-b686-0242c0a87001', 'https://lh3.googleusercontent.com/places/ABCD', CURRENT_TIMESTAMP);