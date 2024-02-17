# campfinder

## アーキテクチャ
レイヤードアーキテクチャ + DDD を採用

採用理由は、アプリケーションの構造を明確に分離して依存関係を管理しやすくすることで、複数人の開発でも品質を維持しやすくするためです。また、ビジネスの複雑さを効率的に扱えるドメイン中心の設計により、小規模から中規模のプロジェクトに適していたからです。

[クリーンアーキテクチャを採用したcampfinderのコード](https://github.com/tusmasoma/clean-architecture-campfinder/tree/main)

## インフラ構成
![campfinder_ver_1_22 drawio-2](https://github.com/tusmasoma/campfinder/assets/104899572/073b3d49-8c7c-4b9f-9227-e4a6a99dee39)

## infra層について
todo: genericsで汎用化

## テスト戦略
### 概要
本プロジェクトでは、品質を保証するためにユニットテストを重視しています。テストの自動化により、開発の効率化とバグの早期発見を目指しています。

### テストレベル
- **ユニットテスト**: 各関数やメソッドの正確性を検証します。すべてのコード変更に対して自動的に実行されます。

### テスト技法
- **モックテスト**: [gomock](https://github.com/golang/mock) ライブラリを使用して、外部依存関係をモック化し、テストの隔離性を向上させます。
- **データベーステスト**: [dockertest](https://github.com/ory/dockertest) ライブラリを利用して、テスト用のDBインスタンスをDockerコンテナで立ち上げ、データベース操作のテストを実行します。

### テストツール
- **ユニットテスト**: [go test](https://golang.org/pkg/testing/)
- **モック生成**: [gomock](https://github.com/golang/mock)
- **データベーステスト**: [dockertest](https://github.com/ory/dockertest)



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

## 今後の開発計画

### E2Eテストの導入
現在、本プロジェクトではユニットテストを中心にテスト戦略を展開していますが、今後の開発においてはエンドツーエンド(E2E)テストの導入を検討しています。E2Eテストにより、ユーザーの視点からアプリケーション全体の動作を検証し、より包括的な品質保証を目指します。

E2Eテストの導入には、テスト自動化ツールやテストシナリオの設計など、様々な準備が必要となるため、段階的にアプローチを進めていく予定です。E2Eテストの導入により、リリース前の最終的な確認作業を効率化し、ユーザーに提供する製品の品質をさらに向上させることを目指しています。

### 認証機能の切り出しとサードパーティ認証の導入
現在、本プロジェクトでは独自の認証システムを使用していますが、今後の開発では認証機能の切り出しとサードパーティ認証の導入を検討しています。このために、Firebase Authentication、Amazon Cognito、Auth0などの認証サービスを利用する可能性があります。これにより、Google、Facebook、Twitterなどの外部サービスを利用した認証が可能となり、ユーザーの利便性が向上します。

認証機能の切り出しにより、認証処理を一元管理しやすくなるため、セキュリティの向上やメンテナンスの効率化が期待できます。また、サードパーティ認証の導入により、ユーザーは複数のアカウントを持つことなく、既存のアカウントでサービスを利用できるようになります。

この取り組みは段階的に進めていく予定であり、まずは認証機能の切り出しを行い、その後サードパーティ認証サービスの導入を検討していきます。
