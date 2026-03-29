<script lang="ts">
	/**
	 * バーコードリーダー入力を受け付けるコンポーネント
	 * USBバーコードリーダー（キーボード入力エミュレーション）を想定
	 *
	 * 使用方法:
	 * <BarcodeInput on:scan={(e) => handleScan(e.detail)} />
	 */

	import { onMount } from 'svelte';

	interface ScanEvent {
		barcode: string;
	}

	let input = $state('');
	let inputRef: HTMLInputElement | undefined;

	// バーコードスキャン完了時のイベント発火
	const dispatchScan = (barcode: string) => {
		const event = new CustomEvent<ScanEvent>('scan', {
			detail: { barcode },
		});
		window.dispatchEvent(event);
	};

	// Enterキーでスキャン完了と判定
	function handleKeyDown(e: KeyboardEvent) {
		if (e.key === 'Enter') {
			const barcode = input.trim();
			if (barcode.length > 0) {
				dispatchScan(barcode);
				input = '';
			}
		}
	}

	onMount(() => {
		// コンポーネントマウント時にフォーカス
		inputRef?.focus();
	});
</script>

<input
	bind:this={inputRef}
	bind:value={input}
	type="text"
	placeholder="バーコードをスキャンしてください"
	on:keydown={handleKeyDown}
	aria-label="バーコードスキャン入力"
	class="barcode-input"
/>

<style>
	.barcode-input {
		width: 100%;
		padding: 0.75rem;
		border: 2px solid #0066cc;
		border-radius: 4px;
		font-size: 1.1rem;
		text-align: center;
		letter-spacing: 0.1em;
	}

	.barcode-input:focus {
		outline: none;
		border-color: #0052a3;
		box-shadow: 0 0 0 3px rgba(0, 102, 204, 0.15);
	}
</style>
