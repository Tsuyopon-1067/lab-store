// See https://svelte.dev/docs/kit/types#app.d.ts
// for information about these interfaces

import type { SvelteWindowAttributes } from 'svelte/elements';

declare global {
	namespace App {
		// interface Error {}
		// interface Locals {}
		// interface PageData {}
		// interface PageState {}
		// interface Platform {}
	}

	namespace svelteHTML {
		interface HTMLAttributes<T> {
			'on:scan'?: (event: CustomEvent<{ barcode: string }>) => void;
		}
	}
}

export {};
