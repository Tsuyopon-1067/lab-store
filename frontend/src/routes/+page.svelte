<script lang="ts">
	import { onMount } from 'svelte';
	import { apiCall, apiCallWithBarcode } from '$lib/api';
	import BarcodeInput from '$lib/components/BarcodeInput.svelte';
	import UserInfo from '$lib/components/UserInfo.svelte';
	import ProductSearch from '$lib/components/ProductSearch.svelte';
	import Cart from '$lib/components/Cart.svelte';
	import PurchaseReceipt from '$lib/components/PurchaseReceipt.svelte';
	import type { UserBalance, Product, CartItem, PurchaseRequest, PurchaseResponse } from '$lib/types';

	type PageState = 'idle' | 'purchasing' | 'receipt';

	// 状態管理
	let pageState: PageState = $state('idle');
	let currentUser: UserBalance | null = $state(null);
	let currentUserBarcode: string = $state('');
	let cart: CartItem[] = $state([]);
	let receipt: PurchaseResponse | null = $state(null);
	let errorMessage = $state('');
	let isLoading = $state(false);

	// タイムアウト設定（ミリ秒）
	const TIMEOUT_MS = 5 * 60 * 1000; // 5分
	let timeoutId: ReturnType<typeof setTimeout> | null = null;

	// タイムアウトをリセット
	function resetTimeout() {
		if (timeoutId) clearTimeout(timeoutId);
		timeoutId = setTimeout(() => {
			resetPurchaseSession();
		}, TIMEOUT_MS);
	}

	// 購買セッションをリセット
	function resetPurchaseSession() {
		pageState = 'idle';
		currentUser = null;
		currentUserBarcode = '';
		cart = [];
		receipt = null;
		errorMessage = '';
		if (timeoutId) clearTimeout(timeoutId);
	}

	// バーコードスキャンハンドラー
	async function handleScan(event: CustomEvent<{ barcode: string }>) {
		resetTimeout();
		const { barcode } = event.detail;

		if (pageState === 'idle') {
			// 利用者スキャン
			await handleUserScan(barcode);
		} else if (pageState === 'purchasing') {
			// 商品スキャン（カート追加）
			await handleProductScan(barcode);
		}
	}

	// 利用者スキャン処理
	async function handleUserScan(barcode: string) {
		isLoading = true;
		errorMessage = '';

		try {
			const balance = await apiCallWithBarcode<UserBalance>(
				'/me/balance',
				barcode,
			);
			currentUser = balance;
			currentUserBarcode = barcode;
			pageState = 'purchasing';
			resetTimeout();
		} catch (err) {
			errorMessage = err instanceof Error ? err.message : '利用者が見つかりません';
		} finally {
			isLoading = false;
		}
	}

	// 商品スキャン処理
	async function handleProductScan(barcode: string) {
		isLoading = true;
		errorMessage = '';

		try {
			const product = await apiCall<Product>(
				`/products/barcode/${encodeURIComponent(barcode)}`,
			);
			addProductToCart(product);
		} catch (err) {
			errorMessage = '商品が見つかりません。管理画面から商品を追加してください';
		} finally {
			isLoading = false;
		}
	}

	// カートに商品を追加
	function addProductToCart(product: Product) {
		const existingItem = cart.find((item) => item.product_id === product.id);

		if (existingItem) {
			existingItem.quantity++;
			existingItem.subtotal = existingItem.quantity * existingItem.unit_price;
		} else {
			const newItem: CartItem = {
				product_id: product.id,
				product_name: product.name,
				quantity: 1,
				unit_price: product.current_price,
				subtotal: product.current_price,
			};
			cart = [...cart, newItem];
		}

		errorMessage = '';
		resetTimeout();
	}

	// 数量変更
	function handleQuantityChange(productId: number, quantity: number) {
		const item = cart.find((item) => item.product_id === productId);
		if (item) {
			if (quantity < 1) {
				cart = cart.filter((item) => item.product_id !== productId);
			} else {
				item.quantity = quantity;
				item.subtotal = quantity * item.unit_price;
			}
		}
		resetTimeout();
	}

	// カートから削除
	function handleRemoveItem(productId: number) {
		cart = cart.filter((item) => item.product_id !== productId);
		resetTimeout();
	}

	// 購入確定
	async function handleConfirmPurchase() {
		if (!currentUser || !currentUserBarcode || cart.length === 0) return;

		isLoading = true;
		errorMessage = '';

		try {
			const purchaseRequest: PurchaseRequest = {
				user_barcode: currentUserBarcode,
				items: cart.map((item) => ({
					product_id: item.product_id,
					quantity: item.quantity,
				})),
			};

			const response = await apiCall<PurchaseResponse>('/purchases', {
				method: 'POST',
				body: JSON.stringify(purchaseRequest),
			});

			receipt = response;
			pageState = 'receipt';
			if (timeoutId) clearTimeout(timeoutId);
		} catch (err) {
			errorMessage = err instanceof Error ? err.message : '購入に失敗しました';
		} finally {
			isLoading = false;
		}
	}

	// レシート表示後、次の購買へ
	function handleReceiptClose() {
		resetPurchaseSession();
	}

	onMount(() => {
		resetTimeout();
		return () => {
			if (timeoutId) clearTimeout(timeoutId);
		};
	});
</script>

<svelte:window on:scan={handleScan} />

<div class="purchase-page">
	<!-- エラー表示 -->
	{#if errorMessage}
		<div class="alert alert-error">
			<strong>エラー:</strong> {errorMessage}
		</div>
	{/if}

	<!-- アイドル状態：利用者スキャン待ち -->
	{#if pageState === 'idle'}
		<div class="idle-state">
			<div class="container">
				<h1>購買画面</h1>
				<p class="subtitle">バーコードリーダーで利用者をスキャンしてください</p>

				<div class="barcode-section">
					<BarcodeInput />
				</div>

				{#if isLoading}
					<div class="loading">処理中...</div>
				{/if}
			</div>
		</div>

	<!-- 購買状態：商品追加＆購入確定 -->
	{:else if pageState === 'purchasing' && currentUser}
		<div class="purchasing-state">
			<div class="container">
				<UserInfo
					balance={currentUser}
					onReset={() => {
						resetPurchaseSession();
					}}
				/>

				<div class="purchase-section">
					<div class="left-panel">
						<ProductSearch
							onProductAdd={addProductToCart}
							onError={(msg) => {
								errorMessage = msg;
							}}
						/>

						<Cart
							items={cart}
							onQuantityChange={handleQuantityChange}
							onRemove={handleRemoveItem}
						/>
					</div>

					<div class="right-panel">
						<div class="action-panel">
							<h3>購入確定</h3>

							{#if cart.length === 0}
								<p class="empty-notice">カートに商品がありません</p>
								<button disabled class="btn-purchase btn-disabled">
									購入確定
								</button>
							{:else}
								<div class="total-section">
									<span class="total-label">合計金額</span>
									<span class="total-amount">
										¥{cart.reduce((sum, item) => sum + item.subtotal, 0).toLocaleString('ja-JP')}
									</span>
								</div>

								<button
									onclick={handleConfirmPurchase}
									disabled={isLoading}
									class="btn-purchase"
								>
									{isLoading ? '処理中...' : '購入確定'}
								</button>
							{/if}
						</div>
					</div>
				</div>
			</div>
		</div>

	<!-- レシート表示状態 -->
	{:else if pageState === 'receipt' && receipt}
		<PurchaseReceipt
			receipt={receipt}
			onClose={handleReceiptClose}
		/>
	{/if}
</div>

<style>
	.purchase-page {
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

	.purchasing-state {
		flex: 1;
		display: flex;
		flex-direction: column;
	}

	.purchase-section {
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

	.action-panel h3 {
		margin: 0 0 1rem 0;
		color: #2c3e50;
		font-size: 1.1rem;
	}

	.empty-notice {
		color: #999;
		text-align: center;
		padding: 1rem 0;
		font-size: 0.95rem;
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

	.btn-purchase {
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

	.btn-purchase:hover:not(:disabled) {
		background-color: #0052a3;
	}

	.btn-purchase:disabled {
		opacity: 0.6;
		cursor: not-allowed;
	}

	.btn-disabled {
		background-color: #ccc;
		cursor: not-allowed;
	}

	@media (max-width: 1024px) {
		.purchase-section {
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

		.purchase-section {
			gap: 1rem;
		}
	}
</style>
