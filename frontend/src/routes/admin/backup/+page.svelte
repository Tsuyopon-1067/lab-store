<script lang="ts">
	import { onMount } from "svelte";
	import { apiCallWithAuth } from "$lib/api";
	import type { Backup } from "$lib/types";

	let backups: Backup[] = $state([]);
	let isLoading = $state(false);
	let errorMessage = $state("");

	onMount(async () => {
		await loadBackups();
	});

	async function loadBackups() {
		isLoading = true;
		errorMessage = "";
		try {
			const data = await apiCallWithAuth<Backup[]>("/backup/list");
			backups = data || [];
		} catch (err) {
			errorMessage =
				err instanceof Error
					? err.message
					: "バックアップ一覧の取得に失敗しました";
		} finally {
			isLoading = false;
		}
	}

	async function handleCreateBackup() {
		isLoading = true;
		errorMessage = "";
		try {
			const backup = await apiCallWithAuth<Backup>("/backup", {
				method: "POST",
			});
			backups = [backup, ...backups];
		} catch (err) {
			errorMessage =
				err instanceof Error
					? err.message
					: "バックアップの作成に失敗しました";
		} finally {
			isLoading = false;
		}
	}

	function formatDate(dateString: string): string {
		const date = new Date(dateString);
		return date.toLocaleString("ja-JP");
	}
</script>

<div class="backup-page">
	<div class="page-header">
		<h1>バックアップ管理</h1>
		<button
			class="btn btn-primary"
			disabled={isLoading}
			onclick={handleCreateBackup}
		>
			{isLoading ? "作成中..." : "バックアップを作成"}
		</button>
	</div>

	{#if errorMessage}
		<div class="alert alert-error">{errorMessage}</div>
	{/if}

	{#if isLoading && backups.length === 0}
		<div class="loading">読み込み中...</div>
	{:else if backups.length === 0}
		<div class="empty-state">
			<p>バックアップレコードがありません</p>
		</div>
	{:else}
		<div class="table-wrapper">
			<table>
				<thead>
					<tr>
						<th>ID</th>
						<th>ファイル名</th>
						<th>作成日時</th>
					</tr>
				</thead>
				<tbody>
					{#each backups as backup (backup.id)}
						<tr>
							<td>{backup.id}</td>
							<td>{backup.filename}</td>
							<td>{formatDate(backup.created_at)}</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	{/if}
</div>

<style>
	.backup-page {
		max-width: 1000px;
	}

	.page-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 2rem;
		gap: 1rem;
	}

	h1 {
		margin: 0;
		font-size: 1.8rem;
		color: #fff;
		flex: 1;
	}

	.btn {
		padding: 0.5rem 1rem;
		border-radius: 4px;
		border: none;
		cursor: pointer;
		font-weight: 500;
		transition: background-color 0.2s;
		font-size: 0.95rem;
	}

	.btn-primary {
		background-color: #0066cc;
		color: white;
	}

	.btn-primary:hover:not(:disabled) {
		background-color: #0052a3;
	}

	.btn:disabled {
		opacity: 0.6;
		cursor: not-allowed;
	}

	.alert {
		padding: 1rem;
		border-radius: 4px;
		margin-bottom: 1rem;
		font-size: 0.95rem;
	}

	.alert-error {
		background-color: #ffebee;
		border: 1px solid #ef5350;
		color: #c62828;
	}

	.loading {
		text-align: center;
		padding: 2rem;
		color: #999;
		font-size: 0.95rem;
	}

	.empty-state {
		text-align: center;
		padding: 2rem;
		background-color: white;
		border: 1px solid #ddd;
		border-radius: 8px;
		color: #666;
		font-size: 0.95rem;
	}

	.table-wrapper {
		background-color: white;
		border-radius: 8px;
		overflow: hidden;
		box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
	}

	table {
		width: 100%;
		border-collapse: collapse;
	}

	thead {
		background-color: #f5f5f5;
		border-bottom: 2px solid #ddd;
	}

	th {
		padding: 1rem;
		text-align: left;
		font-weight: 500;
		color: #2c3e50;
		font-size: 0.9rem;
	}

	td {
		padding: 0.75rem 1rem;
		border-bottom: 1px solid #eee;
		font-size: 0.9rem;
		color: #333;
	}

	tbody tr:last-child td {
		border-bottom: none;
	}

	tbody tr:hover {
		background-color: #f9f9f9;
	}
</style>
