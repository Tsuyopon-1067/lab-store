<script lang="ts">
    import { onMount } from "svelte";
    import { apiCallWithAuth } from "$lib/api";
    import BarcodeInput from "$lib/components/BarcodeInput.svelte";
    import type { AdminUser } from "$lib/types";

    let users: AdminUser[] = $state([]);
    let isLoading = $state(false);
    let errorMessage = $state("");

    // モーダル関連
    let showModal = $state(false);
    let modalMode: "add" | "edit" = "add";
    let modalName = $state("");
    let modalBarcode = $state("");
    let editingUserId: number | null = null;

    // 削除確認モーダル
    let showDeleteConfirm = $state(false);
    let deleteTargetId: number | null = null;
    let deleteTargetName = "";

    onMount(async () => {
        await loadUsers();
    });

    async function loadUsers() {
        isLoading = true;
        errorMessage = "";
        try {
            const data = await apiCallWithAuth<AdminUser[]>("/users");
            users = data || [];
        } catch (err) {
            errorMessage =
                err instanceof Error
                    ? err.message
                    : "利用者一覧の取得に失敗しました";
        } finally {
            isLoading = false;
        }
    }

    function openAddModal() {
        modalMode = "add";
        modalName = "";
        modalBarcode = "";
        editingUserId = null;
        showModal = true;
    }

    function openEditModal(user: AdminUser) {
        modalMode = "edit";
        modalName = user.name;
        modalBarcode = "";
        editingUserId = user.id;
        showModal = true;
    }

    function closeModal() {
        showModal = false;
    }

    async function handleModalSubmit() {
        if (!modalName.trim()) {
            errorMessage = "氏名を入力してください";
            return;
        }

        errorMessage = "";
        isLoading = true;

        try {
            if (modalMode === "add") {
                if (!modalBarcode.trim()) {
                    errorMessage = "バーコードを入力してください";
                    return;
                }
                await apiCallWithAuth("/users", {
                    method: "POST",
                    body: JSON.stringify({
                        name: modalName,
                        barcode: modalBarcode,
                    }),
                });
            } else {
                // edit
                await apiCallWithAuth(`/users/${editingUserId}`, {
                    method: "PUT",
                    body: JSON.stringify({ name: modalName }),
                });
            }
            closeModal();
            await loadUsers();
        } catch (err) {
            errorMessage =
                err instanceof Error ? err.message : "操作に失敗しました";
        } finally {
            isLoading = false;
        }
    }

    function openDeleteConfirm(user: AdminUser) {
        deleteTargetId = user.id;
        deleteTargetName = user.name;
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
            await apiCallWithAuth(`/users/${deleteTargetId}`, {
                method: "DELETE",
            });
            closeDeleteConfirm();
            await loadUsers();
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

    function handleScan(event: CustomEvent<{ barcode: string }>) {
        if (showModal && modalMode === "add") {
            modalBarcode = event.detail.barcode;
        }
    }
</script>

<svelte:window on:scan={handleScan} />

<div class="users-page">
    <div class="page-header">
        <h1>利用者管理</h1>
        <button
            class="btn btn-primary"
            onclick={openAddModal}
            disabled={isLoading}
        >
            + 新規追加
        </button>
    </div>

    {#if errorMessage}
        <div class="alert alert-error">
            <strong>エラー:</strong>
            {errorMessage}
        </div>
    {/if}

    {#if isLoading && users.length === 0}
        <div class="loading">読み込み中...</div>
    {:else if users.length === 0}
        <div class="empty-notice">利用者がいません</div>
    {:else}
        <div class="table-wrapper">
            <table class="users-table">
                <thead>
                    <tr>
                        <th>ID</th>
                        <th>氏名</th>
                        <th>バーコード</th>
                        <th>状態</th>
                        <th>登録日</th>
                        <th>操作</th>
                    </tr>
                </thead>
                <tbody>
                    {#each users as user (user.id)}
                        <tr class:inactive={user.is_active === 0}>
                            <td>{user.id}</td>
                            <td>{user.name}</td>
                            <td class="barcode-cell">{user.barcode}</td>
                            <td>
                                <span
                                    class="status-badge"
                                    class:active={user.is_active === 1}
                                >
                                    {user.is_active === 1 ? "有効" : "無効"}
                                </span>
                            </td>
                            <td>{formatDate(user.created_at)}</td>
                            <td class="action-cell">
                                <button
                                    class="btn-small btn-edit"
                                    onclick={() => openEditModal(user)}
                                    disabled={isLoading}
                                >
                                    編集
                                </button>
                                <button
                                    class="btn-small btn-delete"
                                    onclick={() => openDeleteConfirm(user)}
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
            <h2>{modalMode === "add" ? "利用者を追加" : "利用者を編集"}</h2>
            <button class="btn-close" onclick={closeModal}>✕</button>
        </div>

        <form
            onsubmit={(e) => {
                e.preventDefault();
                handleModalSubmit();
            }}
        >
            <div class="form-group">
                <label for="name">氏名</label>
                <input
                    type="text"
                    id="name"
                    bind:value={modalName}
                    placeholder="山田 太郎"
                    required
                    disabled={isLoading}
                />
            </div>

            {#if modalMode === "add"}
                <div class="form-group">
                    <label>バーコード</label>
                    <BarcodeInput bind:input={modalBarcode} />
                </div>
            {/if}

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
            <h2>利用者を削除</h2>
        </div>

        <div class="modal-body">
            <p>
                <strong>{deleteTargetName}</strong> を削除してもよろしいですか？
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
    .users-page {
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
        color: #fff;
        font-size: 1rem;
    }

    .table-wrapper {
        background-color: white;
        border-radius: 8px;
        border: 1px solid #ddd;
        overflow: hidden;
    }

    .users-table {
        width: 100%;
        border-collapse: collapse;
        font-size: 0.95rem;
    }

    .users-table th {
        background-color: #f5f5f5;
        padding: 1rem;
        text-align: left;
        font-weight: 600;
        border-bottom: 2px solid #ddd;
        color: #2c3e50;
    }

    .users-table td {
        padding: 1rem;
        border-bottom: 1px solid #eee;
    }

    .users-table tbody tr:last-child td {
        border-bottom: none;
    }

    .users-table tbody tr.inactive {
        background-color: #f9f9f9;
        opacity: 0.7;
    }

    .barcode-cell {
        font-family: monospace;
        color: #666;
    }

    .status-badge {
        display: inline-block;
        padding: 0.25rem 0.75rem;
        border-radius: 12px;
        font-size: 0.85rem;
        font-weight: 500;
        background-color: #fee;
        color: #c00;
    }

    .status-badge.active {
        background-color: #e8f5e9;
        color: #2e7d32;
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
        color: #333;
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
        color: #333;
    }

    .modal-body {
        padding: 1.5rem;
    }

    .modal-body p {
        margin: 0.5rem 0;
        color: #333;
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
        color: #333;
    }

    input {
        width: 100%;
        padding: 0.75rem;
        border: 1px solid #ccc;
        border-radius: 4px;
        font-size: 1rem;
        box-sizing: border-box;
    }

    input:focus {
        outline: none;
        border-color: #0066cc;
        box-shadow: 0 0 0 3px rgba(0, 102, 204, 0.1);
    }

    input:disabled {
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

        .modal {
            min-width: 90vw;
            max-width: 90vw;
        }

        .users-table {
            font-size: 0.85rem;
        }

        .users-table th,
        .users-table td {
            padding: 0.75rem;
        }

        .action-cell {
            flex-direction: column;
        }
    }
</style>
