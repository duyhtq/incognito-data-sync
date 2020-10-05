package models

type PToken struct {
	ID      int     `db:"id"`
	TokenID string  `db:"token_id"`
	Name    string  `db:"name"`
	Symbol  string  `db:"symbol"`
	Decimal float64 `db:"decimal"`
	Price   float64 `db:"price"`
}
