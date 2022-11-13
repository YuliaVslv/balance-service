package model

import (
	"time"
)

type Account struct {
	ID      uint32  `json:"user_id"`
	Balance float32 `json:"balance"`
}

type Crediting struct {
	ID    uint32  `json:"user_id"`
	Value float32 `json:"value"`
}

type Reserve struct {
	UserID    uint32  `json:"user_id"`
	ServiceID uint32  `json:"service_id"`
	OrderID   uint32  `json:"order_id"`
	Value     float32 `json:"value"`
}

type Transaction struct {
	Date        time.Time `json:"date"`
	Value       float32   `json:"value"`
	Description string    `json:"description"`
}

type JSONMessage struct {
	Message string `json:"message"`
}
