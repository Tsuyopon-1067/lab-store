/**
 * API型定義
 */

export interface User {
	id: number;
	name: string;
	barcode: string;
	is_active: boolean;
	created_at: string;
}

export interface Product {
	id: number;
	name: string;
	barcode: string;
	is_active: boolean;
	current_price: number;
	created_at: string;
	note?: string;
}

export interface CartItem {
	product_id: number;
	product_name: string;
	quantity: number;
	unit_price: number;
	subtotal: number;
}

export interface UserBalance {
	user_id: number;
	user_name: string;
	purchase_unpaid: number;
	restock_unclaimed: number;
	net_balance: number;
}

export interface PurchaseRequest {
	user_id: number;
	items: Array<{
		product_id: number;
		quantity: number;
		unit_price: number;
	}>;
}

export interface PurchaseResponse {
	id: number;
	user_id: number;
	user_name: string;
	total_amount: number;
	items: CartItem[];
	purchased_at: string;
	updated_balance: UserBalance;
}

export interface RestockItem {
	product_id: number;
	product_name: string;
	quantity: number;
	unit_price: number;
	subtotal: number;
}

export interface RestockRequest {
	user_id: number;
	restocked_at: string;
	note?: string;
	items: Array<{
		product_id: number;
		quantity: number;
		unit_price: number;
	}>;
}

export interface RestockResponse {
	id: number;
	user_id: number;
	user_name: string;
	total_amount: number;
	items: RestockItem[];
	restocked_at: string;
	note?: string;
}

export interface PurchaseHistory {
	id: number;
	purchased_at: string;
	total_amount: number;
	items: CartItem[];
}

export interface RestockHistory {
	id: number;
	total_amount: number;
	restocked_at: string;
	note?: string;
	items: RestockItem[];
}
