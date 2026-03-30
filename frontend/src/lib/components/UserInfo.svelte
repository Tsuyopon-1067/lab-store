<script lang="ts">
	import type { UserBalance } from '$lib/types';

	interface Props {
		balance: UserBalance;
		onReset?: () => void;
	}

	let { balance, onReset }: Props = $props();

	const formatCurrency = (amount: number) => {
		return `¥${amount.toLocaleString('ja-JP')}`;
	};
</script>

<div class="user-info">
	<div class="user-header">
		{#if onReset}
			<button onclick={onReset} class="btn-reset">戻る</button>
		{/if}
		<h2>{balance.user_name}</h2>
	</div>

	<div class="balance-details">
		<div class="balance-item">
			<span class="label">購入未払い額</span>
			<span class="amount">{formatCurrency(balance.purchase_unpaid)}</span>
		</div>

		<div class="balance-item">
			<span class="label">仕入れ立替未精算額</span>
			<span class="amount">{formatCurrency(balance.restock_unclaimed)}</span>
		</div>

		<div class="balance-item total">
			<span class="label">
				{balance.net_balance > 0 ? '今月の支払い予定額' : '今月の受取予定額'}
			</span>
			<span class="amount" class:negative={balance.net_balance < 0}>
				{formatCurrency(Math.abs(balance.net_balance))}
			</span>
		</div>
	</div>
</div>

<style>
	.user-info {
		background-color: white;
		border: 1px solid #ddd;
		border-radius: 8px;
		padding: 1.5rem;
		margin-bottom: 2rem;
	}

	.user-header {
		display: flex;
		align-items: center;
		gap: 1rem;
		margin-bottom: 1.5rem;
		border-bottom: 2px solid #f0f0f0;
		padding-bottom: 1rem;
	}

	h2 {
		margin: 0;
		font-size: 1.8rem;
		color: #2c3e50;
		flex: 1;
		text-align: center;
	}

	.btn-reset {
		background-color: #ccc;
		color: white;
		padding: 0.5rem 1rem;
		border: none;
		border-radius: 4px;
		cursor: pointer;
		font-size: 0.9rem;
		transition: background-color 0.2s;
		flex-shrink: 0;
	}

	.btn-reset:hover {
		background-color: #999;
	}

	.balance-details {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
		gap: 1rem;
	}

	.balance-item {
		padding: 1rem;
		background-color: #f9f9f9;
		border-radius: 4px;
		display: flex;
		justify-content: space-between;
		align-items: center;
	}

	.balance-item.total {
		background-color: #e8f4f8;
		font-weight: 600;
		grid-column: 1 / -1;
	}

	.label {
		color: #666;
		font-size: 0.95rem;
	}

	.amount {
		font-size: 1.3rem;
		color: #0066cc;
		font-weight: 600;
	}

	.amount.negative {
		color: #f97316;
	}
</style>
