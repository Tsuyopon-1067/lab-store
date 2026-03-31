<script lang="ts">
    import type { CartItem } from "$lib/types";

    interface Props {
        items: CartItem[];
        onQuantityChange: (productId: number, quantity: number) => void;
        onRemove: (productId: number) => void;
    }

    let { items, onQuantityChange, onRemove }: Props = $props();

    const formatCurrency = (amount: number) =>
        `¥${amount.toLocaleString("ja-JP")}`;

    const total = items.reduce((sum, item) => sum + item.subtotal, 0);
</script>

<div class="cart">
    <h3>カート内容</h3>

    {#if items.length === 0}
        <div class="empty-cart">
            <p>商品がまだカートに追加されていません</p>
        </div>
    {:else}
        <div class="cart-items">
            {#each items as item (item.product_id)}
                <div class="cart-item">
                    <div class="item-info">
                        <div class="item-name">{item.product_name}</div>
                        <div class="item-price">
                            {formatCurrency(item.unit_price)}
                        </div>
                    </div>

                    <div class="item-controls">
                        <button
                            onclick={() =>
                                onQuantityChange(
                                    item.product_id,
                                    item.quantity - 1,
                                )}
                            class="qty-btn"
                            disabled={item.quantity <= 1}
                        >
                            −
                        </button>
                        <input
                            type="number"
                            value={item.quantity}
                            onchange={(e) => {
                                const qty =
                                    parseInt(
                                        (e.target as HTMLInputElement).value,
                                    ) || 1;
                                onQuantityChange(
                                    item.product_id,
                                    Math.max(1, qty),
                                );
                            }}
                            class="qty-input"
                            min="1"
                        />
                        <button
                            onclick={() =>
                                onQuantityChange(
                                    item.product_id,
                                    item.quantity + 1,
                                )}
                            class="qty-btn"
                        >
                            +
                        </button>
                    </div>

                    <div class="item-subtotal">
                        {formatCurrency(item.subtotal)}
                    </div>

                    <button
                        onclick={() => onRemove(item.product_id)}
                        class="btn-remove"
                        title="削除"
                    >
                        ×
                    </button>
                </div>
            {/each}
        </div>

        <div class="cart-total">
            <span class="label">合計金額</span>
            <span class="total-amount">{formatCurrency(total)}</span>
        </div>
    {/if}
</div>

<style>
    .cart {
        background-color: white;
        border: 1px solid #ddd;
        border-radius: 8px;
        padding: 1.5rem;
        margin-bottom: 2rem;
    }

    h3 {
        margin: 0 0 1rem 0;
        font-size: 1.2rem;
        color: #2c3e50;
    }

    .empty-cart {
        text-align: center;
        padding: 2rem;
        color: #999;
    }

    .cart-items {
        border-bottom: 2px solid #f0f0f0;
        padding-bottom: 1rem;
        margin-bottom: 1rem;
    }

    .cart-item {
        display: grid;
        grid-template-columns: 1fr auto auto auto;
        align-items: center;
        gap: 1rem;
        padding: 1rem;
        border-bottom: 1px solid #f5f5f5;
    }

    .cart-item:last-child {
        border-bottom: none;
    }

    .item-info {
        display: flex;
        flex-direction: column;
    }

    .item-name {
        font-weight: 600;
        color: #2c3e50;
    }

    .item-price {
        font-size: 0.9rem;
        color: #666;
    }

    .item-controls {
        display: flex;
        gap: 0.5rem;
        align-items: center;
    }

    .qty-btn {
        background-color: #f0f0f0;
        border: 1px solid #ddd;
        width: 28px;
        height: 28px;
        border-radius: 4px;
        cursor: pointer;
        font-weight: 600;
        transition: background-color 0.2s;
        padding: 0.05rem 0.2rem 0.2rem;
    }

    .qty-btn:hover:not(:disabled) {
        background-color: #e0e0e0;
    }

    .qty-btn:disabled {
        opacity: 0.5;
        cursor: not-allowed;
    }

    .qty-input {
        width: 50px;
        text-align: center;
        border: 1px solid #ddd;
        border-radius: 4px;
        padding: 0.25rem;
    }

    .item-subtotal {
        text-align: right;
        font-weight: 600;
        color: #0066cc;
        min-width: 100px;
    }

    .btn-remove {
        background-color: #fee;
        color: #c00;
        border: 1px solid #fcc;
        width: 32px;
        height: 32px;
        border-radius: 4px;
        cursor: pointer;
        font-size: 1.2rem;
        transition: background-color 0.2s;
        padding: 0rem 0.2rem 0.2rem;
    }

    .btn-remove:hover {
        background-color: #fdd;
    }

    .cart-total {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 1rem;
        background-color: #f0f8ff;
        border-radius: 4px;
    }

    .label {
        font-weight: 600;
        color: #2c3e50;
    }

    .total-amount {
        font-size: 1.5rem;
        font-weight: 700;
        color: #0066cc;
    }

    @media (max-width: 768px) {
        .cart-item {
            grid-template-columns: 1fr auto;
        }

        .item-controls,
        .item-subtotal {
            grid-column: 1 / -1;
        }
    }
</style>
