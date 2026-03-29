# CLAUDE.md

## プロジェクト概要

研究室等の小規模コミュニティ向け購買会計システム。  
商品のバーコードスキャンによる購入記録、仕入れ立替管理、月次精算を行う。  
1台のLinux PC上でPodmanコンテナとして動作し、ブラウザからアクセスする。

## ドキュメント

実装前に必ず以下をすべて読むこと。

| ファイル | 内容 |
|----------|------|
| `docs/requirements.md` | 要件定義書（機能・API・画面一覧・お金の流れ） |
| `docs/directory_and_db_design.md` | ディレクトリ構成・DBマイグレーション・主要クエリ・Ginルーター構成 |
| `docs/container_design.md` | コンテナ構成・Containerfile・compose.yaml・nginx設定 |

## 技術スタック

| レイヤ | 技術 |
|--------|------|
| バックエンド | Go + Gin |
| データベース | SQLite（WALモード） |
| SQLiteドライバ | modernc.org/sqlite（CGO不要） |
| DBマイグレーション | goose（アプリ起動時に自動実行） |
| フロントエンド | SvelteKit |
| パッケージマネージャ | bun |
| バーコードスキャン | ZXing-js |
| コンテナ | Podman + podman-compose |

## リポジトリ構成

```
purchase-system/
├── CLAUDE.md
├── compose.yaml
├── config.yaml              # Git管理外（config.yaml.templateからコピー）
├── docs/
│   ├── requirements.md
│   ├── directory_and_db_design.md
│   └── container_design.md
├── backend/
│   ├── Containerfile
│   ├── main.go
│   ├── go.mod
│   └── ...
├── frontend/
│   ├── Containerfile
│   ├── nginx.conf
│   ├── package.json
│   ├── bun.lockb
│   ├── svelte.config.js
│   ├── vite.config.js
│   └── src/
├── deploy/
│   ├── purchase-system.service
│   └── config.yaml.template
├── data/                    # Git管理外
└── backup/                  # Git管理外
```

## 重要な設計方針

- 残高はテーブルに持たず、毎回集計クエリで算出する（整合性を保つため）
- 支払い・精算・仕入れ記録の修正・削除は論理削除（deleted_at）で行い、audit_logsに変更前後をJSONで記録する
- 商品の価格変更は product_prices テーブルに履歴として追記し、商品マスターは変更しない
- バーコードと商品は常に1対1で対応する
- gooseのマイグレーションはアプリ起動時（db.Init）に自動実行される
- 管理者認証はセッションCookie＋bcrypt
- 一般利用者の本人確認はX-User-Barcodeヘッダで行う
- APIリクエストはフロントエンドから相対パス /api/... で統一する（nginxが転送）

## ポート

| サービス | ポート |
|----------|--------|
| フロントエンド（nginx） | 3000（ホスト公開） |
| バックエンド（Gin） | 8080（コンテナ内部のみ） |

## 開発時のコマンド

```bash
# バックエンド（ホストで直接起動）
cd backend
go run ./main.go

# フロントエンド（ホストで直接起動・HMR有効）
cd frontend
bun run dev

# コンテナビルド・起動（本番確認時）
podman-compose build
podman-compose up -d

# 管理者パスワード初期設定（初回のみ）
## 開発時
cd backend && go run ./cmd/setup/main.go
## コンテナ起動後
podman-compose exec backend /app/server setup

# マイグレーション状態確認
goose sqlite ./data/purchase.db status
```

## go.mod の主要依存パッケージ

```
github.com/gin-gonic/gin
github.com/pressly/goose/v3
modernc.org/sqlite
golang.org/x/crypto/bcrypt
gopkg.in/yaml.v3
```

## フロントエンドの主要パッケージ

```
@sveltejs/kit
@sveltejs/adapter-static
@zxing/browser
@zxing/library
```
