<script lang="ts">
    import { onMount } from "svelte";
    import { apiCallWithAuth } from "$lib/api";
    import type { AdminRestock, AdminUser } from "$lib/types";

    let restocks: AdminRestock[] = $state([]);
    let users: AdminUser[] = $state([]);
    let isLoading = $state(false);
    let errorMessage = $state("");

    // 削除確認モーダル
    let showDeleteConfirm = $state(false);
    let deleteTargetId: number | null = null;
    let deleteTargetAmount = 0;

    onMount(async () => {
        await loadUsers();
        await loadRestocks();
    });

    async function loadUsers() {
        try {
            const data = await apiCallWithAuth<AdminUser[]>("/users");
            users = data || [];
        } catch (err) {
            console.error("Failed to load users:", err);
        }
    }

    async function loadRestocks() {
        isLoading = true;
        errorMessage = "";
        try {
            const data = await apiCallWithAuth<AdminRestock[]>("/restocks");
            restocks = data || [];
        } catch (err) {
            errorMessage =
                err instanceof Error
                    ? err.message
                    : "仕入れ履歴の取得に失敗しました";
        } finally {
            isLoading = false;
        }
    }

    function getUserName(userId: number): string {
        const user = users.find((u) => u.id === userId);
        return user ? user.name : `ユーザー#${userId}`;
    }

    function isDeleted(restock: AdminRestock): boolean {
        return !!restock.deleted_at;
    }

    function openDeleteConfirm(restock: AdminRestock) {
        deleteTargetId = restock.id;
        deleteTargetAmount = restock.total_amount;
        showDeleteConfirm = true;
    }

    function closeDeleteConfirm() {
        showDeleteConfirm = false;
        deleteTargetId = null;
    }

    async function handleDelete() {
        if (!deleteTargetId) return;

        errorMessage = "";
        isLoading = true;

        try {
            await apiCallWithAuth(`/restocks/${deleteTargetId}`, {
                method: "DELETE",
            });
            closeDeleteConfirm();
            await loadRestocks();
        } catch (err) {
            errorMessage =
                err instanceof Error ? err.message : "削除に失敗しました";
        } finally {
            isLoading = false;
        }
    }

    function formatDate(dateString: string): string {
        return new Date(dateString).toLocaleString("ja-JP");
    }

    function formatCurrency(amount: number): string {
        return `¥${amount.toLocaleString("ja-JP")}`;
    }
</script>

<div class="restocks-page">
    <h1>仕入れ履歴</h1>

    {#if errorMessage}
        <div class="alert alert-error">
            <strong>エラー:</strong>
            {errorMessage}
        </div>
    {/if}

    {#if isLoading && restocks.length === 0}
        <div class="loading">読み込み中...</div>
    {:else if restocks.length === 0}
        <div class="empty-notice">仕入れ履歴がありません</div>
    {:else}
        <div class="table-wrapper">
            <table class="restocks-table">
                <thead>
                    <tr>
                        <th>仕入れID</th>
                        <th>仕入れ者</th>
                        <th>金額</th>
                        <th>仕入れ日時</th>
                        <th>メモ</th>
                        <th>状態</th>
                        <th>操作</th>
                    </tr>
                </thead>
                <tbody>
                    {#each restocks as restock (restock.id)}
                        <tr class:deleted={isDeleted(restock)}>
                            <td>#{restock.id}</td>
                            <td>{getUserName(restock.user_id)}</td>
                            <td class="amount-cell"
                                >{formatCurrency(restock.total_amount)}</td
                            >
                            <td>{formatDate(restock.restocked_at)}</td>
                            <td class="note-cell">{restock.note || "—"}</td>
                            <td>
                                {#if isDeleted(restock)}
                                    <span class="status-badge status-deleted"
                                        >削除済み</span
                                    >
                                {:else}
                                    <span class="status-badge status-active"
                                        >有効</span
                                    >
                                {/if}
                            </td>
                            <td class="action-cell">
                                {#if !isDeleted(restock)}
                                    <button
                                        class="btn-small btn-delete"
                                        onclick={() =>
                                            openDeleteConfirm(restock)}
                                        disabled={isLoading}
                                    >
                                        削除
                                    </button>
                                {/if}
                            </td>
                        </tr>
                    {/each}
                </tbody>
            </table>
        </div>
    {/if}
</div>

<!-- 削除確認モーダル -->
{#if showDeleteConfirm}
    <div class="modal-backdrop" onclick={closeDeleteConfirm}></div>
    <div class="modal modal-small">
        <div class="modal-header">
            <h2>仕入れ記録を削除</h2>
        </div>

        <div class="modal-body">
            <p>
                この仕入れ記録（<strong
                    >{formatCurrency(deleteTargetAmount)}</strong
                >）を削除してもよろしいですか？
            </p>
            <p class="notice">この操作は取り消せません。</p>
        </div>

        <div class="modal-footer">
            <button
                class="btn btn-secondary"
                onclick={closeDeleteConfirm}
                disabled={isLoading}
            >
                キャンセル
            </button>
            <button
                class="btn btn-danger"
                onclick={handleDelete}
                disabled={isLoading}
            >
                {isLoading ? "削除中..." : "削除"}
            </button>
        </div>
    </div>
{/if}

<style>
    .restocks-page {
        max-width: 1200px;
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
    }

    .restocks-table {
        width: 100%;
        border-collapse: collapse;
        font-size: 0.95rem;
    }

    .restocks-table th {
        background-color: #f5f5f5;
        padding: 1rem;
        text-align: left;
        font-weight: 600;
        border-bottom: 2px solid #ddd;
        color: #2c3e50;
    }

    .restocks-table td {
        padding: 1rem;
        border-bottom: 1px solid #eee;
    }

    .restocks-table tbody tr:last-child td {
        border-bottom: none;
    }

    .restocks-table tbody tr.deleted {
        background-color: #f9f9f9;
        opacity: 0.7;
    }

    .restocks-table tbody tr:not(.deleted):hover {
        background-color: #f9f9f9;
    }

    .amount-cell {
        font-weight: 600;
        color: #0066cc;
    }

    .note-cell {
        color: #666;
        font-size: 0.9rem;
        max-width: 200px;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
    }

    .status-badge {
        display: inline-block;
        padding: 0.25rem 0.75rem;
        border-radius: 12px;
        font-size: 0.85rem;
        font-weight: 500;
    }

    .status-active {
        background-color: #e8f5e9;
        color: #2e7d32;
    }

    .status-deleted {
        background-color: #ffebee;
        color: #c62828;
    }

    .action-cell {
        display: flex;
        gap: 0.5rem;
    }

    .btn-small {
        padding: 0.35rem 0.75rem;
        font-size: 0.85rem;
        border: none;
        border-radius: 4px;
        cursor: pointer;
        transition: background-color 0.2s;
    }

    .btn-delete {
        background-color: #d32f2f;
        color: white;
    }

    .btn-delete:hover:not(:disabled) {
        background-color: #b71c1c;
    }

    .btn-small:disabled {
        opacity: 0.6;
        cursor: not-allowed;
    }

    .modal-backdrop {
        position: fixed;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background-color: rgba(0, 0, 0, 0.5);
        z-index: 100;
    }

    .modal {
        position: fixed;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
        background-color: white;
        border-radius: 8px;
        box-shadow: 0 10px 40px rgba(0, 0, 0, 0.2);
        z-index: 101;
        min-width: 400px;
        max-width: 500px;
        max-height: 90vh;
        overflow-y: auto;
    }

    .modal-small {
        min-width: 350px;
        max-width: 450px;
    }

    .modal-header {
        padding: 1.5rem;
        border-bottom: 1px solid #eee;
    }

    .modal-header h2 {
        margin: 0;
        font-size: 1.3rem;
        color: #2c3e50;
    }

    .modal-body {
        padding: 1.5rem;
    }

    .modal-body p {
        margin: 0.5rem 0;
        color: #2c3e50;
    }

    .notice {
        color: #d32f2f;
        font-size: 0.9rem;
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

    .btn-danger {
        background-color: #d32f2f;
        color: white;
    }

    .btn-danger:hover:not(:disabled) {
        background-color: #b71c1c;
    }

    .btn-danger:disabled {
        opacity: 0.6;
        cursor: not-allowed;
    }

    .modal-footer {
        display: flex;
        gap: 1rem;
        padding: 1.5rem;
        border-top: 1px solid #eee;
        justify-content: flex-end;
    }

    .modal-footer button {
        min-width: 100px;
    }

    @media (max-width: 768px) {
        .restocks-table {
            font-size: 0.85rem;
        }

        .restocks-table th,
        .restocks-table td {
            padding: 0.75rem;
        }

        .modal,
        .modal-small {
            min-width: 90vw;
            max-width: 90vw;
        }
    }
</style>
