# クイックスタートガイド

## 5分で本番環境を起動する

### 前提条件

- Linux サーバー（Ubuntu 22.04 LTS以上推奨）
- Podman 4.0+ または Docker 20.10+
- sudo 権限

### ステップ 1: リポジトリをクローン

```bash
git clone https://github.com/your-org/purchase-system.git /opt/purchase-system
cd /opt/purchase-system
```

### ステップ 2: 設定ファイルを作成

```bash
cp deploy/config.yaml.template config.yaml
```

`config.yaml` を本番環境に合わせて編集（オプション）：

```bash
nano config.yaml
```

デフォルト設定でも動作します。

### ステップ 3: デプロイスクリプトで起動

```bash
./deploy.sh up
```

このコマンドが以下を自動実行します：
- ✅ 前提条件チェック
- ✅ ディレクトリ作成（data/, backup/）
- ✅ コンテナのビルド
- ✅ コンテナの起動
- ✅ ヘルスチェック

### ステップ 4: ブラウザでアクセス

```
http://localhost:3000
```

完了です！🎉

---

## 初回ログイン

1. ブラウザで `http://localhost:3000` を開く
2. **管理者ログイン** をクリック
3. パスワードを入力（初回セットアップで設定したもの）
4. ログイン成功

---

## よく使うコマンド

```bash
# ログを表示（リアルタイム）
./deploy.sh logs

# コンテナを再起動
./deploy.sh restart

# コンテナを停止
./deploy.sh down

# ヘルスチェック
./deploy.sh health

# コンテナの状態確認
./deploy.sh ps
```

---

## トラブルシューティング

### Q: `http://localhost:3000` にアクセスできない

**A:** 以下を確認してください：

```bash
# コンテナが起動しているか確認
./deploy.sh ps

# ログでエラーを確認
./deploy.sh logs

# ファイアウォールでポート3000が開いているか確認
sudo firewall-cmd --list-all  # CentOS/RHEL の場合
sudo ufw status               # Ubuntu/Debian の場合
```

### Q: パスワードを忘れた

**A:** 管理者パスワードをリセットできます：

```bash
./deploy.sh setup
```

### Q: データベースが破損した

**A:** バックアップからリストア：

```bash
# 既存のDBをバックアップ
cp data/purchase.db data/purchase.db.backup

# バックアップファイルをリストア
cp backup/manual_20260330_120000.sqlite data/purchase.db

# コンテナを再起動
./deploy.sh restart
```

---

## systemd に登録して自動起動

Linux PC起動時に自動起動したい場合：

```bash
# サービスファイルをコピー
sudo cp deploy/purchase-system.service /etc/systemd/system/

# systemd をリロード
sudo systemctl daemon-reload

# 自動起動を有効化
sudo systemctl enable purchase-system

# サービスを開始
sudo systemctl start purchase-system

# 状態確認
sudo systemctl status purchase-system
```

---

## 詳細ドキュメント

- [本番環境デプロイメントガイド](docs/production_deployment.md)
- [デプロイチェックリスト](docs/production_checklist.md)
- [コンテナ設計](docs/container_design.md)
- [要件定義書](docs/requirements.md)

---

## サポート

問題が発生した場合：

1. `./deploy.sh logs` でエラーメッセージを確認
2. [本番環境デプロイメントガイド](docs/production_deployment.md#10-トラブルシューティング) を参照
3. Issue を報告
