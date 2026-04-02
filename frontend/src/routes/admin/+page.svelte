<script lang="ts">
    import { auth } from "$lib/stores/auth";
</script>

<div class="dashboard">
    <h1>管理ダッシュボード</h1>

    <div class="dashboard-grid">
        <div class="card">
            <h3>利用者管理</h3>
            <p>利用者を登録・編集・削除する</p>
            <a href="/admin/users" class="btn btn-primary">編集</a>
        </div>

        <div class="card">
            <h3>商品管理</h3>
            <p>商品マスターの追加・編集・価格変更</p>
            <a href="/admin/products" class="btn btn-primary">編集</a>
        </div>

        <div class="card">
            <h3>購入履歴</h3>
            <p>全利用者の購入履歴を確認</p>
            <a href="/admin/purchases" class="btn btn-primary">確認</a>
        </div>

        <div class="card">
            <h3>仕入れ履歴</h3>
            <p>仕入れ記録を確認・管理</p>
            <a href="/admin/restocks" class="btn btn-primary">確認</a>
        </div>

        <div class="card">
            <h3>支払い管理</h3>
            <p>購入代金の支払い記録を管理</p>
            <a href="/admin/payments" class="btn btn-primary">管理</a>
        </div>

        <div class="card">
            <h3>立替精算</h3>
            <p>仕入れ立替の精算記録を管理</p>
            <a href="/admin/restock-payments" class="btn btn-primary">管理</a>
        </div>

        <div class="card">
            <h3>月次精算</h3>
            <p>各利用者の収支サマリを確認</p>
            <a href="/admin/summary" class="btn btn-primary">確認</a>
        </div>

        <div class="card">
            <h3>支払い金額一覧</h3>
            <p>利用者ごとの支払い金額を確認</p>
            <a href="/admin/user-payments" class="btn btn-primary">確認</a>
        </div>

        <div class="card">
            <h3>バックアップ</h3>
            <p>手動バックアップを実行</p>
            <a href="/admin/backup" class="btn btn-primary">実行</a>
        </div>

        <div class="card">
            <h3>APIキー管理</h3>
            <p>外部API用のAPIキーを発行・管理する</p>
            <a href="/admin/api-keys" class="btn btn-primary">管理</a>
        </div>

        <div class="card">
            <h3>設定</h3>
            <p>システム設定を管理</p>
            <a href="/admin/settings" class="btn btn-primary">設定</a>
        </div>
    </div>
</div>

<style>
    .dashboard {
        max-width: 1200px;
    }

    h1 {
        margin-bottom: 2rem;
        font-size: 2rem;
        color: #fff;
    }

    .dashboard-grid {
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
        gap: 1.5rem;
    }

    .card {
        background-color: white;
        border: 1px solid #ddd;
        border-radius: 8px;
        padding: 1.5rem;
        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
    }

    .card h3 {
        margin-bottom: 0.5rem;
        color: #2c3e50;
    }

    .card p {
        margin-bottom: 1rem;
        color: #666;
        font-size: 0.95rem;
    }

    .btn {
        display: inline-block;
        padding: 0.5rem 1rem;
        border-radius: 4px;
        text-decoration: none;
        text-align: center;
        cursor: pointer;
        font-weight: 500;
        transition: background-color 0.2s;
    }

    .btn-primary {
        background-color: #0066cc;
        color: white;
    }

    .btn-primary:hover {
        background-color: #0052a3;
        text-decoration: none;
    }
</style>
