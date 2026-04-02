<script lang="ts">
    import { onMount } from "svelte";
    import { apiCallWithAuth } from "$lib/api";
    import type { APIKey, CreateAPIKeyResponse } from "$lib/types";

    let apiKeys: APIKey[] = $state([]);
    let isLoading = $state(false);
    let errorMessage = $state("");

    // Create modal state
    let showCreateModal = $state(false);
    let keyName = $state("");
    let isCreating = $state(false);

    // Key display modal (show raw key after creation)
    let showKeyModal = $state(false);
    let newRawKey = $state("");
    let newKeyName = $state("");

    // Delete confirm modal
    let showDeleteConfirm = $state(false);
    let deleteTargetId: number | null = null;
    let deleteTargetName = $state("");

    onMount(async () => {
        await loadKeys();
    });

    async function loadKeys() {
        isLoading = true;
        errorMessage = "";
        try {
            const response = await apiCallWithAuth<{ keys: APIKey[] }>(
                "/api-keys",
            );
            apiKeys = response.keys || [];
        } catch (err) {
            errorMessage =
                err instanceof Error
                    ? err.message
                    : "キー一覧の取得に失敗しました";
        } finally {
            isLoading = false;
        }
    }

    function openCreateModal() {
        showCreateModal = true;
        keyName = "";
    }

    function closeCreateModal() {
        showCreateModal = false;
        keyName = "";
        isCreating = false;
    }

    async function handleCreateSubmit() {
        if (!keyName.trim()) {
            errorMessage = "キー名を入力してください";
            return;
        }

        isCreating = true;
        errorMessage = "";
        try {
            const response = await apiCallWithAuth<CreateAPIKeyResponse>(
                "/api-keys",
                {
                    method: "POST",
                    body: JSON.stringify({ name: keyName }),
                },
            );

            // Show raw key modal
            newRawKey = response.key;
            newKeyName = response.name;
            showKeyModal = true;
            closeCreateModal();

            // Reload keys
            await loadKeys();
        } catch (err) {
            errorMessage =
                err instanceof Error ? err.message : "キーの発行に失敗しました";
        } finally {
            isCreating = false;
        }
    }

    function closeKeyModal() {
        showKeyModal = false;
        newRawKey = "";
        newKeyName = "";
    }

    function copyToClipboard() {
        navigator.clipboard.writeText(newRawKey);
        alert("APIキーをコピーしました");
    }

    function openDeleteConfirm(key: APIKey) {
        deleteTargetId = key.id;
        deleteTargetName = key.name;
        showDeleteConfirm = true;
    }

    function closeDeleteConfirm() {
        showDeleteConfirm = false;
        deleteTargetId = null;
        deleteTargetName = "";
    }

    async function handleDelete() {
        if (!deleteTargetId) return;

        isLoading = true;
        errorMessage = "";
        try {
            await apiCallWithAuth(`/api-keys/${deleteTargetId}`, {
                method: "DELETE",
            });
            await loadKeys();
        } catch (err) {
            errorMessage =
                err instanceof Error ? err.message : "キーの削除に失敗しました";
        } finally {
            isLoading = false;
            closeDeleteConfirm();
        }
    }

    function formatDate(dateStr: string | null) {
        if (!dateStr) return "-";
        return new Date(dateStr).toLocaleDateString("ja-JP", {
            year: "numeric",
            month: "2-digit",
            day: "2-digit",
            hour: "2-digit",
            minute: "2-digit",
        });
    }
</script>

<div class="api-keys-page">
    <div class="page-header">
        <h1>APIキー管理</h1>
        <button class="btn btn-primary" onclick={openCreateModal}
            >+ 新規発行</button
        >
    </div>

    {#if errorMessage}
        <div class="alert alert-error">{errorMessage}</div>
    {/if}

    {#if isLoading && apiKeys.length === 0}
        <div class="loading">読み込み中...</div>
    {:else if apiKeys.length === 0}
        <div class="empty-notice">発行済みのAPIキーはありません</div>
    {:else}
        <div class="table-wrapper">
            <table>
                <thead>
                    <tr>
                        <th>キー名</th>
                        <th>プレフィックス</th>
                        <th>最終使用日</th>
                        <th>作成日</th>
                        <th>アクション</th>
                    </tr>
                </thead>
                <tbody>
                    {#each apiKeys as key (key.id)}
                        <tr>
                            <td>{key.name}</td>
                            <td class="barcode-cell">{key.key_prefix}***</td>
                            <td>{formatDate(key.last_used_at)}</td>
                            <td>{formatDate(key.created_at)}</td>
                            <td class="action-cell">
                                <button
                                    class="btn btn-small btn-delete"
                                    onclick={() => openDeleteConfirm(key)}
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

<!-- Create API Key Modal -->
{#if showCreateModal}
    <div class="modal-backdrop" onclick={closeCreateModal}></div>
    <div class="modal">
        <div class="modal-header">
            <h2>新しいAPIキーを発行</h2>
        </div>
        <form
            onsubmit={(e) => {
                e.preventDefault();
                handleCreateSubmit();
            }}
        >
            <div class="form-group">
                <label for="key-name">キー名</label>
                <input
                    type="text"
                    id="key-name"
                    bind:value={keyName}
                    placeholder="例: 社外連携API"
                    disabled={isCreating}
                />
            </div>
            <div class="modal-footer">
                <button
                    type="button"
                    class="btn btn-secondary"
                    onclick={closeCreateModal}
                    disabled={isCreating}
                >
                    キャンセル
                </button>
                <button
                    type="submit"
                    class="btn btn-primary"
                    disabled={isCreating || !keyName.trim()}
                >
                    {isCreating ? "発行中..." : "発行"}
                </button>
            </div>
        </form>
    </div>
{/if}

<!-- Raw API Key Display Modal -->
{#if showKeyModal}
    <div class="modal-backdrop" onclick={closeKeyModal}></div>
    <div class="modal modal-large">
        <div class="modal-header">
            <h2>⚠️ APIキーが発行されました</h2>
        </div>
        <div class="modal-body">
            <div class="warning-box">
                <p>
                    <strong>重要:</strong>
                    このキーは今後二度と表示されません。
                    <strong>必ずコピーして安全に保管してください。</strong>
                </p>
            </div>

            <div class="key-display-section">
                <label>キー名: {newKeyName}</label>
                <div class="key-display">
                    <code>{newRawKey}</code>
                    <button
                        type="button"
                        class="btn btn-small btn-primary"
                        onclick={copyToClipboard}
                    >
                        コピー
                    </button>
                </div>
            </div>
        </div>
        <div class="modal-footer">
            <button
                type="button"
                class="btn btn-primary"
                onclick={closeKeyModal}
            >
                確認しました
            </button>
        </div>
    </div>
{/if}

<!-- Delete Confirm Modal -->
{#if showDeleteConfirm && deleteTargetId !== null}
    <div class="modal-backdrop" onclick={closeDeleteConfirm}></div>
    <div class="modal modal-small">
        <div class="modal-header">
            <h2>確認</h2>
        </div>
        <div class="modal-body">
            <p>APIキー「{deleteTargetName}」を削除してもよろしいですか?</p>
            <p style="color: #d32f2f; margin-top: 1rem; font-size: 0.9rem;">
                このキーを使用しているアプリケーションは動作しなくなります。
            </p>
        </div>
        <div class="modal-footer">
            <button
                type="button"
                class="btn btn-secondary"
                onclick={closeDeleteConfirm}
                disabled={isLoading}
            >
                キャンセル
            </button>
            <button
                type="button"
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
    .api-keys-page {
        max-width: 1200px;
    }

    h1 {
        font-size: 2rem;
        margin-bottom: 0;
        color: #fff;
    }

    h2 {
        font-size: 1.3rem;
        margin: 0;
        color: #333;
    }

    .page-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 2rem;
    }

    .alert {
        padding: 1rem;
        border-radius: 4px;
        margin-bottom: 1.5rem;
    }

    .alert-error {
        background-color: #ffebee;
        color: #c62828;
        border: 1px solid #ef5350;
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
        border: 1px solid #ddd;
        border-radius: 8px;
        overflow: hidden;
    }

    table {
        width: 100%;
        border-collapse: collapse;
    }

    thead {
        background-color: #f5f5f5;
        border-bottom: 1px solid #ddd;
    }

    th {
        padding: 1rem;
        text-align: left;
        font-weight: 600;
        color: #333;
        font-size: 0.95rem;
    }

    td {
        padding: 1rem;
        border-bottom: 1px solid #eee;
        color: #666;
    }

    tbody tr:hover {
        background-color: #f9f9f9;
    }

    .barcode-cell {
        font-family: monospace;
        color: #666;
    }

    .action-cell {
        display: flex;
        gap: 0.5rem;
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
        border: none;
        font-size: 0.95rem;
    }

    .btn:disabled {
        opacity: 0.6;
        cursor: not-allowed;
    }

    .btn-primary {
        background-color: #0066cc;
        color: white;
    }

    .btn-primary:hover:not(:disabled) {
        background-color: #0052a3;
    }

    .btn-secondary {
        background-color: #999;
        color: white;
    }

    .btn-secondary:hover:not(:disabled) {
        background-color: #777;
    }

    .btn-danger {
        background-color: #d32f2f;
        color: white;
    }

    .btn-danger:hover:not(:disabled) {
        background-color: #b71c1c;
    }

    .btn-small {
        padding: 0.3rem 0.7rem;
        font-size: 0.85rem;
    }

    .btn-delete {
        background-color: #d32f2f;
        color: white;
    }

    .btn-delete:hover:not(:disabled) {
        background-color: #b71c1c;
    }

    .form-group {
        margin-bottom: 1.5rem;
    }

    .form-group label {
        display: block;
        margin-bottom: 0.5rem;
        font-weight: 500;
        color: #333;
    }

    .form-group input,
    .form-group textarea,
    .form-group select {
        width: 100%;
        padding: 0.75rem;
        border: 1px solid #ddd;
        border-radius: 4px;
        font-size: 1rem;
        font-family: inherit;
    }

    .form-group input:focus,
    .form-group textarea:focus,
    .form-group select:focus {
        outline: none;
        border-color: #0066cc;
        box-shadow: 0 0 0 3px rgba(0, 102, 204, 0.1);
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
        min-width: 500px;
        max-height: 90vh;
        overflow-y: auto;
    }

    .modal-small {
        min-width: 400px;
    }

    .modal-large {
        min-width: 600px;
    }

    .modal-header {
        padding: 1.5rem;
        border-bottom: 1px solid #ddd;
        background-color: #f5f5f5;
    }

    .modal-body {
        padding: 1.5rem;
    }

    .modal-footer {
        padding: 1.5rem;
        border-top: 1px solid #ddd;
        display: flex;
        gap: 1rem;
        justify-content: flex-end;
    }

    .warning-box {
        background-color: #fff3cd;
        border: 1px solid #ffc107;
        border-radius: 4px;
        padding: 1rem;
        margin-bottom: 1.5rem;
        color: #856404;
    }

    .warning-box p {
        margin: 0;
    }

    .key-display-section {
        background-color: #f5f5f5;
        border: 1px solid #ddd;
        border-radius: 4px;
        padding: 1rem;
    }

    .key-display-section label {
        display: block;
        margin-bottom: 0.5rem;
        font-weight: 500;
        color: #333;
    }

    .key-display {
        display: flex;
        gap: 1rem;
        align-items: flex-start;
    }

    .key-display code {
        flex: 1;
        background-color: white;
        border: 1px solid #ddd;
        border-radius: 4px;
        padding: 0.75rem;
        font-family: "Courier New", monospace;
        word-break: break-all;
        color: #333;
    }
</style>
