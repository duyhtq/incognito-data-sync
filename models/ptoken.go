package models

type PToken struct {
	ID      int     `db:"id"`
	TokenID string  `db:"token_id"`
	Name    string  `db:"name"`
	Symbol  string  `db:"symbol"`
	Decimal int     `db:"decimal"`
	Price   float64 `db:"price"`
}
