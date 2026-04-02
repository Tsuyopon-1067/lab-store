<script lang="ts">
	import { goto } from '$app/navigation';
	import { auth } from '$lib/stores/auth';
	import { apiCall } from '$lib/api';

	let password = $state('');
	let isLoading = $state(false);
	let error = $state('');

	async function handleLogin(e: SubmitEvent) {
		e.preventDefault();
		isLoading = true;
		error = '';

		try {
			const response = await apiCall<{ token: string }>('/auth/login', {
				method: 'POST',
				body: JSON.stringify({ password }),
			});

			auth.login(response.token);
			goto('/admin');
		} catch (err) {
			error = err instanceof Error ? err.message : 'ログインに失敗しました';
		} finally {
			isLoading = false;
		}
	}
</script>

<div class="login-container">
	<div class="login-card">
		<h1>Lab Store 管理者ログイン</h1>

		<form on:submit={handleLogin}>
			<div class="form-group">
				<label for="password">パスワード</label>
				<input
					type="password"
					id="password"
					bind:value={password}
					placeholder="管理者パスワードを入力"
					required
					disabled={isLoading}
					autofocus
				/>
			</div>

			{#if error}
				<div class="error-message">
					<strong>エラー:</strong> {error}
				</div>
			{/if}

			<button type="submit" disabled={isLoading} class="btn btn-primary btn-block">
				{isLoading ? 'ログイン中...' : 'ログイン'}
			</button>
		</form>

		<p class="back-link">
			<a href="/">← 購買画面に戻る</a>
		</p>
	</div>
</div>

<style>
	.login-container {
		display: flex;
		justify-content: center;
		align-items: center;
		min-height: 100vh;
		background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
		padding: 1rem;
	}

	.login-card {
		background-color: white;
		border-radius: 8px;
		box-shadow: 0 10px 40px rgba(0, 0, 0, 0.2);
		padding: 2rem;
		width: 100%;
		max-width: 400px;
	}

	h1 {
		text-align: center;
		margin-bottom: 2rem;
		font-size: 1.5rem;
		color: #2c3e50;
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

	input {
		width: 100%;
		padding: 0.75rem;
		border: 1px solid #ccc;
		border-radius: 4px;
		font-size: 1rem;
		transition: border-color 0.2s;
	}

	input:focus {
		outline: none;
		border-color: #667eea;
		box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
	}

	input:disabled {
		background-color: #f5f5f5;
		cursor: not-allowed;
	}

	.error-message {
		background-color: #fee;
		color: #c00;
		padding: 0.75rem;
		border-radius: 4px;
		margin-bottom: 1rem;
		font-size: 0.95rem;
	}

	.btn {
		padding: 0.75rem;
		border: none;
		border-radius: 4px;
		font-size: 1rem;
		font-weight: 500;
		cursor: pointer;
		transition: background-color 0.2s;
	}

	.btn:disabled {
		opacity: 0.6;
		cursor: not-allowed;
	}

	.btn-primary {
		background-color: #667eea;
		color: white;
	}

	.btn-primary:hover:not(:disabled) {
		background-color: #5568d3;
	}

	.btn-block {
		width: 100%;
	}

	.back-link {
		text-align: center;
		margin-top: 1.5rem;
		color: #666;
		font-size: 0.95rem;
	}

	.back-link a {
		color: #667eea;
		text-decoration: none;
	}

	.back-link a:hover {
		text-decoration: underline;
	}
</style>
