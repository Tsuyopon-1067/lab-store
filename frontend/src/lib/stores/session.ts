import { writable } from 'svelte/store';

function createSessionStore() {
	const { subscribe, set } = writable<number>(0);

	return {
		subscribe,
		reset: () => {
			// Trigger reset by incrementing a counter
			set(Date.now());
		},
	};
}

export const session = createSessionStore();
