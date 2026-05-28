<script lang="ts">
    import { onMount } from "svelte";
    import { apiCallWithAuth } from "$lib/api";
    import type { RestockPayment, AdminUser } from "$lib/types";

    let restockPayments: RestockPayment[] = $state([]);
    let users: AdminUser[] = $state([]);
    let isLoading = $state(false);
    let errorMessage = $state("");

    // モーダル関連
    let showModal = $state(false);
    let modalMode: "add" | "edit" = "add";
    let modalUserId = $state<number | null>(null);
    let modalAmount = $state("");
    let modalDate = $state("");
    let modalNote = $state("");
    let editingPaymentId: number | null = null;

    // 削除確認モーダル
    let showDeleteConfirm = $state(false);
    let deleteTargetId: number | null = null;
    let deleteTargetAmount = 0;

    onMount(async () => {
        await loadUsers();
        await loadRestockPayments();
    });

    async function loadUsers() {
        try {
            const data = await apiCallWithAuth<AdminUser[]>("/users");
            users = data || [];
        } catch (err) {
            console.error("Failed to load users:", err);
        }
    }

    async function loadRestockPayments() {
        isLoading = true;
        errorMessage = "";
        try {
            const data =
                await apiCallWithAuth<RestockPayment[]>("/restock-payments");
            restockPayments = data || [];
        } catch (err) {
            errorMessage =
                err instanceof Error
                    ? err.message
                    : "立替精算記録の取得に失敗しました";
        } finally {
            isLoading = false;
        }
    }

    function getUserName(userId: number): string {
        const user = users.find((u) => u.id === userId);
        return user ? user.name : `ユーザー#${userId}`;
    }

    function openAddModal() {
        modalMode = "add";
        modalUserId = null;
        modalAmount = "";
        modalDate = new Date().toISOString().split("T")[0];
        modalNote = "";
        editingPaymentId = null;
        showModal = true;
    }

    function openEditModal(payment: RestockPayment) {
        modalMode = "edit";
        modalUserId = payment.user_id;
        modalAmount = payment.amount.toString();
        modalDate = new Date(payment.settled_at).toISOString().split("T")[0];
        modalNote = payment.note || "";
        editingPaymentId = payment.id;
        showModal = true;
    }

    function closeModal() {
        showModal = false;
    }

    async function handleModalSubmit() {
        if (!modalUserId) {
            errorMessage = "ユーザーを選択してください";
            return;
        }
        if (!modalAmount || modalAmount === "") {
            errorMessage = "金額を入力してください";
            return;
        }
        if (!modalDate) {
            errorMessage = "精算日を入力してください";
            return;
        }

        errorMessage = "";
        isLoading = true;

        try {
            const settledAt = new Date(modalDate).toISOString();
            if (modalMode === "add") {
                await apiCallWithAuth("/restock-payments", {
                    method: "POST",
                    body: JSON.stringify({
                        user_id: modalUserId,
                        amount: parseInt(modalAmount),
                        settled_at: settledAt,
                        note: modalNote || undefined,
                    }),
                });
            } else {
                // edit
                await apiCallWithAuth(`/restock-payments/${editingPaymentId}`, {
                    method: "PUT",
                    body: JSON.stringify({
                        amount: parseInt(modalAmount),
                        settled_at: settledAt,
                        note: modalNote || undefined,
                    }),
                });
            }
            closeModal();
            await loadRestockPayments();
        } catch (err) {
            errorMessage =
                err instanceof Error ? err.message : "操作に失敗しました";
        } finally {
            isLoading = false;
        }
    }

    function openDeleteConfirm(payment: RestockPayment) {
        deleteTargetId = payment.id;
        deleteTargetAmount = payment.amount;
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
            await apiCallWithAuth(`/restock-payments/${deleteTargetId}`, {
                method: "DELETE",
            });
            closeDeleteConfirm();
            await loadRestockPayments();
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

<div class="restock-payments-page">
    <div class="page-header">
        <h1>立替精算管理</h1>
        <button
            class="btn btn-primary"
            onclick={openAddModal}
            disabled={isLoading}
        >
            + 記録を追加
        </button>
    </div>

    {#if errorMessage}
        <div class="alert alert-error">
            <strong>エラー:</strong>
            {errorMessage}
        </div>
    {/if}

    {#if isLoading && restockPayments.length === 0}
        <div class="loading">読み込み中...</div>
    {:else if restockPayments.length === 0}
        <div class="empty-notice">立替精算記録がありません</div>
    {:else}
        <div class="table-wrapper">
            <table class="restock-payments-table">
                <thead>
                    <tr>
                        <th>精算ID</th>
                        <th>仕入れ者</th>
                        <th>金額</th>
                        <th>精算日</th>
                        <th>メモ</th>
                        <th>操作</th>
                    </tr>
                </thead>
                <tbody>
                    {#each restockPayments as payment (payment.id)}
                        <tr>
                            <td>#{payment.id}</td>
                            <td>{getUserName(payment.user_id)}</td>
                            <td class="amount-cell"
                                >{formatCurrency(payment.amount)}</td
                            >
                            <td>{formatDate(payment.settled_at)}</td>
                            <td class="note-cell">{payment.note || "—"}</td>
                            <td class="action-cell">
                                <button
                                    class="btn-small btn-edit"
                                    onclick={() => openEditModal(payment)}
                                    disabled={isLoading}
                                >
                                    編集
                                </button>
                                <button
                                    class="btn-small btn-delete"
                                    onclick={() => openDeleteConfirm(payment)}
                                    disabled={isLoading}
                                >
                                    削除
                                </button>
                            </td>
                        </tr>
                    {/each}
                </tbody>
            </table>
        </div>
    {/if}
</div>

<!-- 追加/編集モーダル -->
{#if showModal}
    <div class="modal-backdrop" onclick={closeModal}></div>
    <div class="modal">
        <div class="modal-header">
            <h2>
                {modalMode === "add"
                    ? "立替精算記録を追加"
                    : "立替精算記録を編集"}
            </h2>
            <button class="btn-close" onclick={closeModal}>✕</button>
        </div>

        <form
            onsubmit={(e) => {
                e.preventDefault();
                handleModalSubmit();
            }}
        >
            <div class="form-group">
                <label for="payment-user">仕入れ者</label>
                <select
                    id="payment-user"
                    bind:value={modalUserId}
                    required
                    disabled={isLoading || modalMode === "edit"}
                >
                    <option value={null}>ユーザーを選択...</option>
                    {#each users as user (user.id)}
                        <option value={user.id}>{user.name}</option>
                    {/each}
                </select>
            </div>

            <div class="form-group">
                <label for="payment-amount">金額（円）</label>
                <input
                    type="number"
                    id="payment-amount"
                    bind:value={modalAmount}
                    placeholder="10000"
                    required
                    min="1"
                    disabled={isLoading}
                />
            </div>

            <div class="form-group">
                <label for="payment-date">精算日</label>
                <input
                    type="date"
                    id="payment-date"
                    bind:value={modalDate}
                    required
                    disabled={isLoading}
                />
            </div>

            <div class="form-group">
                <label for="payment-note">メモ</label>
                <textarea
                    id="payment-note"
                    bind:value={modalNote}
                    placeholder="メモを入力（任意）"
                    disabled={isLoading}
                    rows="3"
                ></textarea>
            </div>

            <div class="modal-footer">
                <button
                    type="button"
                    class="btn btn-secondary"
                    onclick={closeModal}
                    disabled={isLoading}
                >
                    キャンセル
                </button>
                <button
                    type="submit"
                    class="btn btn-primary"
                    disabled={isLoading}
                >
                    {isLoading ? "処理中..." : "保存"}
                </button>
            </div>
        </form>
    </div>
{/if}

<!-- 削除確認モーダル -->
{#if showDeleteConfirm}
    <div class="modal-backdrop" onclick={closeDeleteConfirm}></div>
    <div class="modal modal-small">
        <div class="modal-header">
            <h2>立替精算記録を削除</h2>
        </div>

        <div class="modal-body">
            <p>
                この立替精算記録（<strong
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
    .restock-payments-page {
        max-width: 1200px;
    }

    .page-header {
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
    }

    .restock-payments-table {
        width: 100%;
        border-collapse: collapse;
        font-size: 0.95rem;
    }

    .restock-payments-table th {
        background-color: #f5f5f5;
        padding: 1rem;
        text-align: left;
        font-weight: 600;
        border-bottom: 2px solid #ddd;
        color: #2c3e50;
    }

    .restock-payments-table td {
        padding: 1rem;
        border-bottom: 1px solid #eee;
    }

    .restock-payments-table tbody tr:last-child td {
        border-bottom: none;
    }

    .restock-payments-table tbody tr:hover {
        background-color: #f9f9f9;
    }

    .amount-cell {
        font-weight: 600;
        color: #0066cc;
    }

    .note-cell {
        color: #666;
        font-size: 0.9rem;
        max-width: 150px;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
    }

    .action-cell {
        display: flex;
        gap: 0.5rem;
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

    .btn-small {
        padding: 0.35rem 0.75rem;
        font-size: 0.85rem;
        border: none;
        border-radius: 4px;
        cursor: pointer;
        transition: background-color 0.2s;
    }

    .btn-edit {
        background-color: #0066cc;
        color: white;
    }

    .btn-edit:hover:not(:disabled) {
        background-color: #0052a3;
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
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 1.5rem;
        border-bottom: 1px solid #eee;
    }

    .modal-header h2 {
        margin: 0;
        font-size: 1.3rem;
        color: #2c3e50;
    }

    .btn-close {
        background: none;
        border: none;
        font-size: 1.5rem;
        color: #999;
        cursor: pointer;
        padding: 0;
    }

    .btn-close:hover {
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

    form {
        padding: 1.5rem;
    }

    .form-group {
        margin-bottom: 1.5rem;
    }

    label {
        display: block;
        margin-bottom: 0.5rem;
        font-weight: 500;
        color: #2c3e50;
    }

    input,
    select,
    textarea {
        width: 100%;
        padding: 0.75rem;
        border: 1px solid #ccc;
        border-radius: 4px;
        font-size: 1rem;
        font-family: inherit;
        box-sizing: border-box;
    }

    input:focus,
    select:focus,
    textarea:focus {
        outline: none;
        border-color: #0066cc;
        box-shadow: 0 0 0 3px rgba(0, 102, 204, 0.1);
    }

    input:disabled,
    select:disabled,
    textarea:disabled {
        background-color: #f5f5f5;
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
        .page-header {
            flex-direction: column;
            gap: 1rem;
            align-items: flex-start;
        }

        .modal,
        .modal-small {
            min-width: 90vw;
            max-width: 90vw;
        }

        .restock-payments-table {
            font-size: 0.85rem;
        }

        .restock-payments-table th,
        .restock-payments-table td {
            padding: 0.75rem;
        }

        .action-cell {
            flex-direction: column;
        }

        .btn-small {
            width: 100%;
        }
    }
</style>
