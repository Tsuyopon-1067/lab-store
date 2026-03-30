<script lang="ts">
	import { onMount } from 'svelte';
	import { apiCall, apiCallWithBarcode } from '$lib/api';
	import BarcodeInput from '$lib/components/BarcodeInput.svelte';
	import UserInfo from '$lib/components/UserInfo.svelte';
	import type { UserBalance, PurchaseHistory, RestockHistory } from '$lib/types';

	type PageState = 'idle' | 'viewing';
	type TabType = 'purchases' | 'restocks';

	// 状態管理
	let pageState: PageState = $state('idle');
	let currentUser: UserBalance | null = $state(null);
	let purchases: PurchaseHistory[] = $state([]);
	let restocks: RestockHistory[] = $state([]);
	let activeTab: TabType = $state('purchases');
	let errorMessage = $state('');
	let isLoading = $state(false);

	// タイムアウト設定
	const TIMEOUT_MS = 5 * 60 * 1000; // 5分
	let timeoutId: ReturnType<typeof setTimeout> | null = null;

	function resetTimeout() {
		if (timeoutId) clearTimeout(timeoutId);
		timeoutId = setTimeout(() => {
			resetSession();
		}, TIMEOUT_MS);
	}

	function resetSession() {
		pageState = 'idle';
		currentUser = null;
		purchases = [];
		restocks = [];
		activeTab = 'purchases';
		errorMessage = '';
		if (timeoutId) clearTimeout(timeoutId);
	}

	async function handleScan(event: CustomEvent<{ barcode: string }>) {
		resetTimeout();
		const { barcode } = event.detail;

		if (pageState === 'idle') {
			await handleUserScan(barcode);
		}
	}

	async function handleUserScan(barcode: string) {
		isLoading = true;
		errorMessage = '';

		try {
			// ユーザー情報取得
			const balance = await apiCallWithBarcode<UserBalance>('/me/balance', barcode);
			currentUser = balance;

			// 購入履歴取得
			const purchaseList = await apiCallWithBarcode<PurchaseHistory[]>('/me/purchases', barcode);
			purchases = purchaseList || [];

			// 仕入れ履歴取得
			const restockList = await apiCallWithBarcode<RestockHistory[]>('/me/restocks', barcode);
			restocks = restockList || [];

			pageState = 'viewing';
			activeTab = 'purchases';
			resetTimeout();
		} catch (err) {
			errorMessage = err instanceof Error ? err.message : 'ユーザー情報を取得できません';
		} finally {
			isLoading = false;
		}
	}

	function handleReset() {
		resetSession();
	}

	function formatDate(dateString: string): string {
		return new Date(dateString).toLocaleString('ja-JP');
	}

	function formatCurrency(amount: number): string {
		return `¥${amount.toLocaleString('ja-JP')}`;
	}

	onMount(() => {
		resetTimeout();
		return () => {
			if (timeoutId) clearTimeout(timeoutId);
		};
	});
</script>

<svelte:window on:scan={handleScan} />

<div class="balance-page">
	{#if errorMessage}
		<div class="alert alert-error">
			<strong>エラー:</strong> {errorMessage}
		</div>
	{/if}

	{#if pageState === 'idle'}
		<div class="idle-state">
			<div class="container">
				<h1>残高・履歴確認</h1>
				<p class="subtitle">バーコードをスキャンして確認してください</p>

				<div class="barcode-section">
					<BarcodeInput />
				</div>

				{#if isLoading}
					<div class="loading">処理中...</div>
				{/if}
			</div>
		</div>

	{:else if pageState === 'viewing' && currentUser}
		<div class="viewing-state">
			<div class="container">
				<UserInfo
					balance={currentUser}
					onReset={handleReset}
				/>

				<div class="tabs">
					<button
						class="tab-button"
						class:active={activeTab === 'purchases'}
						onclick={() => (activeTab = 'purchases')}
					>
						購入履歴
					</button>
					<button
						class="tab-button"
						class:active={activeTab === 'restocks'}
						onclick={() => (activeTab = 'restocks')}
					>
						仕入れ履歴
					</button>
				</div>

				<div class="content">
					{#if activeTab === 'purchases'}
						<div class="tab-content">
							<h3>購入履歴</h3>
							{#if purchases.length === 0}
								<p class="empty-notice">購入履歴がありません</p>
							{:else}
								<div class="history-list">
									{#each purchases as purchase (purchase.id)}
										<div class="history-item">
											<div class="history-header">
												<span class="date">{formatDate(purchase.purchased_at)}</span>
												<span class="amount">{formatCurrency(purchase.total_amount)}</span>
											</div>
											<table class="items-table">
												<tbody>
													{#each purchase.items as item (item.product_id)}
														<tr>
															<td class="product-name">{item.product_name}</td>
															<td class="quantity">{item.quantity}</td>
															<td class="unit-price">{formatCurrency(item.unit_price)}</td>
															<td class="subtotal">{formatCurrency(item.subtotal)}</td>
														</tr>
													{/each}
												</tbody>
											</table>
										</div>
									{/each}
								</div>
							{/if}
						</div>

					{:else if activeTab === 'restocks'}
						<div class="tab-content">
							<h3>仕入れ履歴</h3>
							{#if restocks.length === 0}
								<p class="empty-notice">仕入れ履歴がありません</p>
							{:else}
								<div class="history-list">
									{#each restocks as restock (restock.id)}
										<div class="history-item">
											<div class="history-header">
												<span class="date">{formatDate(restock.restocked_at)}</span>
												<span class="amount">{formatCurrency(restock.total_amount)}</span>
											</div>
											{#if restock.note}
												<p class="note">メモ: {restock.note}</p>
											{/if}
											<table class="items-table">
												<tbody>
													{#each restock.items as item (item.product_id)}
														<tr>
															<td class="product-name">{item.product_name}</td>
															<td class="quantity">{item.quantity}</td>
															<td class="unit-price">{formatCurrency(item.unit_price)}</td>
															<td class="subtotal">{formatCurrency(item.subtotal)}</td>
														</tr>
													{/each}
												</tbody>
											</table>
										</div>
									{/each}
								</div>
							{/if}
						</div>
					{/if}
				</div>
			</div>
		</div>
	{/if}
</div>

<style>
	.balance-page {
		flex: 1;
		display: flex;
		flex-direction: column;
		background-color: #f9f9f9;
	}

	.container {
		max-width: 1000px;
		margin: 0 auto;
		padding: 2rem;
		width: 100%;
	}

	.alert {
		max-width: 1000px;
		margin: 0 auto 1.5rem;
		width: 100%;
		padding: 1rem;
		border-radius: 4px;
		font-weight: 500;
	}

	.alert-error {
		background-color: #fee;
		color: #c00;
		border: 1px solid #fcc;
	}

	h1 {
		margin: 0 0 0.5rem 0;
		font-size: 2.5rem;
		text-align: center;
	}

	h3 {
		margin-bottom: 1rem;
		color: #2c3e50;
		font-size: 1.1rem;
	}

	.subtitle {
		text-align: center;
		color: #666;
		margin-bottom: 2rem;
	}

	.idle-state {
		flex: 1;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.barcode-section {
		max-width: 500px;
		margin: 2rem auto 0;
	}

	.loading {
		text-align: center;
		padding: 2rem;
		font-size: 1.1rem;
		color: #0066cc;
		font-weight: 600;
	}

	.viewing-state {
		flex: 1;
		display: flex;
		flex-direction: column;
	}

	.tabs {
		display: flex;
		gap: 0;
		margin-bottom: 1.5rem;
		border-bottom: 2px solid #ddd;
	}

	.tab-button {
		padding: 1rem 2rem;
		background-color: transparent;
		border: none;
		border-bottom: 3px solid transparent;
		color: #666;
		font-size: 1rem;
		font-weight: 600;
		cursor: pointer;
		transition: all 0.2s;
	}

	.tab-button:hover {
		color: #0066cc;
	}

	.tab-button.active {
		color: #0066cc;
		border-bottom-color: #0066cc;
	}

	.content {
		background-color: white;
		border-radius: 8px;
		padding: 1.5rem;
	}

	.tab-content h3 {
		margin-top: 0;
	}

	.empty-notice {
		color: #999;
		text-align: center;
		padding: 2rem 1rem;
		font-size: 0.95rem;
	}

	.history-list {
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}

	.history-item {
		border: 1px solid #eee;
		border-radius: 4px;
		padding: 1rem;
		background-color: #f9f9f9;
	}

	.history-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 1rem;
		padding-bottom: 0.5rem;
		border-bottom: 1px solid #eee;
	}

	.date {
		font-size: 0.9rem;
		color: #666;
	}

	.amount {
		font-size: 1.2rem;
		font-weight: 700;
		color: #0066cc;
	}

	.note {
		color: #666;
		font-size: 0.9rem;
		margin: 0.5rem 0 1rem 0;
		padding: 0.5rem;
		background-color: #fff;
		border-left: 3px solid #ffa500;
		padding-left: 0.75rem;
	}

	.items-table {
		width: 100%;
		border-collapse: collapse;
		font-size: 0.95rem;
	}

	.items-table tbody tr:not(:last-child) {
		border-bottom: 1px solid #eee;
	}

	.items-table td {
		padding: 0.5rem 0;
	}

	.product-name {
		text-align: left;
		color: #2c3e50;
	}

	.quantity,
	.unit-price,
	.subtotal {
		text-align: right;
		color: #666;
	}

	.subtotal {
		color: #0066cc;
		font-weight: 600;
	}

	@media (max-width: 768px) {
		.container {
			padding: 1rem;
		}

		h1 {
			font-size: 1.5rem;
		}

		.tabs {
			gap: 0;
		}

		.tab-button {
			flex: 1;
			text-align: center;
			padding: 0.75rem 1rem;
		}

		.items-table {
			font-size: 0.85rem;
		}

		.items-table td {
			padding: 0.4rem 0;
		}
	}
</style>
