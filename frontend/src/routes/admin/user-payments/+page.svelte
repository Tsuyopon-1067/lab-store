<script lang="ts">
	import { onMount } from 'svelte';
	import { apiCallWithAuth } from '$lib/api';
	import type { UserBalance } from '$lib/types';

	let summaries: UserBalance[] = $state([]);
	let isLoading = $state(false);
	let errorMessage = $state('');

	onMount(async () => {
		await loadSummaries();
	});

	async function loadSummaries() {
		isLoading = true;
		errorMessage = '';
		try {
			const data = await apiCallWithAuth<UserBalance[]>('/users/payment-summary');
			summaries = data || [];
		} catch (err) {
			errorMessage = err instanceof Error ? err.message : '支払い金額一覧の取得に失敗しました';
		} finally {
			isLoading = false;
		}
	}

	function formatAmount(amount: number): string {
		return `¥${amount.toLocaleString('ja-JP')}`;
	}

	async function exportToCSV() {
		try {
			const response = await fetch('/api/users/payment-summary/export', {
				method: 'GET',
				credentials: 'include',
			});

			if (!response.ok) {
				throw new Error('CSV export failed');
			}

			const blob = await response.blob();
			const url = window.URL.createObjectURL(blob);
			const link = document.createElement('a');
			link.href = url;
			link.download = `user_payment_summary_${new Date().toISOString().split('T')[0]}.csv`;
			document.body.appendChild(link);
			link.click();
			document.body.removeChild(link);
			window.URL.revokeObjectURL(url);
		} catch (err) {
			errorMessage = err instanceof Error ? err.message : 'CSV出力に失敗しました';
		}
	}
</script>

<div class="user-payments-page">
	<div class="header">
		<h1>支払い金額一覧</h1>
		{#if summaries.length > 0}
			<button class="btn btn-export" onclick={exportToCSV} disabled={isLoading}>
				📥 CSV出力
			</button>
		{/if}
	</div>

	{#if errorMessage}
		<div class="alert alert-error">
			<strong>エラー:</strong>
			{errorMessage}
		</div>
	{/if}

	{#if isLoading && summaries.length === 0}
		<div class="loading">読み込み中...</div>
	{:else if summaries.length === 0}
		<div class="empty-notice">支払い金額情報がありません</div>
	{:else}
		<div class="table-wrapper">
			<table class="summaries-table">
				<thead>
					<tr>
						<th>ユーザー名</th>
						<th>利用額</th>
						<th>仕入れ金額</th>
						<th>最終的な支払い金額</th>
					</tr>
				</thead>
				<tbody>
					{#each summaries as summary (summary.user_id)}
						<tr class={summary.net_balance > 0 ? 'positive' : summary.net_balance < 0 ? 'negative' : 'neutral'}>
							<td>{summary.user_name}</td>
							<td class="amount">{formatAmount(summary.purchase_unpaid)}</td>
							<td class="amount">{formatAmount(summary.restock_unclaimed)}</td>
							<td class="amount net">{formatAmount(summary.net_balance)}</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	{/if}
</div>

<style>
	.user-payments-page {
		max-width: 1200px;
	}

	.header {
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
		color: #eee;
		font-size: 1rem;
	}

	.table-wrapper {
		background-color: white;
		border-radius: 8px;
		border: 1px solid #ddd;
		overflow-x: auto;
		margin-bottom: 2rem;
	}

	.summaries-table {
		width: 100%;
		border-collapse: collapse;
		font-size: 0.95rem;
	}

	.summaries-table th {
		background-color: #f5f5f5;
		padding: 1rem;
		text-align: left;
		font-weight: 600;
		border-bottom: 2px solid #ddd;
		color: #2c3e50;
	}

	.summaries-table td {
		padding: 1rem;
		border-bottom: 1px solid #eee;
	}

	.summaries-table tbody tr:hover {
		background-color: #f9f9f9;
	}

	.summaries-table .amount {
		text-align: right;
		font-weight: 600;
	}

	.summaries-table .amount.net {
		color: #0066cc;
		font-size: 1.05rem;
	}

	.summaries-table tr.positive .amount.net {
		color: #d32f2f;
	}

	.summaries-table tr.negative .amount.net {
		color: #388e3c;
	}

	.btn {
		padding: 0.5rem 1rem;
		border: none;
		border-radius: 4px;
		font-weight: 500;
		cursor: pointer;
		transition: background-color 0.2s;
	}

	.btn-export {
		background-color: #0066cc;
		color: white;
		border: none;
	}

	.btn-export:hover:not(:disabled) {
		background-color: #0052a3;
	}

	.btn-export:disabled {
		opacity: 0.6;
		cursor: not-allowed;
	}

	@media (max-width: 768px) {
		.summaries-table {
			font-size: 0.85rem;
		}

		.summaries-table th,
		.summaries-table td {
			padding: 0.75rem;
		}

		.header {
			flex-wrap: wrap;
			gap: 1rem;
		}

		.btn-export {
			width: 100%;
		}
	}
</style>
