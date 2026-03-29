import { writable } from 'svelte/store';

export interface AuthState {
	isLoggedIn: boolean;
	adminToken: string | null;
	userBarcode: string | null;
}

const initialState: AuthState = {
	isLoggedIn: false,
	adminToken: null,
	userBarcode: null,
};

function createAuthStore() {
	const { subscribe, set, update } = writable<AuthState>(initialState);

	return {
		subscribe,
		login: (token: string) => {
			update((state) => ({
				...state,
				isLoggedIn: true,
				adminToken: token,
			}));
		},
		logout: () => {
			set(initialState);
		},
		setUserBarcode: (barcode: string) => {
			update((state) => ({
				...state,
				userBarcode: barcode,
			}));
		},
		clearUserBarcode: () => {
			update((state) => ({
				...state,
				userBarcode: null,
			}));
		},
	};
}

export const auth = createAuthStore();
