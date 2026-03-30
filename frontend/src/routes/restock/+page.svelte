<script lang="ts">
	import { onMount } from 'svelte';
	import { apiCall, apiCallWithBarcode } from '$lib/api';
	import BarcodeInput from '$lib/components/BarcodeInput.svelte';
	import UserInfo from '$lib/components/UserInfo.svelte';
	import type { UserBalance, Product, RestockItem, RestockRequest, RestockResponse } from '$lib/types';

	type PageState = 'idle' | 'registering' | 'done';

	// 状態管理
	let pageState: PageState = $state('idle');
	let currentUser: UserBalance | null = $state(null);
	let currentUserBarcode: string = $state('');
	let restockItems: RestockItem[] = $state([]);
	let response: RestockResponse | null = $state(null);
	let errorMessage = $state('');
	let isLoading = $state(false);
	let restockNote = $state('');

	// モーダル状態
	let showModal = $state(false);
	let modalProduct: Product | null = $state(null);
	let modalQuantity = $state(1);
	let modalUnitPrice = $state(0);

	// タイムアウト設定
	const TIMEOUT_MS = 5 * 60 * 1000; // 5分
	let timeoutId: ReturnType<typeof setTimeout> | null = null;

	function resetTimeout() {
		if (timeoutId) clearTimeout(timeoutId);
		timeoutId = setTimeout(() => {
			resetRestockSession();
		}, TIMEOUT_MS);
	}

	function resetRestockSession() {
		pageState = 'idle';
		currentUser = null;
		currentUserBarcode = '';
		restockItems = [];
		response = null;
		errorMessage = '';
		restockNote = '';
		showModal = false;
		modalProduct = null;
		if (timeoutId) clearTimeout(timeoutId);
	}

	async function handleScan(event: CustomEvent<{ barcode: string }>) {
		resetTimeout();
		const { barcode } = event.detail;

		if (pageState === 'idle') {
			await handleUserScan(barcode);
		} else if (pageState === 'registering') {
			await handleProductScan(barcode);
		}
	}

	async function handleUserScan(barcode: string) {
		isLoading = true;
		errorMessage = '';

		try {
			const balance = await apiCallWithBarcode<UserBalance>('/me/balance', barcode);
			currentUser = balance;
			currentUserBarcode = barcode;
			pageState = 'registering';
			resetTimeout();
		} catch (err) {
			errorMessage = err instanceof Error ? err.message : '仕入れ者が見つかりません';
		} finally {
			isLoading = false;
		}
	}

	async function handleProductScan(barcode: string) {
		isLoading = true;
		errorMessage = '';

		try {
			const product = await apiCall<Product>(`/products/barcode/${encodeURIComponent(barcode)}`);
			modalProduct = product;
			modalUnitPrice = product.current_price;
			modalQuantity = 1;
			showModal = true;
			resetTimeout();
		} catch (err) {
			errorMessage = err instanceof Error ? err.message : '商品が見つかりません';
		} finally {
			isLoading = false;
		}
	}

	function addToRestockList() {
		if (!modalProduct) return;

		const existingItem = restockItems.find((item) => item.product_id === modalProduct!.id);

		if (existingItem) {
			existingItem.quantity += modalQuantity;
			existingItem.unit_price = modalUnitPrice;
			existingItem.subtotal = existingItem.quantity * existingItem.unit_price;
		} else {
			const newItem: RestockItem = {
				product_id: modalProduct.id,
				product_name: modalProduct.name,
				quantity: modalQuantity,
				unit_price: modalUnitPrice,
				subtotal: modalQuantity * modalUnitPrice,
			};
			restockItems = [...restockItems, newItem];
		}

		showModal = false;
		modalProduct = null;
		errorMessage = '';
		resetTimeout();
	}

	function updateItemQuantity(productId: number, quantity: number) {
		const item = restockItems.find((item) => item.product_id === productId);
		if (item) {
			if (quantity < 1) {
				restockItems = restockItems.filter((item) => item.product_id !== productId);
			} else {
				item.quantity = quantity;
				item.subtotal = quantity * item.unit_price;
			}
		}
		resetTimeout();
	}

	function removeItem(productId: number) {
		restockItems = restockItems.filter((item) => item.product_id !== productId);
		resetTimeout();
	}

	async function confirmRestock() {
		if (!currentUser || !currentUserBarcode || restockItems.length === 0) return;

		isLoading = true;
		errorMessage = '';

		try {
			const totalAmount = restockItems.reduce((sum, item) => sum + item.subtotal, 0);

			const request = {
				user_barcode: currentUserBarcode,
				total_amount: totalAmount,
				note: restockNote,
				items: restockItems.map((item) => ({
					product_id: item.product_id,
					quantity: item.quantity,
					unit_price: item.unit_price,
				})),
			};

			const result = await apiCall<RestockResponse>('/restocks', {
				method: 'POST',
				body: JSON.stringify(request),
			});

			response = result;
			pageState = 'done';
			if (timeoutId) clearTimeout(timeoutId);
		} catch (err) {
			errorMessage = err instanceof Error ? err.message : '仕入れ登録に失敗しました';
		} finally {
			isLoading = false;
		}
	}

	function handleClose() {
		resetRestockSession();
	}

	onMount(() => {
		resetTimeout();
		return () => {
			if (timeoutId) clearTimeout(timeoutId);
		};
	});
</script>

<svelte:window on:scan={handleScan} />

<div class="restock-page">
	{#if errorMessage}
		<div class="alert alert-error">
			<strong>エラー:</strong> {errorMessage}
		</div>
	{/if}

	{#if pageState === 'idle'}
		<div class="idle-state">
			<div class="container">
				<h1>仕入れ登録</h1>
				<p class="subtitle">仕入れ者のバーコードをスキャンしてください</p>

				<div class="barcode-section">
					<BarcodeInput />
				</div>

				{#if isLoading}
					<div class="loading">処理中...</div>
				{/if}
			</div>
		</div>

	{:else if pageState === 'registering' && currentUser}
		<div class="registering-state">
			<div class="container">
				<UserInfo
					balance={currentUser}
					onReset={() => {
						resetRestockSession();
					}}
				/>

				<div class="restock-section">
					<div class="left-panel">
						<div class="product-input-section">
							<h3>商品追加</h3>
							<p class="hint">商品バーコードをスキャンしてください</p>
							<BarcodeInput />
						</div>

						<div class="restock-list">
							<h3>仕入れリスト</h3>
							{#if restockItems.length === 0}
								<p class="empty-notice">まだ商品が追加されていません</p>
							{:else}
								<table class="items-table">
									<thead>
										<tr>
											<th>商品名</th>
											<th>数量</th>
											<th>単価</th>
											<th>小計</th>
											<th>操作</th>
										</tr>
									</thead>
									<tbody>
										{#each restockItems as item (item.product_id)}
											<tr>
												<td>{item.product_name}</td>
												<td>
													<input
														type="number"
														min="1"
														bind:value={item.quantity}
														onchange={() => {
															item.subtotal = item.quantity * item.unit_price;
															restockItems = restockItems;
														}}
														class="qty-input"
													/>
												</td>
												<td>¥{item.unit_price.toLocaleString('ja-JP')}</td>
												<td>¥{item.subtotal.toLocaleString('ja-JP')}</td>
												<td>
													<button onclick={() => removeItem(item.product_id)} class="btn-remove">削除</button>
												</td>
											</tr>
										{/each}
									</tbody>
								</table>
							{/if}
						</div>
					</div>

					<div class="right-panel">
						<div class="action-panel">
							<h3>仕入れ確定</h3>

							{#if restockItems.length === 0}
								<p class="empty-notice">リストに商品を追加してください</p>
								<button disabled class="btn-confirm btn-disabled">仕入れ確定</button>
							{:else}
								<div class="total-section">
									<span class="total-label">合計金額</span>
									<span class="total-amount">
										¥{restockItems.reduce((sum, item) => sum + item.subtotal, 0).toLocaleString('ja-JP')}
									</span>
								</div>

								<div class="form-group">
									<label for="note">メモ（領収書の内容など）</label>
									<textarea
										id="note"
										bind:value={restockNote}
										placeholder="仕入れの詳細（例：○○店での購入）"
										class="form-textarea"
									></textarea>
								</div>

								<button
									onclick={confirmRestock}
									disabled={isLoading}
									class="btn-confirm"
								>
									{isLoading ? '処理中...' : '仕入れ確定'}
								</button>
							{/if}
						</div>
					</div>
				</div>
			</div>
		</div>

	{:else if pageState === 'done' && response}
		<div class="done-state">
			<div class="container">
				<div class="receipt">
					<h2>仕入れ登録完了</h2>
					<div class="receipt-details">
						<p><strong>仕入れ者:</strong> {response.user_name}</p>
						<p><strong>日時:</strong> {new Date(response.restocked_at).toLocaleString('ja-JP')}</p>
						<p><strong>合計:</strong> ¥{response.total_amount.toLocaleString('ja-JP')}</p>
					</div>

					<table class="items-table">
						<thead>
							<tr>
								<th>商品名</th>
								<th>数量</th>
								<th>単価</th>
								<th>小計</th>
							</tr>
						</thead>
						<tbody>
							{#each response.items as item (item.product_id)}
								<tr>
									<td>{item.product_name}</td>
									<td>{item.quantity}</td>
									<td>¥{item.unit_price.toLocaleString('ja-JP')}</td>
									<td>¥{item.subtotal.toLocaleString('ja-JP')}</td>
								</tr>
							{/each}
						</tbody>
					</table>

					<button onclick={handleClose} class="btn-close">OK</button>
				</div>
			</div>
		</div>
	{/if}

	<!-- 数量・単価入力モーダル -->
	{#if showModal && modalProduct}
		<div class="modal-overlay" onclick={() => (showModal = false)}>
			<div class="modal" onclick={(e) => e.stopPropagation()}>
				<h3>{modalProduct.name}</h3>

				<div class="form-group">
					<label for="quantity">数量</label>
					<input
						id="quantity"
						type="number"
						min="1"
						bind:value={modalQuantity}
						class="form-input"
					/>
				</div>

				<div class="form-group">
					<label for="unitPrice">単価</label>
					<input
						id="unitPrice"
						type="number"
						min="0"
						bind:value={modalUnitPrice}
						class="form-input"
					/>
				</div>

				<div class="modal-actions">
					<button onclick={() => (showModal = false)} class="btn-cancel">キャンセル</button>
					<button onclick={addToRestockList} class="btn-add">追加</button>
				</div>
			</div>
		</div>
	{/if}
</div>

<style>
	.restock-page {
		flex: 1;
		display: flex;
		flex-direction: column;
		background-color: #f9f9f9;
	}

	.container {
		max-width: 1200px;
		margin: 0 auto;
		padding: 2rem;
		width: 100%;
	}

	.alert {
		max-width: 1200px;
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

	.registering-state {
		flex: 1;
		display: flex;
		flex-direction: column;
	}

	.restock-section {
		display: grid;
		grid-template-columns: 1fr 350px;
		gap: 2rem;
		flex: 1;
	}

	.left-panel {
		display: flex;
		flex-direction: column;
		gap: 1.5rem;
	}

	.product-input-section {
		background-color: white;
		border: 1px solid #ddd;
		border-radius: 8px;
		padding: 1.5rem;
	}

	.product-input-section p.hint {
		color: #666;
		font-size: 0.95rem;
		margin-top: 0.5rem;
	}

	.restock-list {
		background-color: white;
		border: 1px solid #ddd;
		border-radius: 8px;
		padding: 1.5rem;
	}

	.empty-notice {
		color: #999;
		text-align: center;
		padding: 1rem 0;
		font-size: 0.95rem;
	}

	.items-table {
		width: 100%;
		border-collapse: collapse;
		margin-top: 1rem;
	}

	.items-table th,
	.items-table td {
		padding: 0.75rem;
		text-align: left;
		border-bottom: 1px solid #eee;
	}

	.items-table th {
		background-color: #f5f5f5;
		font-weight: 600;
		color: #2c3e50;
	}

	.items-table tbody tr:hover {
		background-color: #f9f9f9;
	}

	.qty-input {
		width: 60px;
		padding: 0.3rem;
		border: 1px solid #ddd;
		border-radius: 4px;
		text-align: center;
	}

	.btn-remove {
		padding: 0.3rem 0.6rem;
		background-color: #ff6b6b;
		color: white;
		border: none;
		border-radius: 4px;
		font-size: 0.85rem;
		cursor: pointer;
		transition: background-color 0.2s;
	}

	.btn-remove:hover {
		background-color: #ff5252;
	}

	.right-panel {
		display: flex;
		flex-direction: column;
	}

	.action-panel {
		background-color: white;
		border: 2px solid #0066cc;
		border-radius: 8px;
		padding: 1.5rem;
		position: sticky;
		top: 100px;
	}

	.total-section {
		background-color: #f0f8ff;
		padding: 1rem;
		border-radius: 4px;
		margin-bottom: 1rem;
		display: flex;
		justify-content: space-between;
		align-items: center;
	}

	.total-label {
		font-weight: 600;
		color: #2c3e50;
	}

	.total-amount {
		font-size: 1.5rem;
		font-weight: 700;
		color: #0066cc;
	}

	.form-group {
		margin-bottom: 1.5rem;
	}

	.form-group label {
		display: block;
		margin-bottom: 0.5rem;
		font-weight: 600;
		color: #2c3e50;
	}

	.form-textarea {
		width: 100%;
		padding: 0.75rem;
		border: 1px solid #ddd;
		border-radius: 4px;
		font-size: 0.9rem;
		font-family: inherit;
		box-sizing: border-box;
		resize: vertical;
		min-height: 60px;
	}

	.form-textarea:focus {
		outline: none;
		border-color: #0066cc;
		box-shadow: 0 0 0 3px rgba(0, 102, 204, 0.15);
	}

	.btn-confirm {
		width: 100%;
		padding: 1rem;
		background-color: #0066cc;
		color: white;
		border: none;
		border-radius: 4px;
		font-size: 1.1rem;
		font-weight: 600;
		cursor: pointer;
		transition: background-color 0.2s;
	}

	.btn-confirm:hover:not(:disabled) {
		background-color: #0052a3;
	}

	.btn-confirm:disabled {
		opacity: 0.6;
		cursor: not-allowed;
	}

	.btn-disabled {
		background-color: #ccc;
		cursor: not-allowed;
	}

	.done-state {
		flex: 1;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.receipt {
		background-color: white;
		border: 2px solid #0066cc;
		border-radius: 8px;
		padding: 2rem;
		max-width: 600px;
		width: 100%;
	}

	.receipt h2 {
		text-align: center;
		margin: 0 0 1.5rem 0;
		color: #0066cc;
	}

	.receipt-details {
		background-color: #f9f9f9;
		padding: 1rem;
		border-radius: 4px;
		margin-bottom: 1.5rem;
		border-left: 4px solid #0066cc;
	}

	.receipt-details p {
		margin: 0.5rem 0;
		font-size: 0.95rem;
	}

	.btn-close {
		width: 100%;
		padding: 1rem;
		background-color: #0066cc;
		color: white;
		border: none;
		border-radius: 4px;
		font-size: 1.1rem;
		font-weight: 600;
		cursor: pointer;
		margin-top: 1rem;
	}

	.btn-close:hover {
		background-color: #0052a3;
	}

	/* モーダルスタイル */
	.modal-overlay {
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		background-color: rgba(0, 0, 0, 0.5);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 1000;
	}

	.modal {
		background-color: white;
		border-radius: 8px;
		padding: 2rem;
		max-width: 400px;
		width: 90%;
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
	}

	.modal h3 {
		margin: 0 0 1.5rem 0;
		color: #2c3e50;
	}

	.form-input {
		width: 100%;
		padding: 0.75rem;
		border: 1px solid #ddd;
		border-radius: 4px;
		font-size: 1rem;
		box-sizing: border-box;
	}

	.form-input:focus {
		outline: none;
		border-color: #0066cc;
		box-shadow: 0 0 0 3px rgba(0, 102, 204, 0.15);
	}

	.modal-actions {
		display: flex;
		gap: 1rem;
	}

	.btn-cancel,
	.btn-add {
		flex: 1;
		padding: 0.75rem;
		border: none;
		border-radius: 4px;
		font-size: 1rem;
		font-weight: 600;
		cursor: pointer;
		transition: background-color 0.2s;
	}

	.btn-cancel {
		background-color: #ccc;
		color: white;
	}

	.btn-cancel:hover {
		background-color: #999;
	}

	.btn-add {
		background-color: #0066cc;
		color: white;
	}

	.btn-add:hover {
		background-color: #0052a3;
	}

	@media (max-width: 1024px) {
		.restock-section {
			grid-template-columns: 1fr;
		}

		.action-panel {
			position: static;
		}
	}

	@media (max-width: 768px) {
		.container {
			padding: 1rem;
		}

		h1 {
			font-size: 1.5rem;
		}

		.restock-section {
			gap: 1rem;
		}
	}
</style>
