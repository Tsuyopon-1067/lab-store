<script lang="ts">
	import { apiCall } from '$lib/api';
	import type { Product } from '$lib/types';

	interface Props {
		onProductAdd: (product: Product) => void;
		onError: (message: string) => void;
	}

	let { onProductAdd, onError }: Props = $props();

	let searchQuery = $state('');
	let isSearching = $state(false);
	let searchResults: Product[] = $state([]);
	let showResults = $state(false);

	async function handleBarcodeSearch(barcode: string) {
		isSearching = true;
		try {
			const product = await apiCall<Product>(`/products/barcode/${encodeURIComponent(barcode)}`);
			onProductAdd(product);
			searchQuery = '';
			searchResults = [];
			showResults = false;
		} catch (err) {
			onError(err instanceof Error ? err.message : '商品が見つかりません');
		} finally {
			isSearching = false;
		}
	}

	async function handleSearch() {
		if (!searchQuery.trim()) {
			searchResults = [];
			showResults = false;
			return;
		}

		isSearching = true;
		try {
			const products = await apiCall<Product[]>(
				`/products?search=${encodeURIComponent(searchQuery)}`,
			);
			searchResults = products;
			showResults = true;
		} catch (err) {
			onError(err instanceof Error ? err.message : '検索に失敗しました');
			searchResults = [];
		} finally {
			isSearching = false;
		}
	}

	function handleProductSelect(product: Product) {
		onProductAdd(product);
		searchQuery = '';
		searchResults = [];
		showResults = false;
	}

	// グローバルスキャンイベントをリッスン
	$effect.pre(() => {
		const handleScan = (e: Event) => {
			const scanEvent = e as CustomEvent<{ barcode: string }>;
			const barcode = scanEvent.detail.barcode;
			handleBarcodeSearch(barcode);
		};

		window.addEventListener('scan', handleScan);

		return () => {
			window.removeEventListener('scan', handleScan);
		};
	});
</script>

<div class="product-search">
	<div class="search-controls">
		<input
			type="text"
			bind:value={searchQuery}
			placeholder="商品名で検索..."
			onkeydown={(e) => e.key === 'Enter' && handleSearch()}
			disabled={isSearching}
			class="search-input"
		/>
		<button onclick={handleSearch} disabled={isSearching} class="btn-search">
			{isSearching ? '検索中...' : '検索'}
		</button>
	</div>

	{#if showResults && searchResults.length > 0}
		<div class="search-results">
			{#each searchResults as product (product.id)}
				<button
					onclick={() => handleProductSelect(product)}
					class="result-item"
				>
					<div class="result-name">{product.name}</div>
					<div class="result-price">¥{product.current_price.toLocaleString('ja-JP')}</div>
				</button>
			{/each}
		</div>
	{:else if showResults && searchQuery}
		<div class="no-results">
			検索結果がありません
		</div>
	{/if}
</div>

<style>
	.product-search {
		background-color: white;
		border: 1px solid #ddd;
		border-radius: 8px;
		padding: 1.5rem;
		margin-bottom: 2rem;
	}

	.search-controls {
		display: flex;
		gap: 0.5rem;
	}

	.search-input {
		flex: 1;
		padding: 0.75rem;
		border: 1px solid #ccc;
		border-radius: 4px;
		font-size: 1rem;
	}

	.search-input:focus {
		outline: none;
		border-color: #0066cc;
		box-shadow: 0 0 0 3px rgba(0, 102, 204, 0.1);
	}

	.btn-search {
		background-color: #0066cc;
		color: white;
		padding: 0.75rem 1.5rem;
		border: none;
		border-radius: 4px;
		cursor: pointer;
		font-weight: 500;
		transition: background-color 0.2s;
	}

	.btn-search:hover:not(:disabled) {
		background-color: #0052a3;
	}

	.btn-search:disabled {
		opacity: 0.6;
		cursor: not-allowed;
	}

	.search-results {
		margin-top: 1rem;
		border-top: 1px solid #f0f0f0;
		padding-top: 1rem;
	}

	.result-item {
		display: block;
		width: 100%;
		padding: 0.75rem;
		margin-bottom: 0.5rem;
		border: 1px solid #ddd;
		border-radius: 4px;
		background-color: #f9f9f9;
		cursor: pointer;
		text-align: left;
		transition: background-color 0.2s;
	}

	.result-item:hover {
		background-color: #f0f0f0;
	}

	.result-name {
		font-weight: 600;
		color: #2c3e50;
		margin-bottom: 0.25rem;
	}

	.result-price {
		font-size: 0.9rem;
		color: #0066cc;
	}

	.no-results {
		padding: 1rem;
		text-align: center;
		color: #999;
	}
</style>
