# コンテナ設計書

**文書バージョン**: 1.1  
**作成日**: 2026-03-28

---

## 1. コンテナ構成概要

バックエンドとフロントエンドを別コンテナで構成し、podman-compose で管理する。  
SQLiteのデータファイルはホスト側のディレクトリをボリュームマウントして永続化する。  
開発時はコンテナを使わず、ホストで直接起動する（後述）。

```
ホストPC（Linux）
┌─────────────────────────────────────────────────────┐
│                                                     │
│  ブラウザ → localhost:3000                          │
│                  │                                  │
│  ┌───────────────▼──────────────────────────────┐  │
│  │  podman-compose                              │  │
│  │                                              │  │
│  │  ┌─────────────────┐  ┌──────────────────┐  │  │
│  │  │  frontend        │  │  backend         │  │  │
│  │  │  nginx + static  │  │  Go + Gin        │  │  │
│  │  │  :3000           │→ │  :8080           │  │  │
│  │  └─────────────────┘  └────────┬─────────┘  │  │
│  └────────────────────────────────│─────────────┘  │
│                                   │ volume mount    │
│                          ┌────────▼────────┐        │
│                          │  ./data/         │        │
│                          │  purchase.db     │        │
│                          └─────────────────┘        │
└─────────────────────────────────────────────────────┘
```

### コンテナ一覧

| コンテナ名 | 役割 | 公開ポート |
|-----------|------|-----------|
| purchase-backend | Go API サーバ | 8080（ホスト非公開・コンテナ間のみ） |
| purchase-frontend | SvelteKit静的配信 + /api プロキシ | 3000 → ホスト3000 |

---

## 2. compose.yaml

```yaml
name: purchase-system

services:
  backend:
    build:
      context: ./backend
      dockerfile: Containerfile
    container_name: purchase-backend
    restart: unless-stopped
    volumes:
      - ./data:/app/data
      - ./backup:/app/backup
      - ./config.yaml:/app/config.yaml:ro
    environment:
      - TZ=Asia/Tokyo
    expose:
      - "8080"

  frontend:
    build:
      context: ./frontend
      dockerfile: Containerfile
    container_name: purchase-frontend
    restart: unless-stopped
    ports:
      - "3000:3000"
    depends_on:
      - backend
    environment:
      - TZ=Asia/Tokyo
```

---

## 3. backend/Containerfile

マルチステージビルド。`modernc.org/sqlite` は純Go実装のため CGO_ENABLED=0 でビルド可能。

```dockerfile
# ---- ビルドステージ ----
FROM golang:1.23-bookworm AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o server ./main.go

# ---- 実行ステージ ----
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y \
    tzdata \
    ca-certificates \
    sqlite3 \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /build/server ./server

RUN mkdir -p /app/data /app/backup

EXPOSE 8080

CMD ["/app/server"]
```

---

## 4. frontend/Containerfile

bun でビルドし、nginx で静的配信する。

```dockerfile
# ---- ビルドステージ ----
FROM oven/bun:1 AS builder

WORKDIR /build

COPY package.json bun.lockb ./
RUN bun install --frozen-lockfile

COPY . .
RUN bun run build

# ---- 実行ステージ ----
FROM nginx:1.27-alpine

COPY --from=builder /build/build /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 3000

CMD ["nginx", "-g", "daemon off;"]
```

---

## 5. frontend/nginx.conf

```nginx
server {
    listen 3000;

    root /usr/share/nginx/html;
    index index.html;

    # SPA のルーティング
    location / {
        try_files $uri $uri/ /index.html;
    }

    # /api/* をバックエンドコンテナに転送
    location /api/ {
        proxy_pass         http://backend:8080;
        proxy_http_version 1.1;
        proxy_set_header   Host $host;
        proxy_set_header   X-Real-IP $remote_addr;
        proxy_set_header   X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_read_timeout 30s;
    }
}
```

---

## 6. SvelteKit の設定ファイル

### svelte.config.js

```javascript
import adapter from '@sveltejs/adapter-static';

export default {
  kit: {
    adapter: adapter({
      pages: 'build',
      assets: 'build',
      fallback: 'index.html',  // SPA モード
    }),
  },
};
```

### vite.config.js

```javascript
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
  plugins: [sveltekit()],
  server: {
    // 開発時: /api/* をホストで動いているバックエンドに転送
    proxy: {
      '/api': 'http://localhost:8080',
    },
  },
});
```

APIリクエストは常に相対パス `/api/...` で書く。  
開発時は vite の proxy が、本番コンテナでは nginx が転送するため、環境ごとの切り替えが不要。

---

## 7. systemd サービスファイル（deploy/purchase-system.service）

PC起動時にコンテナを自動起動する。

```ini
[Unit]
Description=Purchase System (podman-compose)
After=network.target

[Service]
Type=oneshot
RemainAfterExit=yes
WorkingDirectory=/opt/purchase-system
ExecStart=/usr/bin/podman-compose up -d
ExecStop=/usr/bin/podman-compose down
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

登録コマンド：

```bash
sudo cp deploy/purchase-system.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable purchase-system
sudo systemctl start purchase-system
```

---

## 8. 開発時の構成（コンテナ不使用）

開発中はホストで直接起動する。コンテナのビルド待ちが不要で効率的。

```bash
# バックエンド
cd backend
go run ./main.go          # localhost:8080 で起動
                          # gooseマイグレーションも自動実行

# フロントエンド（別ターミナル）
cd frontend
bun run dev               # localhost:5173 で起動・HMR有効
                          # /api/* は vite.config.js の proxy で localhost:8080 に転送
```

### 開発時のポート

| サービス | URL |
|----------|-----|
| フロントエンド（HMR） | http://localhost:5173 |
| バックエンドAPI | http://localhost:8080 |

---

## 9. セットアップ手順（初回）

```bash
# 1. リポジトリをクローン
git clone https://github.com/your-org/purchase-system.git
cd purchase-system

# 2. 設定ファイルを作成
cp deploy/config.yaml.template config.yaml

# 3. データディレクトリを作成
mkdir -p data backup

# 4. 【開発時】バックエンドの依存パッケージを取得
cd backend && go mod tidy && cd ..

# 5. 【開発時】フロントエンドの依存パッケージをインストール
cd frontend && bun install && cd ..

# 6. 管理者パスワードを設定（初回のみ）
cd backend && go run ./cmd/setup/main.go && cd ..

# 7. 【本番】コンテナをビルドして起動
podman-compose build
podman-compose up -d
```

---

## 10. よく使うコマンド

```bash
# コンテナ起動・停止
podman-compose up -d
podman-compose down

# ログ確認
podman-compose logs -f backend
podman-compose logs -f frontend

# コード変更後の再ビルド
podman-compose build backend   # バックエンドのみ
podman-compose build           # 全コンテナ
podman-compose up -d

# 手動バックアップ
podman-compose exec backend sqlite3 /app/data/purchase.db \
  ".backup /app/backup/manual_$(date +%Y%m%d_%H%M%S).sqlite"
```

---

## 11. データの永続化

コンテナを削除・再ビルドしてもデータは失われない。

| ホスト側パス | コンテナ側パス | 内容 |
|-------------|--------------|------|
| `./data/` | `/app/data/` | SQLiteデータベース |
| `./backup/` | `/app/backup/` | バックアップファイル |
| `./config.yaml` | `/app/config.yaml` | 設定ファイル（読み取り専用） |
