<script lang="ts">
    import { onMount } from "svelte";
    import { apiCallWithAuth } from "$lib/api";
    import type { SummaryReport } from "$lib/types";

    let summary: SummaryReport | null = $state(null);
    let isLoading = $state(false);
    let errorMessage = $state("");

    // 日付フィルタ
    let fromDate = $state("");
    let toDate = $state("");

    onMount(() => {
        // デフォルト: 今月の初日〜末日
        const today = new Date();
        const firstDay = new Date(today.getFullYear(), today.getMonth(), 1);
        const lastDay = new Date(today.getFullYear(), today.getMonth() + 1, 0);

        fromDate = formatDateForInput(firstDay);
        toDate = formatDateForInput(lastDay);

        loadSummary();
    });

    function formatDateForInput(date: Date): string {
        const year = date.getFullYear();
        const month = String(date.getMonth() + 1).padStart(2, "0");
        const day = String(date.getDate()).padStart(2, "0");
        return `${year}-${month}-${day}`;
    }

    async function loadSummary() {
        isLoading = true;
        errorMessage = "";
        try {
            const params = new URLSearchParams();
            if (fromDate) params.append("from", fromDate);
            if (toDate) params.append("to", toDate);

            const data = await apiCallWithAuth<SummaryReport>(
                `/purchases/summary?${params}`,
            );
            summary = data;
        } catch (err) {
            errorMessage =
                err instanceof Error
                    ? err.message
                    : "月次精算サマリの取得に失敗しました";
        } finally {
            isLoading = false;
        }
    }

    function handleSearch() {
        loadSummary();
    }

    function formatCurrency(amount: number): string {
        return `¥${amount.toLocaleString("ja-JP")}`;
    }

    function getBalanceClass(balance: number): string {
        if (balance > 0) return "balance-positive";
        if (balance < 0) return "balance-negative";
        return "balance-zero";
    }

    function getBalanceLabel(balance: number): string {
        if (balance > 0) return `${formatCurrency(balance)} (支払い)`;
        if (balance < 0) return `${formatCurrency(Math.abs(balance))} (受取)`;
        return "0円";
    }
</script>

<div class="summary-page">
    <h1>月次精算サマリ</h1>

    {#if errorMessage}
        <div class="alert alert-error">
            <strong>エラー:</strong>
            {errorMessage}
        </div>
    {/if}

    <div class="filter-section">
        <div class="filter-controls">
            <div class="date-range">
                <div class="date-input">
                    <label for="from-date">期間開始日</label>
                    <input
                        type="date"
                        id="from-date"
                        bind:value={fromDate}
                        disabled={isLoading}
                    />
                </div>
                <span class="separator">〜</span>
                <div class="date-input">
                    <label for="to-date">期間終了日</label>
                    <input
                        type="date"
                        id="to-date"
                        bind:value={toDate}
                        disabled={isLoading}
                    />
                </div>
            </div>
            <button
                class="btn btn-primary"
                onclick={handleSearch}
                disabled={isLoading}
            >
                {isLoading ? "読み込み中..." : "検索"}
            </button>
        </div>
    </div>

    {#if isLoading && !summary}
        <div class="loading">読み込み中...</div>
    {:else if summary}
        <div class="summary-container">
            <div class="period-info">
                期間: {summary.period.from} 〜 {summary.period.to}
            </div>

            {#if summary.users.length === 0}
                <div class="empty-notice">この期間のデータがありません</div>
            {:else}
                <div class="table-wrapper">
                    <table class="summary-table">
                        <thead>
                            <tr>
                                <th>ユーザー</th>
                                <th>購入合計</th>
                                <th>購入既払い</th>
                                <th>購入未払い</th>
                                <th>仕入れ合計</th>
                                <th>仕入れ既精算</th>
                                <th>仕入れ未精算</th>
                                <th>差引請求額</th>
                            </tr>
                        </thead>
                        <tbody>
                            {#each summary.users as user (user.user_id)}
                                <tr>
                                    <td class="user-name">{user.user_name}</td>
                                    <td class="amount-cell"
                                        >{formatCurrency(
                                            user.purchase_total,
                                        )}</td
                                    >
                                    <td class="amount-cell"
                                        >{formatCurrency(
                                            user.purchase_paid,
                                        )}</td
                                    >
                                    <td class="amount-cell unpaid"
                                        >{formatCurrency(
                                            user.purchase_unpaid,
                                        )}</td
                                    >
                                    <td class="amount-cell"
                                        >{formatCurrency(
                                            user.restock_total,
                                        )}</td
                                    >
                                    <td class="amount-cell"
                                        >{formatCurrency(
                                            user.restock_settled,
                                        )}</td
                                    >
                                    <td class="amount-cell unclaimed"
                                        >{formatCurrency(
                                            user.restock_unclaimed,
                                        )}</td
                                    >
                                    <td
                                        class={`amount-cell balance-cell ${getBalanceClass(user.net_balance)}`}
                                    >
                                        <strong
                                            >{getBalanceLabel(
                                                user.net_balance,
                                            )}</strong
                                        >
                                    </td>
                                </tr>
                            {/each}
                        </tbody>
                    </table>
                </div>

                <div class="legend">
                    <div class="legend-item">
                        <span class="legend-label legend-positive"
                            >支払い義務</span
                        >
                        <p>ユーザーがコミュニティに支払う必要がある金額</p>
                    </div>
                    <div class="legend-item">
                        <span class="legend-label legend-negative"
                            >受取権利</span
                        >
                        <p>コミュニティがユーザーに支払う必要がある金額</p>
                    </div>
                </div>
            {/if}
        </div>
    {/if}
</div>

<style>
    .summary-page {
        max-width: 1400px;
    }

    h1 {
        margin: 0 0 2rem 0;
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

    .filter-section {
        background-color: white;
        padding: 1.5rem;
        border-radius: 8px;
        border: 1px solid #ddd;
        margin-bottom: 2rem;
    }

    .filter-controls {
        display: flex;
        gap: 1.5rem;
        align-items: flex-end;
        flex-wrap: wrap;
    }

    .date-range {
        display: flex;
        gap: 1rem;
        align-items: flex-end;
        flex-wrap: wrap;
    }

    .date-input {
        display: flex;
        flex-direction: column;
        gap: 0.5rem;
    }

    .date-input label {
        font-weight: 500;
        color: #333;
        font-size: 0.9rem;
    }

    .date-input input {
        padding: 0.5rem;
        border: 1px solid #ccc;
        border-radius: 4px;
        font-size: 1rem;
        min-width: 150px;
    }

    .date-input input:focus {
        outline: none;
        border-color: #0066cc;
        box-shadow: 0 0 0 3px rgba(0, 102, 204, 0.1);
    }

    .date-input input:disabled {
        background-color: #f5f5f5;
        cursor: not-allowed;
    }

    .separator {
        color: #666;
        font-weight: 500;
        padding: 0 0.5rem;
    }

    .btn {
        padding: 0.5rem 1rem;
        border: none;
        border-radius: 4px;
        font-weight: 500;
        cursor: pointer;
        transition: background-color 0.2s;
    }

    .btn-primary {
        background-color: #0066cc;
        color: white;
    }

    .btn-primary:hover:not(:disabled) {
        background-color: #0052a3;
    }

    .btn-primary:disabled {
        opacity: 0.6;
        cursor: not-allowed;
    }

    .loading,
    .empty-notice {
        text-align: center;
        padding: 2rem;
        color: #eee;
        font-size: 1rem;
    }

    .summary-container {
        background-color: white;
        border-radius: 8px;
        border: 1px solid #ddd;
        padding: 1.5rem;
    }

    .period-info {
        font-size: 0.95rem;
        color: #666;
        margin-bottom: 1.5rem;
        font-weight: 500;
    }

    .table-wrapper {
        overflow-x: auto;
        margin-bottom: 2rem;
    }

    .summary-table {
        width: 100%;
        border-collapse: collapse;
        font-size: 0.9rem;
    }

    .summary-table th {
        background-color: #f5f5f5;
        padding: 1rem;
        text-align: right;
        font-weight: 600;
        border-bottom: 2px solid #ddd;
        color: #666;
    }

    .summary-table th:first-child {
        text-align: left;
    }

    .summary-table td {
        padding: 1rem;
        text-align: right;
        border-bottom: 1px solid #eee;
    }

    .summary-table tbody tr:last-child td {
        border-bottom: none;
    }

    .summary-table tbody tr:hover {
        background-color: #f9f9f9;
    }

    .user-name {
        text-align: left;
        font-weight: 500;
        color: #2c3e50;
    }

    .amount-cell {
        color: #2c3e50;
    }

    .amount-cell.unpaid {
        color: #d32f2f;
        font-weight: 500;
    }

    .amount-cell.unclaimed {
        color: #ff9800;
        font-weight: 500;
    }

    .balance-cell {
        font-size: 1rem;
    }

    .balance-positive {
        color: #d32f2f;
        background-color: #ffebee;
    }

    .balance-negative {
        color: #2e7d32;
        background-color: #e8f5e9;
    }

    .balance-zero {
        color: #666;
    }

    .legend {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
        gap: 1.5rem;
        margin-top: 2rem;
        padding-top: 1.5rem;
        border-top: 1px solid #eee;
    }

    .legend-item {
        padding: 1rem;
        border-radius: 4px;
    }

    .legend-label {
        display: inline-block;
        padding: 0.35rem 0.75rem;
        border-radius: 12px;
        font-size: 0.85rem;
        font-weight: 600;
        margin-bottom: 0.5rem;
    }

    .legend-positive {
        background-color: #ffebee;
        color: #d32f2f;
    }

    .legend-negative {
        background-color: #e8f5e9;
        color: #2e7d32;
    }

    .legend-item p {
        margin: 0.5rem 0 0 0;
        color: #666;
        font-size: 0.9rem;
    }

    @media (max-width: 768px) {
        .filter-controls {
            flex-direction: column;
            align-items: stretch;
        }

        .date-range {
            flex-direction: column;
            align-items: stretch;
        }

        .date-input input {
            min-width: unset;
            width: 100%;
        }

        .separator {
            text-align: center;
            padding: 0.25rem 0;
        }

        .summary-table {
            font-size: 0.8rem;
        }

        .summary-table th,
        .summary-table td {
            padding: 0.65rem;
        }

        .legend {
            grid-template-columns: 1fr;
        }
    }
</style>
