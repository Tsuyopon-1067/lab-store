<script lang="ts">
	import { onMount } from 'svelte';
	import { apiCallWithAuth } from '$lib/api';
	import type { AdminPurchase, AdminUser } from '$lib/types';

	let purchases: AdminPurchase[] = $state([]);
	let users: AdminUser[] = $state([]);
	let isLoading = $state(false);
	let errorMessage = $state('');

	// フィルタリング
	let selectedUserId = $state<number | null>(null);
	let currentPage = $state(1);
	const itemsPerPage = 50;

	onMount(async () => {
		await loadUsers();
		await loadPurchases();
	});

	async function loadUsers() {
		try {
			const data = await apiCallWithAuth<AdminUser[]>('/users');
			users = data || [];
		} catch (err) {
			console.error('Failed to load users:', err);
		}
	}

	async function loadPurchases() {
		isLoading = true;
		errorMessage = '';
		try {
			const offset = (currentPage - 1) * itemsPerPage;
			const params = new URLSearchParams({
				limit: itemsPerPage.toString(),
				offset: offset.toString()
			});

			if (selectedUserId !== null && selectedUserId !== -1) {
				params.append('user_id', selectedUserId.toString());
			}

			const data = await apiCallWithAuth<AdminPurchase[]>(`/purchases?${params}`);
			purchases = data || [];
		} catch (err) {
			errorMessage = err instanceof Error ? err.message : '購入履歴の取得に失敗しました';
		} finally {
			isLoading = false;
		}
	}

	function getUserName(userId: number): string {
		const user = users.find((u) => u.id === userId);
		return user ? user.name : `ユーザー#${userId}`;
	}

	function handleFilterChange() {
		currentPage = 1;
		loadPurchases();
	}

	function formatDate(dateString: string): string {
		return new Date(dateString).toLocaleString('ja-JP');
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
</script>

<div class="purchases-page">
	<h1>購入履歴</h1>

	{#if errorMessage}
		<div class="alert alert-error">
			<strong>エラー:</strong> {errorMessage}
		</div>
	{/if}

	<div class="filter-section">
		<div class="filter-group">
			<label for="user-filter">ユーザーで絞り込む</label>
			<select id="user-filter" bind:value={selectedUserId} onchange={handleFilterChange} disabled={isLoading}>
				<option value={null}>すべてのユーザー</option>
				{#each users as user (user.id)}
					<option value={user.id}>{user.name}</option>
				{/each}
			</select>
		</div>
	</div>

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
					</tr>
				</thead>
				<tbody>
					{#each purchases as purchase (purchase.id)}
						<tr>
							<td>#{purchase.id}</td>
							<td>{getUserName(purchase.user_id)}</td>
							<td>{formatDate(purchase.purchased_at)}</td>
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

	.filter-group {
		display: flex;
		gap: 1rem;
		align-items: center;
	}

	label {
		font-weight: 500;
		color: #2c3e50;
	}

	select {
		padding: 0.5rem;
		border: 1px solid #ccc;
		border-radius: 4px;
		font-size: 1rem;
		min-width: 200px;
	}

	select:focus {
		outline: none;
		border-color: #0066cc;
		box-shadow: 0 0 0 3px rgba(0, 102, 204, 0.1);
	}

	select:disabled {
		background-color: #f5f5f5;
		cursor: not-allowed;
	}

	.loading,
	.empty-notice {
		text-align: center;
		padding: 2rem;
		color: #2c3e50;
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

	.purchases-table tbody tr:last-child td {
		border-bottom: none;
	}

	.purchases-table tbody tr:hover {
		background-color: #f9f9f9;
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

	.page-info {
		font-weight: 500;
		color: #2c3e50;
	}

	@media (max-width: 768px) {
		.filter-group {
			flex-direction: column;
			align-items: flex-start;
		}

		select {
			width: 100%;
			min-width: unset;
		}

		.purchases-table {
			font-size: 0.85rem;
		}

		.purchases-table th,
		.purchases-table td {
			padding: 0.75rem;
		}

		.pagination {
			flex-wrap: wrap;
		}
	}
</style>
