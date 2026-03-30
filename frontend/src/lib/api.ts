import { auth } from './stores/auth';

const API_BASE = '/api';

interface FetchOptions extends RequestInit {
	headers?: Record<string, string>;
}

/**
 * API呼び出しのラッパー関数
 */
export async function apiCall<T>(
	endpoint: string,
	options: FetchOptions = {},
): Promise<T> {
	const url = `${API_BASE}${endpoint}`;

	const headers: Record<string, string> = {
		'Content-Type': 'application/json',
		...options.headers,
	};

	const response = await fetch(url, {
		...options,
		headers,
	});

	if (!response.ok) {
		let errorMsg = `API error: ${response.status}`;
		try {
			const error = await response.json();
			errorMsg = error.message || error.error || errorMsg;
		} catch {
			// JSON parse failed, try to get text
			try {
				const text = await response.text();
				if (text) {
					errorMsg = text.slice(0, 200); // Limit length
				}
			} catch {
				// Ignore text parsing error
			}
		}
		throw new Error(errorMsg);
	}

	return response.json();
}

/**
 * 管理者認証が必要なAPI呼び出し
 */
export async function apiCallWithAuth<T>(
	endpoint: string,
	options: FetchOptions = {},
): Promise<T> {
	let adminToken: string | null = null;

	// ストアから現在のトークンを取得
	const unsubscribe = auth.subscribe((state) => {
		adminToken = state.adminToken;
	});
	unsubscribe();

	if (!adminToken) {
		throw new Error('認証が必要です');
	}

	return apiCall<T>(endpoint, {
		...options,
		headers: {
			...(options.headers || {}),
			'Authorization': `Bearer ${adminToken}`,
		},
	});
}

/**
 * バーコード認証が必要なAPI呼び出し
 */
export async function apiCallWithBarcode<T>(
	endpoint: string,
	barcode: string,
	options: FetchOptions = {},
): Promise<T> {
	return apiCall<T>(endpoint, {
		...options,
		headers: {
			...(options.headers || {}),
			'X-User-Barcode': barcode,
		},
	});
}
