<script lang="ts">
	import { auth } from '$lib/stores/auth';
</script>

<nav class="navbar">
	<div class="nav-container">
		<div class="nav-brand">
			<a href="/">Lab Store</a>
		</div>

		<ul class="nav-menu">
			<li>
				<a href="/">購買画面</a>
			</li>
			<li>
				<a href="/restock">仕入れ</a>
			</li>
			<li>
				<a href="/balance">残高確認</a>
			</li>

			{#if $auth.isLoggedIn}
				<li class="admin-nav">
					<a href="/admin">管理画面</a>
				</li>
				<li>
					<button
						class="logout-btn"
						onclick={() => {
							auth.logout();
						}}
					>
						ログアウト
					</button>
				</li>
			{:else}
				<li>
					<a href="/admin/login">管理者ログイン</a>
				</li>
			{/if}
		</ul>
	</div>
</nav>

<style>
	.navbar {
		background-color: #fff;
		border-bottom: 1px solid #ddd;
		box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
		position: sticky;
		top: 0;
		z-index: 100;
	}

	.nav-container {
		max-width: 1200px;
		margin: 0 auto;
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 0 1rem;
		height: 60px;
	}

	.nav-brand a {
		font-size: 1.5rem;
		font-weight: bold;
		color: #0066cc;
		text-decoration: none;
	}

	.nav-brand a:hover {
		text-decoration: none;
	}

	.nav-menu {
		display: flex;
		list-style: none;
		gap: 2rem;
		align-items: center;
	}

	.nav-menu a {
		color: #333;
		text-decoration: none;
		padding: 0.5rem 0;
		border-bottom: 2px solid transparent;
		transition: border-bottom-color 0.2s;
	}

	.nav-menu a:hover {
		border-bottom-color: #0066cc;
		text-decoration: none;
	}

	.admin-nav a {
		color: #f97316;
		font-weight: 500;
	}

	.logout-btn {
		background: none;
		border: none;
		color: #666;
		cursor: pointer;
		padding: 0;
		text-decoration: underline;
		font-size: 1rem;
	}

	.logout-btn:hover {
		color: #333;
	}

	@media (max-width: 768px) {
		.nav-menu {
			gap: 1rem;
		}

		.nav-container {
			flex-direction: column;
			height: auto;
			padding: 0.5rem 1rem;
		}
	}
</style>
