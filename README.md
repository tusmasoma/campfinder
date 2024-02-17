# campfinder
campfinder

## アーキテクチャ
レイヤードアーキテクチャ + DDD を採用

採用理由は、アプリケーションの構造を明確に分離して依存関係を管理しやすくすることで、複数人の開発でも品質を維持しやすくするためです。また、ビジネスの複雑さを効率的に扱えるドメイン中心の設計により、小規模から中規模のプロジェクトに適していたからです。

[クリーンアーキテクチャを採用したcampfinderのコード](https://github.com/tusmasoma/clean-architecture-campfinder/tree/main)

## インフラ構成
![campfinder_ver_1_22 drawio-2](https://github.com/tusmasoma/campfinder/assets/104899572/073b3d49-8c7c-4b9f-9227-e4a6a99dee39)

## curl集
現在は、コストの観点からデプロイを停止しています。

#### User
```
curl -v -X 'POST' \
  'https://production.campfinderjp.com/api/user/create' \
  -H 'accept: */*' \
  -H 'Content-Type: application/json' \
  -d '{
  "email": "example@gmail.com",
  "password": "example12345"
}'
```

```
curl -v -X 'POST' \
  'https://production.campfinderjp.com/api/user/login' \
  -H 'accept: */*' \
  -H 'Content-Type: application/json' \
  -d '{
  "email": "example@gmail.com",
  "password": "example12345"
}'
```

```
curl -v -X GET 'https://production.campfinderjp.com/api/user/logout'
```

#### Spot
```
curl -v -X GET 'https://production.campfinderjp.com/api/spot?category=campsite'         
```

```
curl -X 'POST' \
  'https://production.campfinderjp.com/api/spot/create' \
  -H 'accept: */*' \
  -H 'Content-Type: application/json' \
  -d '{
  "category": "campsite",
  "name": "旭川市21世紀の森ふれあい広場",
  "address": "北海道旭川市東旭川町瑞穂888",
  "lat": 43.7172721,
  "lng": 142.6674615,
  "period": "2022年5月1日(日)～11月30日(水)",
  "phone": "0166-76-2108",
  "price": "有料",
  "description": "旭川市21世紀の森ふれあい広場は、ペーパンダムの周辺に整備された多目的公園、旭川市21世紀の森に隣接するキャンプ場です。",
  "iconpath": "/static/img/campsiteflag.jpeg"
}'
```

#### Image
```
curl -X GET 'https://production.campfinderjp.com/api/img?spot_id='
```

```
curl -v -X POST \  
  'https://production.campfinderjp.com/api/img/create' \
  -H 'accept: */*' \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer ' \
  -d '{
    "spotID": "",
    "url": "https://lh3.googleusercontent.com/places/AJQcZqLYx0Skfw6yLIKZEwYt4sugN-O3dQ7RTra-jVe6lnTDXj1iW5IPuBXLspKbvoRI8pb5PGvkee4nZMtcsveQpQ4QS3TuruFOpL4=s1600-w400"
}'
```

