<script lang="ts">
	import favicon from '$lib/assets/favicon.svg';
	import Navigation from '$lib/components/Navigation.svelte';

	let { children } = $props();
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<div class="app">
	<Navigation />
	<main class="main-content">
		{@render children()}
	</main>
</div>

<style>
	:global(*) {
		margin: 0;
		padding: 0;
		box-sizing: border-box;
	}

	:global(html, body) {
		height: 100%;
		font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial,
			'Noto Sans', sans-serif, 'Apple Color Emoji', 'Segoe UI Emoji', 'Segoe UI Symbol';
		font-size: 16px;
		line-height: 1.5;
		color: #333;
		background-color: #f5f5f5;
	}

	:global(#app) {
		display: flex;
		flex-direction: column;
		min-height: 100vh;
	}

	:global(a) {
		color: #0066cc;
		text-decoration: none;
	}

	:global(a:hover) {
		text-decoration: underline;
	}

	:global(button) {
		cursor: pointer;
		border: none;
		border-radius: 4px;
		padding: 0.5rem 1rem;
		font-size: 1rem;
		font-weight: 500;
		transition: background-color 0.2s;
	}

	:global(input, textarea, select) {
		padding: 0.5rem;
		border: 1px solid #ccc;
		border-radius: 4px;
		font-size: 1rem;
	}

	:global(input:focus, textarea:focus, select:focus) {
		outline: none;
		border-color: #0066cc;
		box-shadow: 0 0 0 3px rgba(0, 102, 204, 0.1);
	}

	.app {
		display: flex;
		flex-direction: column;
		flex: 1;
	}

	.main-content {
		flex: 1;
		display: flex;
		flex-direction: column;
	}
</style>
