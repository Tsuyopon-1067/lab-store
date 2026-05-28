<script lang="ts">
  import type { PurchaseResponse } from "$lib/types";
  import ctrlBarcode from "$lib/assets/barcodes/ctrlBarcode.svg?url";

  interface Props {
    receipt: PurchaseResponse;
    onClose: () => void;
  }

  let { receipt, onClose }: Props = $props();

  const formatCurrency = (amount: number) =>
    `¥${amount.toLocaleString("ja-JP")}`;
  const formatDateTime = (dateStr: string) => {
    const date = new Date(dateStr);
    return date.toLocaleString("ja-JP");
  };
</script>

<div class="receipt-overlay">
  <div class="receipt-card">
    <h2>購入完了</h2>
    <div class="receipt-header">
      <div class="receipt-item">
        <span class="label">購入者</span>
        <span class="value">{receipt.user_name}</span>
      </div>
      <div class="receipt-item">
        <span class="label">購入日時</span>
        <span class="value">{formatDateTime(receipt.purchased_at)}</span>
      </div>
    </div>

    <div class="control-barcode-section">
      <p class="control-barcode-label">最初の画面に戻る</p>
      <img
        src={ctrlBarcode}
        alt="最初の画面に戻るためのバーコード"
        class="control-barcode-image"
      />
    </div>

    <div class="receipt-header">
      <div class="receipt-item">
        <span class="label">購入者</span>
        <span class="value">{receipt.user_name}</span>
      </div>
      <div class="receipt-item">
        <span class="label">購入日時</span>
        <span class="value">{formatDateTime(receipt.purchased_at)}</span>
      </div>
    </div>

    <div class="receipt-items">
      <h3>購入内容</h3>
      <table>
        <thead>
          <tr>
            <th>商品名</th>
            <th class="text-right">単価</th>
            <th class="text-right">数量</th>
            <th class="text-right">小計</th>
          </tr>
        </thead>
        <tbody>
          {#each receipt.items as item (item.product_id)}
            <tr>
              <td>{item.product_name}</td>
              <td class="text-right">{formatCurrency(item.unit_price)}</td>
              <td class="text-right">{item.quantity}</td>
              <td class="text-right">{formatCurrency(item.subtotal)}</td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>

    <div class="receipt-total">
      <span class="label">合計金額</span>
      <span class="total">{formatCurrency(receipt.total_amount)}</span>
    </div>

    <div class="receipt-balance">
      <h3>更新後の残高</h3>
      <div class="balance-item">
        <span class="label">購入未払い額</span>
        <span class="value">
          {formatCurrency(receipt.updated_balance.purchase_unpaid)}
        </span>
      </div>
      <div class="balance-item">
        <span class="label">仕入れ立替未精算額</span>
        <span class="value">
          {formatCurrency(receipt.updated_balance.restock_unclaimed)}
        </span>
      </div>
      <div class="balance-item net">
        <span class="label">
          {receipt.updated_balance.net_balance > 0
            ? "今月の支払い予定額"
            : "今月の受取予定額"}
        </span>
        <span
          class="value"
          class:negative={receipt.updated_balance.net_balance < 0}
        >
          {formatCurrency(Math.abs(receipt.updated_balance.net_balance))}
        </span>
      </div>
    </div>

    <button onclick={onClose} class="btn-close"> 最初の画面に戻る </button>
  </div>
</div>

<style>
  .receipt-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-color: rgba(0, 0, 0, 0.7);
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 1000;
    padding: 1rem;
  }

  .receipt-card {
    background-color: white;
    border-radius: 8px;
    padding: 2rem;
    max-width: 600px;
    width: 100%;
    max-height: 90vh;
    overflow-y: auto;
    box-shadow: 0 10px 40px rgba(0, 0, 0, 0.3);
  }

  h2 {
    text-align: center;
    margin: 0 0 2rem 0;
    font-size: 1.8rem;
    color: #2c3e50;
  }

  h3 {
    margin: 1.5rem 0 1rem 0;
    font-size: 1rem;
    color: #2c3e50;
    border-bottom: 2px solid #f0f0f0;
    padding-bottom: 0.5rem;
  }

  .receipt-header {
    background-color: #f9f9f9;
    border-radius: 4px;
    padding: 1rem;
    margin-bottom: 1.5rem;
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1rem;
  }

  .receipt-item {
    display: flex;
    justify-content: space-between;
    flex-direction: column;
  }

  .label {
    color: #666;
    font-size: 0.9rem;
    margin-bottom: 0.25rem;
  }

  .value {
    font-weight: 600;
    color: #2c3e50;
  }

  .receipt-items {
    margin-bottom: 1.5rem;
  }

  table {
    width: 100%;
    border-collapse: collapse;
    margin-top: 0.5rem;
  }

  thead {
    background-color: #f0f0f0;
  }

  th {
    padding: 0.5rem;
    text-align: left;
    font-weight: 600;
    font-size: 0.9rem;
    color: #666;
  }

  td {
    padding: 0.75rem 0.5rem;
    border-bottom: 1px solid #f0f0f0;
  }

  .text-right {
    text-align: right;
  }

  .receipt-total {
    background-color: #f0f8ff;
    padding: 1rem;
    border-radius: 4px;
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1.5rem;
  }

  .total {
    font-size: 1.5rem;
    font-weight: 700;
    color: #0066cc;
  }

  .receipt-balance {
    background-color: #f9f9f9;
    padding: 1rem;
    border-radius: 4px;
    margin-bottom: 2rem;
  }

  .receipt-balance .balance-item {
    display: flex;
    justify-content: space-between;
    margin-bottom: 0.75rem;
    padding: 0.5rem 0;
  }

  .balance-item.net {
    font-weight: 600;
    padding: 0.75rem;
    background-color: #e8f4f8;
    border-radius: 4px;
  }

  .value.negative {
    color: #f97316;
  }

  .btn-close {
    width: 100%;
    padding: 1rem;
    background-color: #0066cc;
    color: white;
    border: none;
    border-radius: 4px;
    font-size: 1.1rem;
    font-weight: 600;
    cursor: pointer;
    transition: background-color 0.2s;
  }

  .btn-close:hover {
    background-color: #0052a3;
  }

  .control-barcode-section {
    margin-top: 2rem;
    padding-top: 1.5rem;
    border-top: 1px solid #e0e0e0;
    text-align: center;
  }

  .control-barcode-label {
    margin: 0 0 1rem 0;
    color: #666;
    font-size: 0.9rem;
    font-weight: 600;
  }

  .control-barcode-image {
    width: 100%;
    height: auto;
    max-width: 250px;
  }

  @media (max-width: 600px) {
    .receipt-card {
      padding: 1rem;
    }

    .receipt-header {
      grid-template-columns: 1fr;
    }

    table {
      font-size: 0.9rem;
    }

    th,
    td {
      padding: 0.25rem;
    }
  }
</style>
