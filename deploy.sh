#!/bin/bash

# 購買会計システム本番環境デプロイスクリプト
# 使用方法: ./deploy.sh [build|up|down|restart|logs]

set -e

COMPOSE_CMD=${1:-help}
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

cd "$SCRIPT_DIR"

# 色出力用
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 前提条件チェック
check_requirements() {
    log_info "前提条件をチェック中..."

    # docker-compose または podman-compose をチェック
    if ! command -v podman-compose &> /dev/null; then
        if ! command -v docker-compose &> /dev/null; then
            log_error "podman-compose または docker-compose がインストールされていません"
            exit 1
        fi
        COMPOSE="docker-compose"
    else
        COMPOSE="podman-compose"
    fi

    log_info "使用するコンテナツール: $COMPOSE"

    # 設定ファイルの確認
    if [ ! -f config.yaml ]; then
        log_error "config.yaml が見つかりません"
        log_info "以下のコマンドで作成してください:"
        echo "  cp deploy/config.yaml.template config.yaml"
        exit 1
    fi

    # データディレクトリの確認
    if [ ! -d data ]; then
        log_info "data/ ディレクトリを作成します..."
        mkdir -p data
    fi

    if [ ! -d backup ]; then
        log_info "backup/ ディレクトリを作成します..."
        mkdir -p backup
    fi

    log_info "前提条件チェック完了 ✓"
}

# コンテナのビルド
build_containers() {
    log_info "コンテナをビルド中..."
    $COMPOSE build
    log_info "ビルド完了 ✓"
}

# コンテナの起動
start_containers() {
    log_info "コンテナを起動中..."
    $COMPOSE up -d
    log_info "起動完了 ✓"

    # 起動確認
    sleep 2
    log_info "ステータス確認:"
    $COMPOSE ps
}

# コンテナの停止
stop_containers() {
    log_info "コンテナを停止中..."
    $COMPOSE down
    log_info "停止完了 ✓"
}

# コンテナの再起動
restart_containers() {
    log_info "コンテナを再起動中..."
    $COMPOSE restart
    log_info "再起動完了 ✓"
}

# ログ表示
show_logs() {
    log_info "ログを表示中 (Ctrl+C で終了)..."
    $COMPOSE logs -f
}

# 管理者パスワード初期化
setup_admin() {
    log_info "管理者パスワードを初期化中..."
    check_requirements
    $COMPOSE run --rm backend /app/server setup
    log_info "初期化完了 ✓"
}

# SSL/TLS 証明書セットアップ
setup_ssl() {
    log_info "SSL/TLS 証明書をセットアップ中..."

    # 前提条件チェック
    if ! command -v certbot &> /dev/null; then
        log_error "certbot がインストールされていません"
        log_info "以下のコマンドでインストールしてください:"
        echo "  sudo apt-get install -y certbot python3-certbot-nginx"
        exit 1
    fi

    # config.yaml から情報を取得
    if [ ! -f config.yaml ]; then
        log_error "config.yaml が見つかりません"
        exit 1
    fi

    # ドメイン名をユーザーに入力させる
    read -p "ドメイン名を入力してください (例: example.com): " DOMAIN
    read -p "Let's Encrypt登録用メールアドレスを入力してください: " EMAIL

    if [ -z "$DOMAIN" ] || [ -z "$EMAIL" ]; then
        log_error "ドメイン名とメールアドレスが必須です"
        exit 1
    fi

    # certbot ディレクトリを作成
    mkdir -p certbot

    log_info "Let's Encrypt から証明書を取得中 (domain: $DOMAIN)..."

    # webroot モードで証明書を取得
    sudo certbot certonly \
        --webroot \
        -w ./certbot \
        -d "$DOMAIN" \
        --email "$EMAIL" \
        --agree-tos \
        --non-interactive \
        --expand

    if [ $? -eq 0 ]; then
        log_info "証明書取得成功 ✓"

        # シンボリックリンクを作成（nginx が参照する場所）
        sudo mkdir -p /etc/letsencrypt/live/default
        sudo ln -sf "/etc/letsencrypt/live/$DOMAIN/fullchain.pem" /etc/letsencrypt/live/default/fullchain.pem
        sudo ln -sf "/etc/letsencrypt/live/$DOMAIN/privkey.pem" /etc/letsencrypt/live/default/privkey.pem

        # パーミッション設定
        sudo chmod -R 755 /etc/letsencrypt/live
        sudo chmod -R 755 /etc/letsencrypt/archive

        log_info "config.yaml を更新してください:"
        echo "  https.enabled: true"
        echo "  https.domain: $DOMAIN"
        echo "  https.email: $EMAIL"

        log_info "その後 './deploy.sh up' で再起動してください"
    else
        log_error "証明書取得に失敗しました"
        exit 1
    fi
}

# SSL/TLS 証明書更新テスト
renew_ssl() {
    log_info "SSL/TLS 証明書の更新テストを実行中..."

    if ! command -v certbot &> /dev/null; then
        log_error "certbot がインストールされていません"
        exit 1
    fi

    sudo certbot renew --webroot -w ./certbot --dry-run

    if [ $? -eq 0 ]; then
        log_info "更新テスト成功 ✓"
        log_info "以下のコマンドで定期更新を設定してください:"
        echo "  sudo systemctl enable certbot.timer"
        echo "  sudo systemctl start certbot.timer"
    else
        log_error "更新テストに失敗しました"
        exit 1
    fi
}

# ヘルスチェック
health_check() {
    log_info "ヘルスチェック実施中..."

    # コンテナの起動状態確認
    if ! $COMPOSE ps | grep -q "purchase-backend"; then
        log_error "バックエンドコンテナが起動していません"
        return 1
    fi

    if ! $COMPOSE ps | grep -q "purchase-frontend"; then
        log_error "フロントエンドコンテナが起動していません"
        return 1
    fi

    # ポート確認
    if ! curl -s http://localhost:3000/health &> /dev/null; then
        log_warn "http://localhost:3000/health へのアクセスに失敗しました"
    else
        log_info "http://localhost:3000/health ✓"
    fi

    log_info "ヘルスチェック完了 ✓"
}

# メインの処理
case $COMPOSE_CMD in
    build)
        check_requirements
        build_containers
        ;;
    up)
        check_requirements
        build_containers
        start_containers
        health_check
        ;;
    down)
        stop_containers
        ;;
    restart)
        restart_containers
        health_check
        ;;
    logs)
        show_logs
        ;;
    setup)
        setup_admin
        ;;
    setup-ssl)
        setup_ssl
        ;;
    renew-ssl)
        renew_ssl
        ;;
    health)
        health_check
        ;;
    ps)
        check_requirements
        $COMPOSE ps
        ;;
    help|"")
        cat << EOF
購買会計システム本番環境デプロイスクリプト

使用方法:
  ./deploy.sh [コマンド]

コマンド:
  up           コンテナをビルド・起動（初回推奨）
  build        コンテナのみビルド
  down         コンテナを停止・削除
  restart      コンテナを再起動
  logs         ログを表示（リアルタイム）
  setup        管理者パスワード初期化
  setup-ssl    SSL/TLS証明書をセットアップ（Let's Encrypt）
  renew-ssl    SSL/TLS証明書の更新テスト
  health       ヘルスチェック実施
  ps           コンテナ一覧表示
  help         このヘルプを表示

例:
  ./deploy.sh up       # 初回セットアップ・起動
  ./deploy.sh logs     # ログ表示
  ./deploy.sh down     # 停止

注意:
  - 初回実行時は config.yaml を作成してください
  - Linux サーバー上で実行してください（macOS では Docker デーモンが必要）
EOF
        ;;
    *)
        log_error "不明なコマンド: $COMPOSE_CMD"
        echo "詳細: ./deploy.sh help"
        exit 1
        ;;
esac
