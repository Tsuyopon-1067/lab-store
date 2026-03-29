<script lang="ts">
	import { auth } from '$lib/stores/auth';
	import { goto } from '$app/navigation';

	let { children } = $props();

	// 管理画面にアクセス時、未認証ならログイン画面へリダイレクト
	$effect(() => {
		if (!$auth.isLoggedIn && !window.location.pathname.includes('/admin/login')) {
			goto('/admin/login');
		}
	});
</script>

<div class="admin-layout">
	<aside class="admin-sidebar">
		<nav class="admin-menu">
			<h3>管理メニュー</h3>
			<ul>
				<li><a href="/admin">ダッシュボード</a></li>
				<li><a href="/admin/users">利用者管理</a></li>
				<li><a href="/admin/products">商品管理</a></li>
				<li><a href="/admin/purchases">購入履歴</a></li>
				<li><a href="/admin/restocks">仕入れ履歴</a></li>
				<li><a href="/admin/payments">支払い管理</a></li>
				<li><a href="/admin/restock-payments">立替精算</a></li>
				<li><a href="/admin/summary">月次精算</a></li>
				<li><a href="/admin/backup">バックアップ</a></li>
				<li><a href="/admin/settings">設定</a></li>
			</ul>
		</nav>
	</aside>

	<div class="admin-content">
		{@render children()}
	</div>
</div>

<style>
	.admin-layout {
		display: flex;
		gap: 0;
		flex: 1;
	}

	.admin-sidebar {
		width: 200px;
		background-color: #2c3e50;
		color: white;
		padding: 1rem 0;
		border-right: 1px solid #34495e;
		min-height: calc(100vh - 60px);
	}

	.admin-menu h3 {
		padding: 0 1rem;
		font-size: 0.9rem;
		color: #95a5a6;
		text-transform: uppercase;
		margin-bottom: 1rem;
	}

	.admin-menu ul {
		list-style: none;
	}

	.admin-menu a {
		display: block;
		padding: 0.75rem 1rem;
		color: #ecf0f1;
		text-decoration: none;
		transition: background-color 0.2s;
	}

	.admin-menu a:hover {
		background-color: #34495e;
		text-decoration: none;
	}

	.admin-content {
		flex: 1;
		padding: 2rem;
		background-color: #f5f5f5;
		overflow-y: auto;
	}

	@media (max-width: 768px) {
		.admin-layout {
			flex-direction: column;
		}

		.admin-sidebar {
			width: 100%;
			min-height: auto;
			border-right: none;
			border-bottom: 1px solid #34495e;
		}

		.admin-menu ul {
			display: flex;
			flex-wrap: wrap;
		}

		.admin-menu a {
			flex: 1;
		}
	}
</style>
