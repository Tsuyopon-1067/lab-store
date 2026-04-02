package repository

import (
	"database/sql"
	"purchase-system/model"
	"time"
)

type SummaryRepository struct {
	db *sql.DB
}

func NewSummaryRepository(db *sql.DB) *SummaryRepository {
	return &SummaryRepository{db: db}
}

// GetSummary は期間指定での月次精算サマリを取得する
func (r *SummaryRepository) GetSummary(from, to time.Time) (*model.SummaryReport, error) {
	rows, err := r.db.Query(`
		SELECT
			u.id,
			u.name,
			COALESCE(SUM(pi.quantity * pi.unit_price), 0) AS purchase_total,
			COALESCE(pay.paid, 0) AS purchase_paid,
			COALESCE(SUM(pi.quantity * pi.unit_price), 0) - COALESCE(pay.paid, 0) AS purchase_unpaid,
			COALESCE(rs.restock_total, 0) AS restock_total,
			COALESCE(rp.settled, 0) AS restock_settled,
			COALESCE(rs.restock_total, 0) - COALESCE(rp.settled, 0) AS restock_unclaimed
		FROM users u
		LEFT JOIN purchases pu ON pu.user_id = u.id AND pu.purchased_at BETWEEN ? AND ?
		LEFT JOIN purchase_items pi ON pi.purchase_id = pu.id
		LEFT JOIN (
			SELECT user_id, SUM(amount) AS paid
			FROM payments
			WHERE paid_at BETWEEN ? AND ? AND deleted_at IS NULL
			GROUP BY user_id
		) pay ON pay.user_id = u.id
		LEFT JOIN (
			SELECT user_id, SUM(total_amount) AS restock_total
			FROM restocks
			WHERE restocked_at BETWEEN ? AND ? AND deleted_at IS NULL
			GROUP BY user_id
		) rs ON rs.user_id = u.id
		LEFT JOIN (
			SELECT user_id, SUM(amount) AS settled
			FROM restock_payments
			WHERE settled_at BETWEEN ? AND ? AND deleted_at IS NULL
			GROUP BY user_id
		) rp ON rp.user_id = u.id
		WHERE u.is_active = 1
		GROUP BY u.id
		ORDER BY u.id
	`, from, to, from, to, from, to, from, to)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*model.SummaryItem
	for rows.Next() {
		item := &model.SummaryItem{}
		if err := rows.Scan(
			&item.UserID, &item.UserName,
			&item.PurchaseTotal, &item.PurchasePaid, &item.PurchaseUnpaid,
			&item.RestockTotal, &item.RestockSettled, &item.RestockUnclaimed,
		); err != nil {
			return nil, err
		}

		// 差引残高 = 購入未払い - 仕入れ未精算
		item.NetBalance = item.PurchaseUnpaid - item.RestockUnclaimed

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	report := &model.SummaryReport{
		Users: items,
	}
	report.Period.From = from.Format("2006-01-02")
	report.Period.To = to.Format("2006-01-02")

	return report, nil
}

// GetAllUsersBalance は全有効利用者の支払い金額サマリーを取得する（期間制限なし）
func (r *SummaryRepository) GetAllUsersBalance() ([]*model.BalanceSummary, error) {
	rows, err := r.db.Query(`
		SELECT
			u.id,
			u.name,
			COALESCE(SUM(pi.quantity * pi.unit_price), 0) AS purchase_total,
			COALESCE(pay.paid, 0) AS purchase_paid,
			COALESCE(rs.restock_total, 0) AS restock_total,
			COALESCE(rp.settled, 0) AS restock_settled
		FROM users u
		LEFT JOIN purchases pu ON pu.user_id = u.id AND pu.deleted_at IS NULL
		LEFT JOIN purchase_items pi ON pi.purchase_id = pu.id
		LEFT JOIN (
			SELECT user_id, SUM(amount) AS paid
			FROM payments
			WHERE deleted_at IS NULL
			GROUP BY user_id
		) pay ON pay.user_id = u.id
		LEFT JOIN (
			SELECT user_id, SUM(total_amount) AS restock_total
			FROM restocks
			WHERE deleted_at IS NULL
			GROUP BY user_id
		) rs ON rs.user_id = u.id
		LEFT JOIN (
			SELECT user_id, SUM(amount) AS settled
			FROM restock_payments
			WHERE deleted_at IS NULL
			GROUP BY user_id
		) rp ON rp.user_id = u.id
		WHERE u.is_active = 1
		GROUP BY u.id
		ORDER BY u.id
	`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var balances []*model.BalanceSummary
	for rows.Next() {
		var userID int
		var userName string
		var purchaseTotal, purchasePaid, restockTotal, restockSettled int

		if err := rows.Scan(
			&userID, &userName,
			&purchaseTotal, &purchasePaid,
			&restockTotal, &restockSettled,
		); err != nil {
			return nil, err
		}

		balance := &model.BalanceSummary{
			UserID:           userID,
			UserName:         userName,
			PurchaseUnpaid:   purchaseTotal - purchasePaid,
			RestockUnclaimed: restockTotal - restockSettled,
		}
		balance.NetBalance = balance.PurchaseUnpaid - balance.RestockUnclaimed

		balances = append(balances, balance)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return balances, nil
}
