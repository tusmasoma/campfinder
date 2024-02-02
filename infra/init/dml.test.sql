-- dml.test.sql
USE `campfinderdb`;

-- Userデータのセットアップ
INSERT INTO User (id, name, email, password) VALUES
('5fe0e237-6b49-11ee-b686-0242c0a87001', 'test', 'test@gmail.com', 'password123');

-- Spotデータのセットアップ
INSERT INTO Spot (id, category, name, address, lat, lng, period, phone, price, description, iconpath) VALUES
('5c5323e9-c78f-4dac-94ef-d34ab5ea8fed', 'campsite', '旭川市21世紀の森ふれあい広場', '北海道旭川市東旭川町瑞穂4288', 43.7172721, 142.6674615, '2022年5月1日(日)〜11月30日(水)', '0166-76-2108', '有料。ログハウス大人290円〜750円、高校生以下180〜460円', '旭川市21世紀の森ふれあい広場は、ペーパンダムの周辺に整備された多目的公園、旭川市21世紀の森に隣接するキャンプ場です。', '/static/img/campsiteflag.jpeg'),
('5c5323e9-c78f-4dac-94ef-d34ab5ea8abc', 'campsite', '千代田の丘キャンプ場', '北海道上川郡美瑛町字水沢春日台第一', 43.5436008, 142.4912747, '毎年4月下旬～10月上旬', '0166-92-1718', '有料。', '千代田の丘キャンプ場は、美瑛町のファームズ千代田内に位置し、自然豊かな牧場空間を楽しめます。', '/static/img/campsiteflag.jpeg'),
('5c5323e9-c78f-4dac-94ef-d34ab5ea8def', 'campsite', 'とままえ夕陽ヶ丘未来港公園', '北海道苫前郡苫前町字栄浜313', 44.3153234, 141.6563455, '管理棟は7月中旬～8月中旬', '0164-64-2212', '不明', 'とままえ夕陽ヶ丘公園は、日本海に面した位置にある開放感あふれる公園です。', '/static/img/campsiteflag.jpeg');
