<script lang="ts">
	import { onMount } from 'svelte';
	import { apiCallWithAuth } from '$lib/api';
	import type { ProductWithPrice, ProductPrice } from '$lib/types';

	let products: ProductWithPrice[] = $state([]);
	let isLoading = $state(false);
	let errorMessage = $state('');

	// 追加/編集モーダル
	let showModal = $state(false);
	let modalMode: 'add' | 'edit' = 'add';
	let modalName = $state('');
	let modalBarcode = $state('');
	let modalPrice = $state('');
	let modalNote = $state('');
	let editingProductId: number | null = null;

	// 価格変更モーダル
	let showPriceModal = $state(false);
	let priceProductId: number | null = null;
	let priceProductName = '';
	let newPrice = $state('');

	// 価格履歴モーダル
	let showPriceHistoryModal = $state(false);
	let priceHistory: ProductPrice[] = $state([]);
	let priceHistoryProductName = '';

	onMount(async () => {
		await loadProducts();
	});

	async function loadProducts() {
		isLoading = true;
		errorMessage = '';
		try {
			const data = await apiCallWithAuth<ProductWithPrice[]>('/products');
			products = data || [];
		} catch (err) {
			errorMessage = err instanceof Error ? err.message : '商品一覧の取得に失敗しました';
		} finally {
			isLoading = false;
		}
	}

	function openAddModal() {
		modalMode = 'add';
		modalName = '';
		modalBarcode = '';
		modalPrice = '';
		modalNote = '';
		editingProductId = null;
		showModal = true;
	}

	function openEditModal(product: ProductWithPrice) {
		modalMode = 'edit';
		modalName = product.name;
		modalBarcode = '';
		modalPrice = product.price.toString();
		modalNote = product.note || '';
		editingProductId = product.id;
		showModal = true;
	}

	function closeModal() {
		showModal = false;
	}

	async function handleModalSubmit() {
		if (!modalName.trim()) {
			errorMessage = '商品名を入力してください';
			return;
		}

		errorMessage = '';
		isLoading = true;

		try {
			if (modalMode === 'add') {
				if (!modalBarcode.trim()) {
					errorMessage = 'バーコードを入力してください';
					return;
				}
				if (!modalPrice.trim()) {
					errorMessage = '価格を入力してください';
					return;
				}
				await apiCallWithAuth('/products', {
					method: 'POST',
					body: JSON.stringify({
						name: modalName,
						barcode: modalBarcode,
						price: parseInt(modalPrice),
						note: modalNote || undefined
					})
				});
			} else {
				// edit
				await apiCallWithAuth(`/products/${editingProductId}`, {
					method: 'PUT',
					body: JSON.stringify({ name: modalName, note: modalNote || undefined })
				});
			}
			closeModal();
			await loadProducts();
		} catch (err) {
			errorMessage = err instanceof Error ? err.message : '操作に失敗しました';
		} finally {
			isLoading = false;
		}
	}

	function openPriceModal(product: ProductWithPrice) {
		priceProductId = product.id;
		priceProductName = product.name;
		newPrice = product.price.toString();
		showPriceModal = true;
	}

	function closePriceModal() {
		showPriceModal = false;
	}

	async function handlePriceChange() {
		if (!newPrice.trim() || !priceProductId) {
			errorMessage = '価格を入力してください';
			return;
		}

		errorMessage = '';
		isLoading = true;

		try {
			await apiCallWithAuth(`/products/${priceProductId}/change-price`, {
				method: 'POST',
				body: JSON.stringify({ price: parseInt(newPrice) })
			});
			closePriceModal();
			await loadProducts();
		} catch (err) {
			errorMessage = err instanceof Error ? err.message : '価格変更に失敗しました';
		} finally {
			isLoading = false;
		}
	}

	async function openPriceHistory(product: ProductWithPrice) {
		priceHistoryProductName = product.name;
		showPriceHistoryModal = true;
		isLoading = true;
		errorMessage = '';

		try {
			const data = await apiCallWithAuth<ProductPrice[]>(`/products/${product.id}/prices`);
			priceHistory = data || [];
		} catch (err) {
			errorMessage = err instanceof Error ? err.message : '価格履歴の取得に失敗しました';
		} finally {
			isLoading = false;
		}
	}

	function closePriceHistoryModal() {
		showPriceHistoryModal = false;
	}

	function formatDate(dateString: string): string {
		return new Date(dateString).toLocaleString('ja-JP');
	}

	function formatCurrency(amount: number): string {
		return `¥${amount.toLocaleString('ja-JP')}`;
	}
</script>

<div class="products-page">
	<div class="page-header">
		<h1>商品管理</h1>
		<button class="btn btn-primary" onclick={openAddModal} disabled={isLoading}>
			+ 新規追加
		</button>
	</div>

	{#if errorMessage}
		<div class="alert alert-error">
			<strong>エラー:</strong> {errorMessage}
		</div>
	{/if}

	{#if isLoading && products.length === 0}
		<div class="loading">読み込み中...</div>
	{:else if products.length === 0}
		<div class="empty-notice">商品がいません</div>
	{:else}
		<div class="table-wrapper">
			<table class="products-table">
				<thead>
					<tr>
						<th>ID</th>
						<th>商品名</th>
						<th>バーコード</th>
						<th>現在価格</th>
						<th>状態</th>
						<th>メモ</th>
						<th>操作</th>
					</tr>
				</thead>
				<tbody>
					{#each products as product (product.id)}
						<tr class:inactive={product.is_active === 0}>
							<td>{product.id}</td>
							<td>{product.name}</td>
							<td class="barcode-cell">{product.barcode}</td>
							<td class="price-cell">{formatCurrency(product.price)}</td>
							<td>
								<span class="status-badge" class:active={product.is_active === 1}>
									{product.is_active === 1 ? '有効' : '無効'}
								</span>
							</td>
							<td class="note-cell">{product.note || '—'}</td>
							<td class="action-cell">
								<button
									class="btn-small btn-edit"
									onclick={() => openEditModal(product)}
									disabled={isLoading}
								>
									編集
								</button>
								<button
									class="btn-small btn-price"
									onclick={() => openPriceModal(product)}
									disabled={isLoading}
								>
									価格変更
								</button>
								<button
									class="btn-small btn-history"
									onclick={() => openPriceHistory(product)}
									disabled={isLoading}
								>
									履歴
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
			<h2>{modalMode === 'add' ? '商品を追加' : '商品を編集'}</h2>
			<button class="btn-close" onclick={closeModal}>✕</button>
		</div>

		<form onsubmit={(e) => { e.preventDefault(); handleModalSubmit(); }}>
			<div class="form-group">
				<label for="product-name">商品名</label>
				<input
					type="text"
					id="product-name"
					bind:value={modalName}
					placeholder="商品名"
					required
					disabled={isLoading}
				/>
			</div>

			{#if modalMode === 'add'}
				<div class="form-group">
					<label for="product-barcode">バーコード</label>
					<input
						type="text"
						id="product-barcode"
						bind:value={modalBarcode}
						placeholder="4912345678901"
						required
						disabled={isLoading}
					/>
				</div>

				<div class="form-group">
					<label for="product-price">初期価格（円）</label>
					<input
						type="number"
						id="product-price"
						bind:value={modalPrice}
						placeholder="500"
						required
						min="0"
						disabled={isLoading}
					/>
				</div>
			{/if}

			<div class="form-group">
				<label for="product-note">メモ</label>
				<textarea
					id="product-note"
					bind:value={modalNote}
					placeholder="メモを入力（任意）"
					disabled={isLoading}
					rows="3"
				></textarea>
			</div>

			<div class="modal-footer">
				<button type="button" class="btn btn-secondary" onclick={closeModal} disabled={isLoading}>
					キャンセル
				</button>
				<button type="submit" class="btn btn-primary" disabled={isLoading}>
					{isLoading ? '処理中...' : '保存'}
				</button>
			</div>
		</form>
	</div>
{/if}

<!-- 価格変更モーダル -->
{#if showPriceModal}
	<div class="modal-backdrop" onclick={closePriceModal}></div>
	<div class="modal modal-small">
		<div class="modal-header">
			<h2>価格を変更</h2>
			<button class="btn-close" onclick={closePriceModal}>✕</button>
		</div>

		<form onsubmit={(e) => { e.preventDefault(); handlePriceChange(); }}>
			<div class="form-group">
				<label>商品名</label>
				<p class="product-name-display">{priceProductName}</p>
			</div>

			<div class="form-group">
				<label for="new-price">新しい価格（円）</label>
				<input
					type="number"
					id="new-price"
					bind:value={newPrice}
					required
					min="0"
					disabled={isLoading}
				/>
			</div>

			<div class="modal-footer">
				<button type="button" class="btn btn-secondary" onclick={closePriceModal} disabled={isLoading}>
					キャンセル
				</button>
				<button type="submit" class="btn btn-primary" disabled={isLoading}>
					{isLoading ? '処理中...' : '変更'}
				</button>
			</div>
		</form>
	</div>
{/if}

<!-- 価格履歴モーダル -->
{#if showPriceHistoryModal}
	<div class="modal-backdrop" onclick={closePriceHistoryModal}></div>
	<div class="modal modal-large">
		<div class="modal-header">
			<h2>価格履歴 - {priceHistoryProductName}</h2>
			<button class="btn-close" onclick={closePriceHistoryModal}>✕</button>
		</div>

		<div class="modal-body">
			{#if priceHistory.length === 0}
				<p class="empty-notice">価格履歴がありません</p>
			{:else}
				<div class="price-history-list">
					{#each priceHistory as price (price.id)}
						<div class="price-history-item">
							<div class="price-value">{formatCurrency(price.price)}</div>
							<div class="price-dates">
								<span>開始: {formatDate(price.valid_from)}</span>
								{#if price.valid_to}
									<span>終了: {formatDate(price.valid_to)}</span>
								{:else}
									<span class="current">現在有効</span>
								{/if}
							</div>
						</div>
					{/each}
				</div>
			{/if}
		</div>

		<div class="modal-footer">
			<button class="btn btn-secondary" onclick={closePriceHistoryModal}>
				閉じる
			</button>
		</div>
	</div>
{/if}

<style>
	.products-page {
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
		color: #2c3e50;
		font-size: 1rem;
	}

	.table-wrapper {
		background-color: white;
		border-radius: 8px;
		border: 1px solid #ddd;
		overflow-x: auto;
	}

	.products-table {
		width: 100%;
		border-collapse: collapse;
		font-size: 0.95rem;
	}

	.products-table th {
		background-color: #f5f5f5;
		padding: 1rem;
		text-align: left;
		font-weight: 600;
		border-bottom: 2px solid #ddd;
		color: #2c3e50;
	}

	.products-table td {
		padding: 1rem;
		border-bottom: 1px solid #eee;
	}

	.products-table tbody tr:last-child td {
		border-bottom: none;
	}

	.products-table tbody tr.inactive {
		background-color: #f9f9f9;
		opacity: 0.7;
	}

	.barcode-cell {
		font-family: monospace;
		color: #666;
		font-size: 0.9rem;
	}

	.price-cell {
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
		flex-wrap: wrap;
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

	.btn-small {
		padding: 0.35rem 0.65rem;
		font-size: 0.8rem;
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

	.btn-price {
		background-color: #ff9800;
		color: white;
	}

	.btn-price:hover:not(:disabled) {
		background-color: #f57c00;
	}

	.btn-history {
		background-color: #009688;
		color: white;
	}

	.btn-history:hover:not(:disabled) {
		background-color: #00796b;
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

	.modal-large {
		min-width: 600px;
		max-width: 700px;
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

	.product-name-display {
		margin: 0;
		padding: 0.75rem;
		background-color: #f5f5f5;
		border-radius: 4px;
		color: #2c3e50;
		font-weight: 500;
	}

	.price-history-list {
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}

	.price-history-item {
		padding: 1rem;
		border: 1px solid #eee;
		border-radius: 4px;
		background-color: #f9f9f9;
	}

	.price-value {
		font-size: 1.3rem;
		font-weight: 600;
		color: #0066cc;
		margin-bottom: 0.5rem;
	}

	.price-dates {
		display: flex;
		gap: 1rem;
		font-size: 0.9rem;
		color: #666;
	}

	.price-dates .current {
		background-color: #e8f5e9;
		color: #2e7d32;
		padding: 0.25rem 0.5rem;
		border-radius: 3px;
		font-weight: 500;
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

	input,
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
	textarea:focus {
		outline: none;
		border-color: #0066cc;
		box-shadow: 0 0 0 3px rgba(0, 102, 204, 0.1);
	}

	input:disabled,
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
		.modal-small,
		.modal-large {
			min-width: 90vw;
			max-width: 90vw;
		}

		.products-table {
			font-size: 0.85rem;
		}

		.products-table th,
		.products-table td {
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
