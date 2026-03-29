import { redirect } from '@sveltejs/kit';
import type { Actions } from './$types';

export const actions = {
	logout: async () => {
		// クライアント側で認証状態をクリアするため、ここでは単純にリダイレクト
		// APIへの logout 呼び出しはクライアント側で行う可能性がある
		throw redirect(303, '/');
	},
} satisfies Actions;
