# ディレクトリ構成・DB初期化設計書

**文書バージョン**: 1.3  
**作成日**: 2026-03-28

---

## 1. リポジトリ全体構成

```
purchase-system/
├── backend/
├── frontend/
├── data/
├── backup/
├── scripts/
├── deploy/
├── .gitignore
└── README.md
```

---

## 2. バックエンド（Go）ディレクトリ構成

```
backend/
├── main.go
├── go.mod
├── go.sum
│
├── config/
│   └── config.go
│
├── db/
│   ├── db.go
│   └── migrations/
│       ├── 001_create_users.sql
│       ├── 002_create_products.sql
│       ├── 003_create_product_prices.sql
│       ├── 004_create_purchases.sql
│       ├── 005_create_purchase_items.sql
│       ├── 006_create_payments.sql
│       ├── 007_create_restocks.sql
│       ├── 008_create_restock_items.sql
│       ├── 009_create_restock_payments.sql
│       ├── 010_create_audit_logs.sql
│       ├── 011_create_admin.sql
│       ├── 012_create_sessions.sql
│       └── 013_create_backups.sql
│
├── handler/
│   ├── auth.go
│   ├── user.go
│   ├── product.go
│   ├── purchase.go
│   ├── payment.go           # 購入代金支払い記録
│   ├── restock.go           # 仕入れ登録・履歴
│   ├── restock_payment.go   # 仕入れ立替精算記録
│   ├── summary.go           # 月次精算サマリ
│   ├── me.go                # 一般利用者向け自分の情報
│   └── backup.go
│
├── middleware/
│   ├── auth.go              # 管理者セッション検証
│   └── barcode.go           # X-User-Barcodeヘッダで利用者特定
│
├── model/
│   ├── user.go
│   ├── product.go
│   ├── purchase.go
│   ├── payment.go
│   ├── restock.go
│   ├── restock_payment.go
│   ├── audit_log.go
│   └── admin.go
│
├── repository/
│   ├── user.go
│   ├── product.go
│   ├── purchase.go
│   ├── payment.go
│   ├── restock.go
│   ├── restock_payment.go
│   ├── audit_log.go
│   └── backup.go
│
├── service/
│   ├── user.go
│   ├── product.go
│   ├── purchase.go
│   ├── payment.go           # 支払い記録・修正時に audit_log も書く
│   ├── restock.go
│   ├── restock_payment.go   # 精算記録・修正時に audit_log も書く
│   ├── summary.go           # 残高集計ロジック
│   └── backup.go
│
└── cmd/
    └── setup/
        └── main.go
```

---

## 3. フロントエンド（SvelteKit）ディレクトリ構成

```
frontend/
├── package.json
├── svelte.config.js
├── vite.config.js
│
├── src/
│   ├── app.html
│   │
│   ├── lib/
│   │   ├── api/
│   │   │   ├── users.ts
│   │   │   ├── products.ts
│   │   │   ├── purchases.ts
│   │   │   ├── payments.ts
│   │   │   ├── restocks.ts
│   │   │   ├── restock_payments.ts
│   │   │   ├── summary.ts
│   │   │   └── me.ts
│   │   │
│   │   ├── components/
│   │   │   ├── Cart.svelte
│   │   │   ├── BarcodeScanner.svelte
│   │   │   ├── ProductSearch.svelte
│   │   │   ├── ConfirmDialog.svelte
│   │   │   ├── PurchaseHistory.svelte
│   │   │   ├── RestockForm.svelte       # 仕入れ登録フォーム
│   │   │   └── BalanceSummary.svelte    # 残高サマリ表示
│   │   │
│   │   └── stores/
│   │       ├── cart.ts
│   │       ├── restockCart.ts           # 仕入れカート
│   │       └── session.ts
│   │
│   └── routes/
│       ├── +layout.svelte
│       ├── +page.svelte                 # 購買画面（メイン）
│       │
│       ├── restock/
│       │   └── +page.svelte             # 仕入れ画面
│       │
│       ├── history/
│       │   └── +page.svelte             # 残高・履歴確認画面（バーコード認証）
│       │
│       └── admin/
│           ├── +layout.svelte
│           ├── login/
│           │   └── +page.svelte
│           ├── products/
│           │   └── +page.svelte
│           ├── users/
│           │   └── +page.svelte
│           ├── purchases/
│           │   └── +page.svelte         # 購入履歴
│           ├── restocks/
│           │   └── +page.svelte         # 仕入れ履歴・修正・削除
│           ├── payments/
│           │   └── +page.svelte         # 購入代金支払い管理
│           ├── restock-payments/
│           │   └── +page.svelte         # 仕入れ立替精算管理
│           ├── summary/
│           │   └── +page.svelte         # 月次精算サマリ
│           ├── backup/
│           │   └── +page.svelte
│           └── settings/
│               └── +page.svelte
│
└── static/
    └── favicon.png
```

---

## 4. DBマイグレーション（goose）

### db.go での組み込み

```go
package db

import (
    "database/sql"
    "embed"

    "github.com/pressly/goose/v3"
    _ "modernc.org/sqlite"
)

//go:embed migrations/*.sql
var migrations embed.FS

func Init(dbPath string) (*sql.DB, error) {
    db, err := sql.Open("sqlite", dbPath)
    if err != nil {
        return nil, err
    }
    if _, err := db.Exec("PRAGMA journal_mode=WAL;"); err != nil {
        return nil, err
    }
    goose.SetBaseFS(migrations)
    if err := goose.SetDialect("sqlite3"); err != nil {
        return nil, err
    }
    if err := goose.Up(db, "migrations"); err != nil {
        return nil, err
    }
    return db, nil
}
```

---

## 5. SQLマイグレーションファイル

### 001_create_users.sql

```sql
-- +goose Up
CREATE TABLE users (
    id         INTEGER  PRIMARY KEY AUTOINCREMENT,
    name       TEXT     NOT NULL,
    barcode    TEXT     NOT NULL UNIQUE,
    is_active  INTEGER  NOT NULL DEFAULT 1,
    created_at DATETIME NOT NULL DEFAULT (datetime('now', 'localtime'))
);

-- +goose Down
DROP TABLE users;
```

### 002_create_products.sql

```sql
-- +goose Up
CREATE TABLE products (
    id         INTEGER  PRIMARY KEY AUTOINCREMENT,
    name       TEXT     NOT NULL,
    barcode    TEXT     NOT NULL UNIQUE,
    is_active  INTEGER  NOT NULL DEFAULT 1,
    note       TEXT,
    created_at DATETIME NOT NULL DEFAULT (datetime('now', 'localtime'))
);

-- +goose Down
DROP TABLE products;
```

### 003_create_product_prices.sql

```sql
-- +goose Up
CREATE TABLE product_prices (
    id         INTEGER  PRIMARY KEY AUTOINCREMENT,
    product_id INTEGER  NOT NULL REFERENCES products(id),
    price      INTEGER  NOT NULL,
    valid_from DATETIME NOT NULL DEFAULT (datetime('now', 'localtime')),
    valid_to   DATETIME
);

CREATE INDEX idx_product_prices_product_id ON product_prices(product_id);

-- +goose Down
DROP INDEX idx_product_prices_product_id;
DROP TABLE product_prices;
```

### 004_create_purchases.sql

```sql
-- +goose Up
CREATE TABLE purchases (
    id           INTEGER  PRIMARY KEY AUTOINCREMENT,
    user_id      INTEGER  NOT NULL REFERENCES users(id),
    purchased_at DATETIME NOT NULL DEFAULT (datetime('now', 'localtime'))
);

CREATE INDEX idx_purchases_user_id      ON purchases(user_id);
CREATE INDEX idx_purchases_purchased_at ON purchases(purchased_at);

-- +goose Down
DROP INDEX idx_purchases_purchased_at;
DROP INDEX idx_purchases_user_id;
DROP TABLE purchases;
```

### 005_create_purchase_items.sql

```sql
-- +goose Up
CREATE TABLE purchase_items (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    purchase_id INTEGER NOT NULL REFERENCES purchases(id),
    product_id  INTEGER NOT NULL REFERENCES products(id),
    quantity    INTEGER NOT NULL DEFAULT 1,
    unit_price  INTEGER NOT NULL
);

CREATE INDEX idx_purchase_items_purchase_id ON purchase_items(purchase_id);

-- +goose Down
DROP INDEX idx_purchase_items_purchase_id;
DROP TABLE purchase_items;
```

### 006_create_payments.sql

```sql
-- +goose Up
-- 購入代金の支払い記録（利用者 → コミュニティへの現金支払い）
CREATE TABLE payments (
    id          INTEGER  PRIMARY KEY AUTOINCREMENT,
    user_id     INTEGER  NOT NULL REFERENCES users(id),
    amount      INTEGER  NOT NULL,
    paid_at     DATETIME NOT NULL,
    note        TEXT,
    created_by  TEXT     NOT NULL DEFAULT 'admin',
    updated_by  TEXT,
    deleted_at  DATETIME,            -- 論理削除用（NULLなら有効）
    created_at  DATETIME NOT NULL DEFAULT (datetime('now', 'localtime')),
    updated_at  DATETIME
);

CREATE INDEX idx_payments_user_id ON payments(user_id);

-- +goose Down
DROP INDEX idx_payments_user_id;
DROP TABLE payments;
```

### 007_create_restocks.sql

```sql
-- +goose Up
-- 仕入れヘッダ（仕入れ者が立て替えた記録）
CREATE TABLE restocks (
    id            INTEGER  PRIMARY KEY AUTOINCREMENT,
    user_id       INTEGER  NOT NULL REFERENCES users(id),  -- 仕入れ者（立替者）
    total_amount  INTEGER  NOT NULL,                        -- 実際の支払い合計（円）
    restocked_at  DATETIME NOT NULL,
    note          TEXT,
    deleted_at    DATETIME,
    created_at    DATETIME NOT NULL DEFAULT (datetime('now', 'localtime')),
    updated_at    DATETIME
);

CREATE INDEX idx_restocks_user_id      ON restocks(user_id);
CREATE INDEX idx_restocks_restocked_at ON restocks(restocked_at);

-- +goose Down
DROP INDEX idx_restocks_restocked_at;
DROP INDEX idx_restocks_user_id;
DROP TABLE restocks;
```

### 008_create_restock_items.sql

```sql
-- +goose Up
-- 仕入れ明細（どの商品を何個仕入れたか）
CREATE TABLE restock_items (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    restock_id  INTEGER NOT NULL REFERENCES restocks(id),
    product_id  INTEGER NOT NULL REFERENCES products(id),
    quantity    INTEGER NOT NULL DEFAULT 1,
    unit_price  INTEGER NOT NULL   -- 仕入れ単価（円）
);

CREATE INDEX idx_restock_items_restock_id ON restock_items(restock_id);

-- +goose Down
DROP INDEX idx_restock_items_restock_id;
DROP TABLE restock_items;
```

### 009_create_restock_payments.sql

```sql
-- +goose Up
-- 仕入れ立替精算記録（コミュニティ → 仕入れ者への支払い）
CREATE TABLE restock_payments (
    id          INTEGER  PRIMARY KEY AUTOINCREMENT,
    user_id     INTEGER  NOT NULL REFERENCES users(id),  -- 仕入れ者（受取人）
    amount      INTEGER  NOT NULL,
    settled_at  DATETIME NOT NULL,
    note        TEXT,
    created_by  TEXT     NOT NULL DEFAULT 'admin',
    updated_by  TEXT,
    deleted_at  DATETIME,
    created_at  DATETIME NOT NULL DEFAULT (datetime('now', 'localtime')),
    updated_at  DATETIME
);

CREATE INDEX idx_restock_payments_user_id ON restock_payments(user_id);

-- +goose Down
DROP INDEX idx_restock_payments_user_id;
DROP TABLE restock_payments;
```

### 010_create_audit_logs.sql

```sql
-- +goose Up
-- 支払い記録・精算記録の修正・削除の操作ログ
CREATE TABLE audit_logs (
    id          INTEGER  PRIMARY KEY AUTOINCREMENT,
    table_name  TEXT     NOT NULL,   -- 'payments' or 'restock_payments' or 'restocks'
    record_id   INTEGER  NOT NULL,
    action      TEXT     NOT NULL,   -- 'update' or 'delete'
    before_json TEXT     NOT NULL,   -- 変更前のレコード内容（JSON）
    after_json  TEXT,                -- 変更後のレコード内容（JSON）、削除時はNULL
    operated_at DATETIME NOT NULL DEFAULT (datetime('now', 'localtime'))
);

CREATE INDEX idx_audit_logs_table_record ON audit_logs(table_name, record_id);

-- +goose Down
DROP INDEX idx_audit_logs_table_record;
DROP TABLE audit_logs;
```

### 011_create_admin.sql

```sql
-- +goose Up
CREATE TABLE admin (
    id            INTEGER  PRIMARY KEY AUTOINCREMENT,
    password_hash TEXT     NOT NULL,
    updated_at    DATETIME NOT NULL DEFAULT (datetime('now', 'localtime'))
);

-- +goose Down
DROP TABLE admin;
```

### 012_create_sessions.sql

```sql
-- +goose Up
CREATE TABLE admin_sessions (
    id         INTEGER  PRIMARY KEY AUTOINCREMENT,
    token      TEXT     NOT NULL UNIQUE,
    created_at DATETIME NOT NULL DEFAULT (datetime('now', 'localtime')),
    expires_at DATETIME NOT NULL
);

CREATE INDEX idx_admin_sessions_token ON admin_sessions(token);

-- +goose Down
DROP INDEX idx_admin_sessions_token;
DROP TABLE admin_sessions;
```

### 013_create_backups.sql

```sql
-- +goose Up
CREATE TABLE backups (
    id         INTEGER  PRIMARY KEY AUTOINCREMENT,
    filename   TEXT     NOT NULL,
    created_at DATETIME NOT NULL DEFAULT (datetime('now', 'localtime'))
);

-- +goose Down
DROP TABLE backups;
```

---

## 6. 主要クエリ例

### バーコードで有効商品と現在価格を取得

```sql
SELECT p.id, p.name, p.barcode, pp.price
FROM products p
JOIN product_prices pp ON pp.product_id = p.id
WHERE p.barcode   = ?
  AND p.is_active = 1
  AND pp.valid_to IS NULL;
```

### 利用者の残高サマリを算出

```sql
-- 購入未払い額
SELECT COALESCE(SUM(pi.quantity * pi.unit_price), 0)
FROM purchases pu
JOIN purchase_items pi ON pi.purchase_id = pu.id
WHERE pu.user_id = ?;

-- 購入既払い額
SELECT COALESCE(SUM(amount), 0)
FROM payments
WHERE user_id = ? AND deleted_at IS NULL;

-- 仕入れ立替合計額
SELECT COALESCE(SUM(total_amount), 0)
FROM restocks
WHERE user_id = ? AND deleted_at IS NULL;

-- 仕入れ既精算額
SELECT COALESCE(SUM(amount), 0)
FROM restock_payments
WHERE user_id = ? AND deleted_at IS NULL;
```

### 月次精算サマリ（全利用者）

```sql
SELECT
    u.id,
    u.name,
    COALESCE(SUM(pi.quantity * pi.unit_price), 0)  AS purchase_total,
    COALESCE(pay.paid, 0)                           AS purchase_paid,
    COALESCE(SUM(pi.quantity * pi.unit_price), 0)
        - COALESCE(pay.paid, 0)                     AS purchase_unpaid,
    COALESCE(rs.restock_total, 0)                   AS restock_total,
    COALESCE(rp.settled, 0)                         AS restock_settled,
    COALESCE(rs.restock_total, 0)
        - COALESCE(rp.settled, 0)                   AS restock_unclaimed
FROM users u
LEFT JOIN purchases pu ON pu.user_id = u.id
    AND pu.purchased_at BETWEEN ? AND ?
LEFT JOIN purchase_items pi ON pi.purchase_id = pu.id
LEFT JOIN (
    SELECT user_id, SUM(amount) AS paid
    FROM payments
    WHERE paid_at BETWEEN ? AND ? AND deleted_at IS NULL
    GROUP BY user_id
) pay ON pay.user_id = u.id
LEFT JOIN (
    SELECT user_id, SUM(total_amount) AS restock_total
    FROM restocks
    WHERE restocked_at BETWEEN ? AND ? AND deleted_at IS NULL
    GROUP BY user_id
) rs ON rs.user_id = u.id
LEFT JOIN (
    SELECT user_id, SUM(amount) AS settled
    FROM restock_payments
    WHERE settled_at BETWEEN ? AND ? AND deleted_at IS NULL
    GROUP BY user_id
) rp ON rp.user_id = u.id
WHERE u.is_active = 1
GROUP BY u.id;
```

---

## 7. Gin ルーター構成（main.go）

```go
r := gin.Default()

// 認証不要
r.GET("/api/users/barcode/:code",    handler.GetUserByBarcode(db))
r.GET("/api/products/barcode/:code", handler.GetProductByBarcode(db))
r.GET("/api/products",               handler.ListProducts(db))
r.POST("/api/purchases",             handler.CreatePurchase(db))
r.POST("/api/restocks",              handler.CreateRestock(db))   // 仕入れ登録

// 一般利用者（X-User-Barcodeヘッダで本人確認）
me := r.Group("/api/me", middleware.BarcodeAuth(db))
{
    me.GET("/purchases", handler.GetMyPurchases(db))
    me.GET("/restocks",  handler.GetMyRestocks(db))
    me.GET("/balance",   handler.GetMyBalance(db))
}

// 管理者（セッションCookie必須）
admin := r.Group("/api", middleware.AdminAuth(db))
{
    admin.POST("/auth/login",  handler.Login(db))
    admin.POST("/auth/logout", handler.Logout(db))

    admin.GET("/users",            handler.ListUsers(db))
    admin.POST("/users",           handler.CreateUser(db))
    admin.PUT("/users/:id",        handler.UpdateUser(db))
    admin.DELETE("/users/:id",     handler.DeleteUser(db))

    admin.POST("/products",                  handler.CreateProduct(db))
    admin.PUT("/products/:id",               handler.UpdateProduct(db))
    admin.POST("/products/:id/change-price", handler.ChangePrice(db))
    admin.GET("/products/:id/prices",        handler.GetPriceHistory(db))

    admin.GET("/purchases",         handler.ListPurchases(db))
    admin.GET("/purchases/summary", handler.GetSummary(db))

    admin.GET("/restocks",         handler.ListRestocks(db))
    admin.PUT("/restocks/:id",     handler.UpdateRestock(db))
    admin.DELETE("/restocks/:id",  handler.DeleteRestock(db))

    admin.POST("/payments",          handler.CreatePayment(db))
    admin.GET("/payments",           handler.ListPayments(db))
    admin.PUT("/payments/:id",       handler.UpdatePayment(db))
    admin.DELETE("/payments/:id",    handler.DeletePayment(db))

    admin.POST("/restock-payments",        handler.CreateRestockPayment(db))
    admin.GET("/restock-payments",         handler.ListRestockPayments(db))
    admin.PUT("/restock-payments/:id",     handler.UpdateRestockPayment(db))
    admin.DELETE("/restock-payments/:id",  handler.DeleteRestockPayment(db))

    admin.POST("/backup",     handler.CreateBackup(db))
    admin.GET("/backup/list", handler.ListBackups(db))
}
```

---

## 8. go.mod 主要依存パッケージ

```
github.com/gin-gonic/gin          # HTTPフレームワーク
github.com/pressly/goose/v3       # DBマイグレーション
modernc.org/sqlite                 # SQLiteドライバ（CGO不要）
golang.org/x/crypto/bcrypt         # パスワードハッシュ
gopkg.in/yaml.v3                   # 設定ファイル読み込み
```

---

## 9. .gitignore

```
data/
backup/
server
frontend/.svelte-kit/
frontend/build/
config.yaml
frontend/node_modules/
```
