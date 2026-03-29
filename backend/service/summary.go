package service

import (
	"database/sql"
	"purchase-system/model"
)

type SummaryService struct {
	db *sql.DB
}

func NewSummaryService(db *sql.DB) *SummaryService {
	return &SummaryService{db: db}
}

// GetUserBalance はユーザーの残高サマリを取得
func (s *SummaryService) GetUserBalance(userID int) (*model.BalanceSummary, error) {
	user := &model.User{}
	err := s.db.QueryRow(
		"SELECT id, name FROM users WHERE id = ?",
		userID,
	).Scan(&user.ID, &user.Name)
	if err != nil {
		return nil, err
	}

	// 購入未払い額 = SUM(purchase_items.quantity * unit_price) - SUM(payments.amount where deleted_at IS NULL)
	purchaseTotal := 0
	err = s.db.QueryRow(`
		SELECT COALESCE(SUM(pi.quantity * pi.unit_price), 0)
		FROM purchases pu
		JOIN purchase_items pi ON pi.purchase_id = pu.id
		WHERE pu.user_id = ?
	`, userID).Scan(&purchaseTotal)
	if err != nil {
		return nil, err
	}

	paymentTotal := 0
	err = s.db.QueryRow(`
		SELECT COALESCE(SUM(amount), 0)
		FROM payments
		WHERE user_id = ? AND deleted_at IS NULL
	`, userID).Scan(&paymentTotal)
	if err != nil {
		return nil, err
	}

	purchaseUnpaid := purchaseTotal - paymentTotal

	// 仕入れ未精算額 = SUM(restocks.total_amount where deleted_at IS NULL) - SUM(restock_payments.amount where deleted_at IS NULL)
	restockTotal := 0
	err = s.db.QueryRow(`
		SELECT COALESCE(SUM(total_amount), 0)
		FROM restocks
		WHERE user_id = ? AND deleted_at IS NULL
	`, userID).Scan(&restockTotal)
	if err != nil {
		return nil, err
	}

	restockPaymentTotal := 0
	err = s.db.QueryRow(`
		SELECT COALESCE(SUM(amount), 0)
		FROM restock_payments
		WHERE user_id = ? AND deleted_at IS NULL
	`, userID).Scan(&restockPaymentTotal)
	if err != nil {
		return nil, err
	}

	restockUnclaimed := restockTotal - restockPaymentTotal

	// 差引残高 = 購入未払い - 仕入れ未精算
	netBalance := purchaseUnpaid - restockUnclaimed

	return &model.BalanceSummary{
		UserID:           userID,
		UserName:         user.Name,
		PurchaseUnpaid:   purchaseUnpaid,
		RestockUnclaimed: restockUnclaimed,
		NetBalance:       netBalance,
	}, nil
}
