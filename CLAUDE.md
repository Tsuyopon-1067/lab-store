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

## コーディングパターン

### バックエンド（Go/Gin）

#### ハンドラー（handler/）
```go
// クロージャーパターン：db/configを受け取り gin.HandlerFunc を返す
func CreateProduct(db *sql.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req CreateProductRequest
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        
        repo := repository.NewProductRepository(db)
        // 処理...
        
        c.JSON(http.StatusOK, result)
    }
}
```
- 全ハンドラーはクロージャー形式（db/cfg を受け取る）
- リクエストバリデーション：`c.ShouldBindJSON(&req)`、`binding:"required"` タグ使用
- エラーレスポンス：`gin.H{"error": "message"}` で統一
- 成功時：対応する HTTP ステータス（201 Created、200 OK 等）

#### リポジトリ（repository/）
```go
type UserRepository struct { db *sql.DB }

func NewUserRepository(db *sql.DB) *UserRepository { ... }

func (r *UserRepository) Get(id int) (*model.User, error) { ... }
func (r *UserRepository) Create(req *model.CreateUserRequest) (*model.User, error) { ... }
```
- 構造体に `db *sql.DB` のみを持つ
- `New*Repository` コンストラクタで初期化
- メソッドレシーバーは ポインタレシーバー `(r *XxxRepository)`

#### ミドルウェア（middleware/）
```go
func AdminAuth(db *sql.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 認証ロジック
        if err := verify() {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "..."})
            c.Abort()
            return
        }
        c.Set("key", value)
        c.Next()
    }
}
```
- ハンドラーと同じクロージャーパターン
- 認証失敗時は `c.Abort()` + JSON レスポンス
- 成功時は `c.Set()` で値を格納してから `c.Next()`

#### マイグレーション（db/migrations/）
```sql
-- +goose Up
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    ...
);

-- +goose Down
DROP TABLE users;
```
- goose 形式（`-- +goose Up/Down` で囲む）
- 連番ファイル命名（`001_*.sql`, `002_*.sql` ...）

### フロントエンド（SvelteKit）

#### ページコンポーネント（routes/admin/）
```svelte
<script lang="ts">
    import { onMount } from 'svelte';
    import { apiCallWithAuth } from '$lib/api';
    import type { Product } from '$lib/types';

    let products = $state([]);
    let isLoading = $state(true);
    let errorMessage = $state('');

    onMount(async () => {
        await loadData();
    });

    async function loadData() {
        isLoading = true;
        errorMessage = '';
        try {
            const data = await apiCallWithAuth<Product[]>('/products');
            products = data;
        } catch (err) {
            errorMessage = err instanceof Error ? err.message : 'Failed to load';
        } finally {
            isLoading = false;
        }
    }

    async function handleCreate() {
        try {
            const result = await apiCallWithAuth<Product>('/products', {
                method: 'POST',
                body: JSON.stringify({ name: '...', ... })
            });
            // 成功後は loadData() で再取得
            await loadData();
        } catch (err) {
            errorMessage = err instanceof Error ? err.message : 'Failed to create';
        }
    }
</script>

<div class="page">
    {#if isLoading}
        <div>読み込み中...</div>
    {/if}
    
    <!-- コンテンツ -->
</div>

<style>
    /* スコープ付きスタイル */
</style>
```
- Svelte 5 runes：`$state()`, `$effect()` を使用
- API呼び出し：`apiCallWithAuth()` を使用（セッションCookie自動付加）
- 非認証 API（バーコード）：`apiCallWithBarcode(endpoint, barcode, options)`
- 状態管理：`isLoading`, `errorMessage` など `$state` で宣言
- `onMount` で初期データ読み込み
- エラーは `try-catch` でキャッチして `errorMessage` に格納
- スタイルは各ファイルにスコープ付き `<style>` で定義（外部 CSS ライブラリなし）

#### 型定義（lib/types.ts）
```typescript
export interface Product {
    id: number;
    name: string;
    barcode: string;
    is_active: boolean;
    current_price: number;
    created_at: string;
    note?: string;
}

export interface CreateProductRequest {
    name: string;
    barcode: string;
    price: number;
    note?: string;
}
```
- バックエンドの struct フィールド名を snake_case で対応

#### API呼び出し（lib/api.ts）
```typescript
// 管理者（セッション）
const data = await apiCallWithAuth<ResponseType>('/endpoint', {
    method: 'POST',
    body: JSON.stringify(req)
});

// 一般利用者（バーコード）
const data = await apiCallWithBarcode<ResponseType>('/endpoint', barcodeValue, {
    method: 'GET'
});
```
- `apiCallWithAuth` は Cookie に格納された admin_session を自動付加
- `apiCallWithBarcode` は `X-User-Barcode` ヘッダを自動付加
- エラー時は Error を throw（呼び出し元で catch）

---

## Git ワークフロー

### ブランチ戦略
- 機能追加・バグ修正は新しいブランチを作成して行う（例: `feature/settings-management`, `fix/login-timeout`）
- main ブランチへのマージは PR を通じて行う

### コミットメッセージ
- **必ず英語で作成する**
- 1 行目：変更内容の要約（70 文字以内）
- 2 行目以降：詳細な説明（必要に応じて）
- 末尾に署名を含める：`Co-Authored-By: Claude Haiku 4.5 <noreply@anthropic.com>`

例：
```
add settings table migration

Create a single-row settings table to store system configuration.

Co-Authored-By: Claude Haiku 4.5 <noreply@anthropic.com>
```

### コミット粒度
- 1 つのコミットは 1 つの論理的な変更を表す
- マイグレーション、モデル、リポジトリ、ハンドラなど機能を構成する各要素を段階的にコミットする
- テストやドキュメント更新は対応する機能変更と同じコミットまたは別のコミットに分ける
- コミットは小さく保つ（レビュー可能なサイズ）

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
