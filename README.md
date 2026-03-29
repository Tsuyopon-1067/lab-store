# 購買会計システム

研究室等の小規模コミュニティ向け購買管理・会計システム。

## ドキュメント

- [要件定義書](docs/requirements.md)
- [ディレクトリ構成・DB設計](docs/directory_and_db_design.md)
- [コンテナ設計](docs/container_design.md)

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
