package model

type SummaryItem struct {
	UserID             int    `json:"user_id"`
	UserName           string `json:"user_name"`
	PurchaseTotal      int    `json:"purchase_total"`
	PurchasePaid       int    `json:"purchase_paid"`
	PurchaseUnpaid     int    `json:"purchase_unpaid"`
	RestockTotal       int    `json:"restock_total"`
	RestockSettled     int    `json:"restock_settled"`
	RestockUnclaimed   int    `json:"restock_unclaimed"`
	NetBalance         int    `json:"net_balance"`
}

type SummaryReport struct {
	Period struct {
		From string `json:"from"`
		To   string `json:"to"`
	} `json:"period"`
	Users []*SummaryItem `json:"users"`
}

type GetSummaryRequest struct {
	From string `form:"from"`
	To   string `form:"to"`
}
