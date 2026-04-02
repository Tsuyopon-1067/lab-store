# 購買会計システム

研究室等の小規模コミュニティ向け購買管理・会計システム。

## ドキュメント

### 開発・設計
- [要件定義書](docs/requirements.md)
- [ディレクトリ構成・DB設計](docs/directory_and_db_design.md)
- [コンテナ設計](docs/container_design.md)

### 本番環境
- [本番環境デプロイメントガイド](docs/production_deployment.md)
- [デプロイチェックリスト](docs/production_checklist.md)

## 技術スタック

- **バックエンド**: Go + Gin + SQLite（goose でマイグレーション）
- **フロントエンド**: SvelteKit + bun + ZXing-js
- **コンテナ**: Podman + podman-compose

## セットアップ

### 前提条件

- Go 1.23+
- bun
- Podman + podman-compose（本番環境のみ）

### 初回セットアップ

```bash
# 1. 設定ファイルを作成
cp deploy/config.yaml.template config.yaml

# 2. データディレクトリを作成
mkdir -p data backup

# 3. バックエンドの依存パッケージを取得
cd backend && go mod tidy && cd ..

# 4. フロントエンドの依存パッケージをインストール
cd frontend && bun install && cd ..

# 5. 管理者パスワードを設定（初回のみ）
cd backend && go run ./cmd/setup/main.go && cd ..
```

### 開発時の起動

```bash
# バックエンド（localhost:8080）
cd backend && go run ./main.go

# フロントエンド（localhost:5173、別ターミナル）
cd frontend && bun run dev
```

### 本番環境（コンテナ）

```bash
# コンテナをビルドして起動
podman-compose build
podman-compose up -d

# systemd に登録して自動起動
sudo cp deploy/purchase-system.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable --now purchase-system
```

## ブラウザアクセス

- 開発時: http://localhost:5173
- 本番時: http://localhost:3000

## 外部API（APIキー認証）

### 概要

ボット・外部システム向けの API キー認証を使用した外部 API エンドポイントを提供しています。

### APIキーの生成

1. 管理画面にログイン
2. 管理者メニュー → API キー管理ページへアクセス
3. **新規作成** ボタンから、APIキーに付ける名前を入力
4. 表示された生キー（64文字）を**安全に保管** - このタイミングでのみ表示されます

### 使用方法

外部からのリクエストに `X-API-Key` ヘッダを含めて API を呼び出します。

```bash
curl -H "X-API-Key: YOUR_API_KEY" \
  http://localhost:8080/v1/users/4912345678901/balance
```

### 利用可能なエンドポイント

#### 1. ユーザーの残高を取得（バーコード指定）

```
GET /v1/users/{barcode}/balance
```

**レスポンス例:**
```json
{
  "user_id": 1,
  "user_name": "山田 太郎",
  "purchase_unpaid": 1000,
  "restock_unclaimed": 800,
  "net_balance": 200
}
```

#### 2. ユーザーの残高を取得（名前指定）

```
GET /v1/users/by-name/{name}/balance
```

**レスポンス例:**
```json
{
  "user_id": 1,
  "user_name": "山田 太郎",
  "purchase_unpaid": 1000,
  "restock_unclaimed": 800,
  "net_balance": 200
}
```

#### 3. 月次レポートを取得

```
GET /v1/reports/monthly?from=2026-03-01&to=2026-03-31
```

**クエリパラメータ:**
- `from`: 集計開始日（YYYY-MM-DD形式）
- `to`: 集計終了日（YYYY-MM-DD形式）

**レスポンス例:**
```json
{
  "period": {
    "from": "2026-03-01",
    "to": "2026-03-31"
  },
  "users": [
    {
      "user_id": 1,
      "user_name": "山田 太郎",
      "purchase_total": 1000,
      "purchase_paid": 800,
      "purchase_unpaid": 200,
      "restock_total": 500,
      "restock_settled": 500,
      "restock_unclaimed": 0,
      "net_balance": 200
    }
  ]
}
```

#### 4. 商品一覧を取得（外部向け）

```
GET /v1/products
```

**レスポンス例:**
```json
[
  {
    "id": 1,
    "name": "コーヒー",
    "barcode": "4912345678901",
    "current_price": 100,
    "is_active": 1
  }
]
```

#### 5. 全利用者の支払い金額サマリーを取得

```
GET /v1/users/payment-summary
```

全有効利用者の支払い金額サマリーを取得します（期間フィルターなし、全期間の合計）。

**レスポンス例:**
```json
[
  {
    "user_id": 1,
    "user_name": "山田 太郎",
    "purchase_unpaid": 1000,
    "restock_unclaimed": 500,
    "net_balance": 500
  },
  {
    "user_id": 2,
    "user_name": "鈴木 花子",
    "purchase_unpaid": 2000,
    "restock_unclaimed": 0,
    "net_balance": 2000
  }
]
```

**フィールド説明:**
- `user_id`: 利用者ID
- `user_name`: 利用者名
- `purchase_unpaid`: 未払い利用額（購入合計 - 支払い済み合計）
- `restock_unclaimed`: 未精算仕入れ立替額（仕入れ立替合計 - 精算済み合計）
- `net_balance`: 最終的な支払い金額（purchase_unpaid - restock_unclaimed）
  - 正の値：利用者がコミュニティに支払う必要がある金額
  - 負の値：コミュニティが利用者に支払う必要がある金額

#### 6. 支払い金額サマリーをCSV出力

```
GET /v1/users/payment-summary/export
```

全利用者の支払い金額サマリーをCSVファイル形式でダウンロードします。

**CSVヘッダー:**
```
ユーザー名,利用額,仕入れ金額,最終的な支払い金額
```

**使用例:**
```bash
curl -H "X-API-Key: YOUR_API_KEY" \
  http://localhost:8080/v1/users/payment-summary/export \
  -o user_payment_summary.csv
```

### セキュリティ

- API キーは **SHA256 ハッシュ** で保存されます
- 生キーは**作成時のみ**表示されます
- キープレフィックス（先頭8文字）で管理可能
- `is_active` フラグで有効/無効を切り替え可能
- `last_used_at` で最後の使用日時を追跡可能
- 不要なキーは削除してください

### API キーの管理

管理画面から以下の操作が可能です：

- **一覧表示**: 作成済みのすべての API キーを確認
- **削除**: 不要なキーを削除（削除後は復元不可）
- **有効/無効**: キーを有効/無効に切り替え
