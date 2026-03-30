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
	product_name?: string;
	quantity: number;
	unit_price: number;
	subtotal?: number;
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
	restock_id: number;
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

// --- Admin用型 ---

export interface AdminUser {
	id: number;
	name: string;
	barcode: string;
	is_active: number; // 0 or 1
	created_at: string;
}

export interface ProductWithPrice {
	id: number;
	name: string;
	barcode: string;
	is_active: number;
	note?: string;
	price: number;
	created_at: string;
}

export interface ProductPrice {
	id: number;
	product_id: number;
	price: number;
	valid_from: string;
	valid_to: string | null;
}

export interface AdminPurchase {
	id: number;
	user_id: number;
	purchased_at: string;
}

export interface AdminRestock {
	id: number;
	user_id: number;
	total_amount: number;
	restocked_at: string;
	note?: string;
	deleted_at?: string;
	created_at: string;
}

export interface Payment {
	id: number;
	user_id: number;
	amount: number;
	paid_at: string;
	note?: string;
	created_at: string;
}

export interface PaymentResponse {
	id: number;
	user_id: number;
	user_name: string;
	amount: number;
	paid_at: string;
	note?: string;
}

export interface RestockPayment {
	id: number;
	user_id: number;
	amount: number;
	settled_at: string;
	note?: string;
	created_at: string;
}

export interface RestockPaymentResponse {
	id: number;
	user_id: number;
	user_name: string;
	amount: number;
	settled_at: string;
	note?: string;
}

export interface SummaryUser {
	user_id: number;
	user_name: string;
	purchase_total: number;
	purchase_paid: number;
	purchase_unpaid: number;
	restock_total: number;
	restock_settled: number;
	restock_unclaimed: number;
	net_balance: number;
}

export interface SummaryReport {
	period: { from: string; to: string };
	users: SummaryUser[];
}
