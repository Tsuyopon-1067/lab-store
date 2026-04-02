<script lang="ts">
    import { onMount } from "svelte";
    import { apiCallWithAuth } from "$lib/api";
    import type { AdminPurchaseWithItems } from "$lib/types";

    let purchases: AdminPurchaseWithItems[] = $state([]);
    let isLoading = $state(false);
    let errorMessage = $state("");

    // フィルタリング
    let selectedUserId = $state<number | null>(null);
    let currentPage = $state(1);
    const itemsPerPage = 50;

    onMount(async () => {
        await loadPurchases();
    });

    async function loadPurchases() {
        isLoading = true;
        errorMessage = "";
        try {
            const offset = (currentPage - 1) * itemsPerPage;
            const params = new URLSearchParams({
                limit: itemsPerPage.toString(),
                offset: offset.toString(),
            });

            if (selectedUserId !== null && selectedUserId !== -1) {
                params.append("user_id", selectedUserId.toString());
            }

            const data = await apiCallWithAuth<AdminPurchaseWithItems[]>(
                `/purchases?${params}`,
            );
            purchases = data || [];
        } catch (err) {
            errorMessage =
                err instanceof Error
                    ? err.message
                    : "購入履歴の取得に失敗しました";
        } finally {
            isLoading = false;
        }
    }

    function handleFilterChange() {
        currentPage = 1;
        loadPurchases();
    }

    function formatDate(dateString: string): string {
        return new Date(dateString).toLocaleString("ja-JP");
    }

    function formatAmount(amount: number): string {
        return `¥${amount.toLocaleString("ja-JP")}`;
    }

    function goToPreviousPage() {
        if (currentPage > 1) {
            currentPage--;
            loadPurchases();
        }
    }

    function goToNextPage() {
        currentPage++;
        loadPurchases();
    }

    async function deletePurchase(purchaseId: number) {
        if (!confirm(`購入 #${purchaseId} を削除しますか？この操作は取り消せません。`)) {
            return;
        }

        try {
            await apiCallWithAuth(`/purchases/${purchaseId}`, {
                method: "DELETE",
            });
            // Refresh the list after deletion
            await loadPurchases();
        } catch (err) {
            errorMessage =
                err instanceof Error
                    ? err.message
                    : "削除に失敗しました";
        }
    }

    async function exportToCSV() {
        try {
            const response = await fetch("/api/purchases/export", {
                method: "GET",
                credentials: "include",
            });

            if (!response.ok) {
                throw new Error("CSV export failed");
            }

            const blob = await response.blob();
            const url = window.URL.createObjectURL(blob);
            const link = document.createElement("a");
            link.href = url;
            link.download = `purchases_${new Date().toISOString().split("T")[0]}.csv`;
            document.body.appendChild(link);
            link.click();
            document.body.removeChild(link);
            window.URL.revokeObjectURL(url);
        } catch (err) {
            errorMessage =
                err instanceof Error ? err.message : "CSV出力に失敗しました";
        }
    }
</script>

<div class="purchases-page">
    <div class="header">
        <h1>購入履歴</h1>
        {#if purchases.length > 0}
            <button class="btn btn-export" onclick={exportToCSV} disabled={isLoading}>
                📥 CSV出力
            </button>
        {/if}
    </div>

    {#if errorMessage}
        <div class="alert alert-error">
            <strong>エラー:</strong>
            {errorMessage}
        </div>
    {/if}

    {#if isLoading && purchases.length === 0}
        <div class="loading">読み込み中...</div>
    {:else if purchases.length === 0}
        <div class="empty-notice">購入履歴がありません</div>
    {:else}
        <div class="table-wrapper">
            <table class="purchases-table">
                <thead>
                    <tr>
                        <th>購入ID</th>
                        <th>ユーザー</th>
                        <th>購入日時</th>
                        <th>合計金額</th>
                        <th>操作</th>
                    </tr>
                </thead>
                <tbody>
                    {#each purchases as purchase (purchase.id)}
                        <tr class="purchase-row">
                            <td>#{purchase.id}</td>
                            <td>{purchase.user_name}</td>
                            <td>{formatDate(purchase.purchased_at)}</td>
                            <td class="amount">{formatAmount(purchase.total_amount)}</td>
                            <td>
                                <button
                                    class="btn btn-delete"
                                    onclick={() => deletePurchase(purchase.id)}
                                    disabled={isLoading}
                                >
                                    削除
                                </button>
                            </td>
                        </tr>
                        <tr class="items-row">
                            <td colspan="5">
                                <details>
                                    <summary class="details-summary">商品明細を表示</summary>
                                    <div class="items-table">
                                        <div class="items-header">
                                            <div class="col-product">商品</div>
                                            <div class="col-quantity">数量</div>
                                            <div class="col-price">単価</div>
                                            <div class="col-subtotal">小計</div>
                                        </div>
                                        {#each purchase.items as item (item.product_id)}
                                            <div class="items-row-content">
                                                <div class="col-product">{item.product_name}</div>
                                                <div class="col-quantity">{item.quantity}</div>
                                                <div class="col-price">{formatAmount(item.unit_price)}</div>
                                                <div class="col-subtotal">{formatAmount(item.subtotal)}</div>
                                            </div>
                                        {/each}
                                    </div>
                                </details>
                            </td>
                        </tr>
                    {/each}
                </tbody>
            </table>
        </div>

        <div class="pagination">
            <button
                class="btn btn-secondary"
                onclick={goToPreviousPage}
                disabled={currentPage === 1 || isLoading}
            >
                ← 前へ
            </button>
            <span class="page-info">ページ {currentPage}</span>
            <button
                class="btn btn-secondary"
                onclick={goToNextPage}
                disabled={purchases.length < itemsPerPage || isLoading}
            >
                次へ →
            </button>
        </div>
    {/if}
</div>

<style>
    .purchases-page {
        max-width: 1200px;
    }

    .header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 2rem;
    }

    h1 {
        margin: 0;
        font-size: 1.8rem;
        color: #fff;
    }

    .alert {
        padding: 1rem;
        margin-bottom: 1.5rem;
        border-radius: 4px;
        font-weight: 500;
    }

    .alert-error {
        background-color: #fee;
        color: #c00;
        border: 1px solid #fcc;
    }

    .loading,
    .empty-notice {
        text-align: center;
        padding: 2rem;
        color: #eee;
        font-size: 1rem;
    }

    .table-wrapper {
        background-color: white;
        border-radius: 8px;
        border: 1px solid #ddd;
        overflow-x: auto;
        margin-bottom: 2rem;
    }

    .purchases-table {
        width: 100%;
        border-collapse: collapse;
        font-size: 0.95rem;
    }

    .purchases-table th {
        background-color: #f5f5f5;
        padding: 1rem;
        text-align: left;
        font-weight: 600;
        border-bottom: 2px solid #ddd;
        color: #2c3e50;
    }

    .purchases-table td {
        padding: 1rem;
        border-bottom: 1px solid #eee;
    }

    .purchases-table tbody .purchase-row:hover {
        background-color: #f9f9f9;
    }

    .purchases-table .amount {
        font-weight: 600;
        color: #0066cc;
    }

    .items-row {
        background-color: #fafafa;
    }

    .items-row td {
        padding: 0.5rem 1rem;
        border-bottom: 1px solid #eee;
    }

    .details-summary {
        cursor: pointer;
        color: #0066cc;
        font-weight: 500;
        user-select: none;
        padding: 0.5rem 0;
    }

    .details-summary:hover {
        text-decoration: underline;
    }

    .items-table {
        margin-top: 1rem;
        border: 1px solid #ddd;
        border-radius: 4px;
        background-color: #fff;
    }

    .items-header {
        display: grid;
        grid-template-columns: 1fr 80px 100px 100px;
        gap: 1rem;
        padding: 1rem;
        background-color: #f5f5f5;
        border-bottom: 1px solid #ddd;
        font-weight: 600;
        color: #2c3e50;
        font-size: 0.9rem;
    }

    .items-row-content {
        display: grid;
        grid-template-columns: 1fr 80px 100px 100px;
        gap: 1rem;
        padding: 0.75rem 1rem;
        border-bottom: 1px solid #eee;
        align-items: center;
        font-size: 0.9rem;
    }

    .items-row-content:last-child {
        border-bottom: none;
    }

    .col-product {
        text-align: left;
    }

    .col-quantity {
        text-align: center;
    }

    .col-price {
        text-align: right;
    }

    .col-subtotal {
        text-align: right;
        font-weight: 500;
        color: #0066cc;
    }

    .pagination {
        display: flex;
        justify-content: center;
        align-items: center;
        gap: 1rem;
    }

    .btn {
        padding: 0.5rem 1rem;
        border: none;
        border-radius: 4px;
        font-weight: 500;
        cursor: pointer;
        transition: background-color 0.2s;
    }

    .btn-secondary {
        background-color: #f0f0f0;
        color: #2c3e50;
        border: 1px solid #ddd;
    }

    .btn-secondary:hover:not(:disabled) {
        background-color: #e0e0e0;
    }

    .btn-secondary:disabled {
        opacity: 0.6;
        cursor: not-allowed;
    }

    .btn-delete {
        background-color: #fee;
        color: #c00;
        border: 1px solid #fcc;
        padding: 0.4rem 0.8rem;
        font-size: 0.9rem;
    }

    .btn-delete:hover:not(:disabled) {
        background-color: #fdd;
    }

    .btn-export {
        background-color: #0066cc;
        color: white;
        border: none;
    }

    .btn-export:hover:not(:disabled) {
        background-color: #0052a3;
    }

    .btn-export:disabled {
        opacity: 0.6;
        cursor: not-allowed;
    }

    .page-info {
        font-weight: 500;
        color: #2c3e50;
    }

    @media (max-width: 768px) {
        .purchases-table {
            font-size: 0.85rem;
        }

        .purchases-table th,
        .purchases-table td {
            padding: 0.75rem;
        }

        .items-header {
            grid-template-columns: 1fr 60px 80px 80px;
            gap: 0.5rem;
            padding: 0.75rem;
            font-size: 0.85rem;
        }

        .items-row-content {
            grid-template-columns: 1fr 60px 80px 80px;
            gap: 0.5rem;
            padding: 0.5rem 0.75rem;
        }

        .pagination {
            flex-wrap: wrap;
        }
    }
</style>
