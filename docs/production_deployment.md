# 本番環境デプロイメントガイド

このドキュメントはLinux サーバー上への本番環境デプロイを説明します。

---

## 前提条件

- Linux サーバー（Ubuntu 22.04 LTS 以上推奨）
- Podman 4.0+ または Docker 20.10+
- podman-compose または docker-compose
- 管理者権限（root または sudo）

---

## 1. リポジトリのセットアップ

```bash
# リポジトリをクローン（または転送）
git clone <repository-url> /opt/purchase-system
cd /opt/purchase-system

# 必要なディレクトリを作成
mkdir -p data backup

# 設定ファイルを作成
cp deploy/config.yaml.template config.yaml

# config.yaml を本番用に編集
# - server.port: 8080（変更不可、コンテナ内部）
# - database.path: /app/data/purchase.db（自動）
# - session.timeout_minutes: 本番環境に合わせて調整
# - backup.path: /app/backup（自動）
nano config.yaml
```

---

## 2. 初期管理者パスワード設定

**本番環境にデプロイする前に、必ず管理者パスワードを設定してください。**

開発環境で設定済みの場合は、以下をスキップして既存の `data/purchase.db` を本番環境にコピーしてください。

```bash
# 開発環境で（まだの場合）
cd backend
go run ./cmd/setup/main.go
cd ..

# または本番環境のコンテナで実行
podman-compose build
podman-compose run --rm backend /app/server setup
```

---

## 3. コンテナのビルド・起動

### Podman を使う場合（推奨・CentOS/RHEL/Fedora）

```bash
# インストール
sudo dnf install -y podman podman-compose

# コンテナをビルド
podman-compose build

# バックグラウンドで起動
podman-compose up -d

# ログ確認
podman-compose logs -f
```

### Docker を使う場合（Ubuntu/Debian）

```bash
# インストール
sudo apt-get install -y docker.io docker-compose

# dockerグループにユーザーを追加（optionalだが推奨）
sudo usermod -aG docker $USER

# コンテナをビルド
docker-compose build

# バックグラウンドで起動
docker-compose up -d

# ログ確認
docker-compose logs -f
```

---

## 4. 起動確認

```bash
# コンテナの状態確認
podman-compose ps
# または
docker-compose ps

# ポートが開いているか確認
netstat -tlnp | grep 3000
# または
ss -tlnp | grep 3000

# ヘルスチェック
curl http://localhost:3000/
```

ブラウザで `http://localhost:3000` を開いて、フロントエンドが表示されることを確認してください。

---

## 5. systemd サービス登録（自動起動設定）

PC起動時に自動的にコンテナを起動する設定。

```bash
# サービスファイルをコピー
sudo cp deploy/purchase-system.service /etc/systemd/system/

# systemd をリロード
sudo systemctl daemon-reload

# サービスを有効化・起動
sudo systemctl enable purchase-system
sudo systemctl start purchase-system

# ステータス確認
sudo systemctl status purchase-system

# ログ確認
sudo journalctl -u purchase-system -f
```

---

## 6. ファイアウォール設定

ブラウザからアクセスするには、ファイアウォールでポート 3000 と 443 (HTTPS) を開く必要があります。

### UFW（Ubuntu/Debian）

```bash
sudo ufw allow 3000/tcp
sudo ufw allow 443/tcp
sudo ufw reload
```

### Firewalld（CentOS/RHEL）

```bash
sudo firewall-cmd --permanent --add-port=3000/tcp
sudo firewall-cmd --permanent --add-port=443/tcp
sudo firewall-cmd --reload
```

---

## 7. SSL/TLS 設定（本番必須）

### 概要

本システムのコンテナ内の nginx は HTTPS で通信します。Let's Encrypt の無料証明書を使用します。

**構成:**
```
ブラウザ → https://example.com:443 → nginx (HTTPS)
                                     ↓
                                   backend:8080 (HTTP内部)
```

### 7.1 certbot のインストール

```bash
# Ubuntu/Debian
sudo apt-get install -y certbot

# CentOS/RHEL
sudo dnf install -y certbot
```

### 7.2 deploy.sh でワンコマンドセットアップ

```bash
# SSL/TLSをセットアップ（推奨）
./deploy.sh setup-ssl
```

以下を入力します：
- ドメイン名（例: `api.example.com`）
- Let's Encrypt登録用メールアドレス

スクリプトが自動で以下を実行します：
1. ✅ certbot で証明書を取得
2. ✅ /etc/letsencrypt にマウント
3. ✅ nginx が参照できるシンボリックリンクを作成

### 7.3 config.yaml を更新

セットアップ後、`config.yaml` を編集して HTTPS を有効化：

```yaml
https:
  enabled: true
  domain: "api.example.com"
  email: "admin@example.com"
```

### 7.4 コンテナを再起動

```bash
./deploy.sh restart
```

### 7.5 HTTPS でアクセス確認

```bash
curl https://api.example.com/
# または
curl --insecure https://localhost:443/
```

ブラウザで `https://api.example.com` を開いて、🔒 マークが表示されることを確認してください。

### 7.6 証明書の自動更新設定

Let's Encrypt 証明書は **90日で期限切れ**になるため、自動更新が必須です。

```bash
# 更新テストを実施
./deploy.sh renew-ssl

# systemd タイマーで自動更新を有効化
sudo systemctl enable certbot.timer
sudo systemctl start certbot.timer

# 自動更新の状態確認
sudo systemctl status certbot.timer
sudo systemctl list-timers --all | grep certbot
```

certbot は **毎日** 証明書を確認し、30日以内に期限切れの場合は自動的に更新します。

### 7.7 トラブルシューティング

**証明書取得に失敗した場合:**

```bash
# ファイアウォールでポート 80 が開いているか確認
sudo ufw allow 80/tcp
sudo ufw reload

# または certbot を手動実行
sudo certbot certonly --standalone -d api.example.com
```

**nginx で証明書が見つからないエラー:**

```bash
# パーミッション確認
sudo ls -la /etc/letsencrypt/live/

# nginx コンテナを再起動
./deploy.sh restart frontend
```

**期限切れ近い証明書:**

```bash
# 証明書一覧と有効期限を確認
sudo certbot certificates

# 手動更新
sudo certbot renew --force-renewal
```

---

## 8. バックアップ戦略

### 自動バックアップ

コンテナ内部で `db/backups.go` で実装されているバックアップ機能が動作しています。

### 手動バックアップ

```bash
# バックアップディレクトリを確認
ls -la ./backup/

# 手動でバックアップを取得
podman-compose exec backend sqlite3 /app/data/purchase.db \
  ".backup /app/backup/manual_$(date +%Y%m%d_%H%M%S).sqlite"

# バックアップを外部に転送
scp -r ./backup/ backup-server:/backups/purchase-system/
```

### データベースのリストア

```bash
# バックアップファイルをコンテナ環境にコピー
cp /path/to/backup.sqlite ./data/purchase.db

# コンテナを再起動
podman-compose restart backend
```

---

## 9. よく使うコマンド

| コマンド | 説明 |
|---------|------|
| `./deploy.sh up` | コンテナビルド・起動 |
| `./deploy.sh setup-ssl` | SSL/TLS証明書セットアップ |
| `./deploy.sh logs` | ログ表示 |
| `./deploy.sh restart` | コンテナ再起動 |
| `./deploy.sh down` | コンテナ停止 |

### 詳細なコマンド例

```bash
# ログ確認（全コンテナ）
podman-compose logs -f

# 特定コンテナのログ
podman-compose logs -f backend
podman-compose logs -f frontend

# コンテナの停止
podman-compose stop

# コンテナの再起動
podman-compose restart

# コンテナの削除（データは保持）
podman-compose down

# コンテナのアップグレード（コード変更後）
podman-compose build backend
podman-compose up -d

# 全コンテナの再構築
podman-compose build
podman-compose up -d
```

---

## 10. トラブルシューティング

### ポート 3000 がすでに使用されている

```bash
# 既存のプロセスを確認・停止
lsof -i :3000
kill -9 <PID>

# または compose.yaml のポート定義を変更
# ports:
#   - "8000:3000"  # ホスト 8000 → コンテナ 3000
```

### データベースロック

```bash
# WAL ファイルを削除（最後の手段）
rm ./data/purchase.db-shm
rm ./data/purchase.db-wal

# コンテナを再起動
podman-compose restart backend
```

### 管理者パスワード忘却

```bash
# 初期化スクリプトを再実行
podman-compose run --rm backend /app/server setup
```

---

## 11. 監視・ヘルスチェック

### ヘルスチェックエンドポイント

```bash
curl http://localhost:3000/health
# または
curl http://localhost:8080/health
```

### systemd ログ監視

```bash
sudo journalctl -u purchase-system -f -n 50
```

### Prometheus メトリクス（将来実装予定）

```
GET http://localhost:8080/metrics
```

---

## 12. アップグレード手順

新しいバージョンへのアップグレード：

```bash
# コードを最新にプル
git pull origin main

# 新しい設定があれば確認
diff config.yaml deploy/config.yaml.template

# コンテナを再ビルド・再起動
podman-compose build
podman-compose up -d

# ログで エラーがないか確認
podman-compose logs -f
```

---

## 付録：環境別推奨値

### 開発環境
- session.timeout_minutes: 120（デバッグに十分な時間）
- backup.interval_minutes: 0（自動バックアップ無効）
- https.enabled: false

### 本番環境
- session.timeout_minutes: 30（セキュリティ）
- backup.interval_minutes: 60（1時間ごとにバックアップ）
- https.enabled: true

### 自動バックアップ

コンテナ内部で `db/backups.go` で実装されているバックアップ機能が動作しています。

### 手動バックアップ

```bash
# バックアップディレクトリを確認
ls -la ./backup/

# 手動でバックアップを取得
podman-compose exec backend sqlite3 /app/data/purchase.db \
  ".backup /app/backup/manual_$(date +%Y%m%d_%H%M%S).sqlite"

# バックアップを外部に転送
scp -r ./backup/ backup-server:/backups/purchase-system/
```

### データベースのリストア

```bash
# バックアップファイルをコンテナ環境にコピー
cp /path/to/backup.sqlite ./data/purchase.db

# コンテナを再起動
podman-compose restart backend
```

