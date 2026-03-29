import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
  plugins: [sveltekit()],
  server: {
    port: 3000,
    proxy: {
      // 開発時: /api/* をバックエンド（go run）に転送
      '/api': 'http://localhost:8080',
    },
  },
});
