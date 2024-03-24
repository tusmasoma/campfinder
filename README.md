# campfinder

## アーキテクチャ
レイヤードアーキテクチャ + DDD を採用

採用理由は、アプリケーションの構造を明確に分離して依存関係を管理しやすくすることで、複数人の開発でも品質を維持しやすくするためです。また、ビジネスの複雑さを効率的に扱えるドメイン中心の設計により、小規模から中規模のプロジェクトに適していたからです。

[クリーンアーキテクチャを採用したcampfinderのコード](https://github.com/tusmasoma/clean-architecture-campfinder/tree/main)

## インフラ構成
![campfinder_ver_1_22 drawio-2](https://github.com/tusmasoma/campfinder/assets/104899572/073b3d49-8c7c-4b9f-9227-e4a6a99dee39)

## Development
### Format
goimports, gofmt
```makefile
make fmt
```
### Generate
go generate
```makefile
make generate
```
### Build Go
```makefile
make build
```
### Lint
```makefile
make lint
```
これは ./docker/back/... を対象として golangci-lint を実行するため、実用的ではありません。 実際には、次のように PKG 変数を指定し、package を限定した状態で実行することをお勧めします。
```makefile
make lint PKG=./docker/back/infra/...
```
後述の make lint-diff では差分のみを対象とするため、既存のコードには利用できません。 既存のコードの品質改善を行いたい場合には make lint が有用です。

### Lint (diff)
```makefile
make lint-diff
```
PKG 指定もできます。
```
make lint-diff PKG="./docker/back/infra/..."
```
make lint だと既存の指摘が多く、追加したコードに対する解析結果が判別しにくいため、 develop branch との差分の解析結果を表示する lint-diff を用意しています。 利用するには事前に reviewdog のインストールが必要です。
```shell
curl -sfL https://raw.githubusercontent.com/reviewdog/reviewdog/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
```
### Test
```makefile
make test
```
## infra層について
todo: genericsで汎用化

## タスク管理

本プロジェクトでは、GitHubのIssueを活用してタスク管理を行うことにします。開発に関連するすべてのタスクは、Issueとして登録され、進捗管理や議論が行われます。(2024/2/18)

### Issueの活用方法
- **新しいタスクの登録**: 新しい機能の追加やバグの修正など、開発に関連する新しいタスクは、GitHubのIssueとして登録してください。Issueには、タスクの背景、目的、具体的なタスクリスト、受け入れ基準などを明記してください。
- **タスクの進捗管理**: Issueにはラベルやマイルストーンを付与して、タスクの進捗状況を管理します。また、プルリクエストとIssueをリンクさせることで、コードの変更とタスクの進捗を関連付けます。
- **議論とフィードバック**: Issueのコメント機能を使用して、タスクに関する議論やフィードバックを行います。チームメンバー間のコミュニケーションを促進し、より良い解決策を見つけるために活用してください。

参考: https://qiita.com/tkmd35/items/9612c03dc60b1c516969
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

### データベースマイグレーションの導入
現在、本プロジェクトではデータベーススキーマの変更を手動で管理していますが、今後の開発ではデータベースマイグレーションツールの導入を検討しています。これにより、スキーマの変更をより効率的かつ安全に行うことができるようになります。

データベースマイグレーションツールを導入することで、以下の利点が期待できます：
- **バージョン管理**: スキーマの変更履歴を追跡し、特定のバージョンへのロールバックが可能になります。
- **自動化**: マイグレーションの実行を自動化することで、手作業によるエラーを減らします。
- **環境の統一**: 開発、テスト、本番環境でのスキーマの差異を排除し、一貫性を保ちます。

検討しているマイグレーションツールには、Flyway、Liquibase、Goのマイグレーションライブラリなどがあります。ツールの選定は、プロジェクトの技術スタックやチームのニーズに基づいて行われます。

マイグレーションの導入プロセスは段階的に進めていく予定であり、まずは開発環境での試験的な導入を行い、その後テスト環境および本番環境へと展開していきます。
