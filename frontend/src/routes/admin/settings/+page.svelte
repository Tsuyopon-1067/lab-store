<script lang="ts">
	import { onMount } from 'svelte';
	import { apiCallWithAuth } from '$lib/api';
	import type { Setting } from '$lib/types';

	// --- loaded state ---
	let isLoading = $state(true);
	let isSaving = $state(false);
	let errorMessage = $state('');
	let successMessage = $state('');

	// --- form fields ---
	let sessionTimeoutMinutes = $state(30);
	let backupPath = $state('../backup');
	let backupIntervalMinutes = $state(60);

	onMount(async () => {
		await loadSettings();
	});

	async function loadSettings() {
		isLoading = true;
		errorMessage = '';
		try {
			const data = await apiCallWithAuth<Setting>('/settings');
			sessionTimeoutMinutes = data.session_timeout_minutes;
			backupPath = data.backup_path;
			backupIntervalMinutes = data.backup_interval_minutes;
		} catch (err) {
			errorMessage = err instanceof Error ? err.message : '設定の読み込みに失敗しました';
		} finally {
			isLoading = false;
		}
	}

	async function handleSave() {
		errorMessage = '';
		successMessage = '';

		// Client-side validation
		if (sessionTimeoutMinutes < 1) {
			errorMessage = 'セッションタイムアウトは1分以上を入力してください';
			return;
		}
		if (backupIntervalMinutes < 10) {
			errorMessage = 'バックアップ間隔は10分以上を入力してください';
			return;
		}
		if (!backupPath.trim()) {
			errorMessage = 'バックアップパスを入力してください';
			return;
		}

		isSaving = true;
		try {
			const req = {
				session_timeout_minutes: sessionTimeoutMinutes,
				backup_path: backupPath.trim(),
				backup_interval_minutes: backupIntervalMinutes
			};
			const updated = await apiCallWithAuth<Setting>('/settings', {
				method: 'PUT',
				body: JSON.stringify(req)
			});
			// Sync form from server response (server may have clamped values)
			sessionTimeoutMinutes = updated.session_timeout_minutes;
			backupPath = updated.backup_path;
			backupIntervalMinutes = updated.backup_interval_minutes;
			successMessage = '設定を保存しました';
			// Auto-dismiss after 3 seconds
			setTimeout(() => {
				successMessage = '';
			}, 3000);
		} catch (err) {
			errorMessage = err instanceof Error ? err.message : '設定の保存に失敗しました';
		} finally {
			isSaving = false;
		}
	}
</script>

<div class="settings-page">
	<h1>設定</h1>

	{#if errorMessage}
		<div class="alert alert-error">{errorMessage}</div>
	{/if}

	{#if successMessage}
		<div class="alert alert-success">{successMessage}</div>
	{/if}

	{#if isLoading}
		<div class="loading">読み込み中...</div>
	{:else}
		<form
			onsubmit={(e) => {
				e.preventDefault();
				handleSave();
			}}
		>
			<section class="settings-section">
				<h2>セッション</h2>
				<div class="form-group">
					<label for="session-timeout"> セッションタイムアウト（分） </label>
					<div class="input-hint">
						管理者セッションの有効期限です。変更は次回ログイン時から反映されます。
					</div>
					<input
						type="number"
						id="session-timeout"
						bind:value={sessionTimeoutMinutes}
						min="1"
						max="1440"
						disabled={isSaving}
					/>
				</div>
			</section>

			<section class="settings-section">
				<h2>バックアップ</h2>
				<div class="form-group">
					<label for="backup-path">バックアップ保存先パス</label>
					<div class="input-hint">
						バックアップファイルを保存するディレクトリのパスです（相対パスまたは絶対パス）。
					</div>
					<input
						type="text"
						id="backup-path"
						bind:value={backupPath}
						placeholder="../backup"
						disabled={isSaving}
					/>
				</div>
				<div class="form-group">
					<label for="backup-interval"> バックアップ間隔（分） </label>
					<div class="input-hint">
						自動バックアップの実行間隔です（最小10分）。変更は次の実行サイクルから反映されます。
					</div>
					<input
						type="number"
						id="backup-interval"
						bind:value={backupIntervalMinutes}
						min="10"
						disabled={isSaving}
					/>
				</div>
			</section>

			<div class="form-actions">
				<button type="submit" class="btn btn-primary" disabled={isSaving}>
					{isSaving ? '保存中...' : '保存'}
				</button>
			</div>
		</form>
	{/if}
</div>

<style>
	.settings-page {
		max-width: 700px;
	}

	h1 {
		margin: 0 0 2rem 0;
		font-size: 1.8rem;
		color: #fff;
	}

	h2 {
		margin: 0 0 1.25rem 0;
		font-size: 1.1rem;
		color: #2c3e50;
		font-weight: 600;
	}

	.alert {
		padding: 1rem;
		border-radius: 4px;
		margin-bottom: 1.5rem;
		font-size: 0.95rem;
	}

	.alert-error {
		background-color: #ffebee;
		color: #c62828;
		border: 1px solid #ef5350;
	}

	.alert-success {
		background-color: #e8f5e9;
		color: #2e7d32;
		border: 1px solid #66bb6a;
	}

	.loading {
		text-align: center;
		padding: 2rem;
		color: #eee;
	}

	.settings-section {
		background-color: white;
		border: 1px solid #ddd;
		border-radius: 8px;
		padding: 1.5rem 2rem;
		margin-bottom: 1.5rem;
		box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
	}

	.form-group {
		margin-bottom: 1.25rem;
	}

	.form-group:last-child {
		margin-bottom: 0;
	}

	.form-group label {
		display: block;
		font-weight: 500;
		color: #333;
		margin-bottom: 0.35rem;
		font-size: 0.95rem;
	}

	.input-hint {
		font-size: 0.82rem;
		color: #888;
		margin-bottom: 0.5rem;
	}

	.form-group input {
		width: 100%;
		padding: 0.6rem 0.75rem;
		border: 1px solid #ddd;
		border-radius: 4px;
		font-size: 1rem;
		font-family: inherit;
		box-sizing: border-box;
	}

	.form-group input:focus {
		outline: none;
		border-color: #0066cc;
		box-shadow: 0 0 0 3px rgba(0, 102, 204, 0.1);
	}

	.form-group input:disabled {
		background-color: #f5f5f5;
		cursor: not-allowed;
	}

	.form-actions {
		display: flex;
		justify-content: flex-end;
	}

	.btn {
		padding: 0.6rem 1.5rem;
		border-radius: 4px;
		border: none;
		cursor: pointer;
		font-weight: 500;
		font-size: 1rem;
		transition: background-color 0.2s;
	}

	.btn:disabled {
		opacity: 0.6;
		cursor: not-allowed;
	}

	.btn-primary {
		background-color: #0066cc;
		color: white;
	}

	.btn-primary:hover:not(:disabled) {
		background-color: #0052a3;
	}
</style>
