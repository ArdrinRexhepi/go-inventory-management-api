package models

import "time"

type RestockHistory struct {
	ID           int       `json:"id"`
	ItemID       int       `json:"item_id"`
	RestockAmount int      `json:"restock_amount"`
	RestockedAt  time.Time `json:"restocked_at"`
}
