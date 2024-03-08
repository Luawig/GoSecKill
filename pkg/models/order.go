package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	UserId    uint `gorm:"not null" json:"user_id"`
	ProductId uint `gorm:"not null" json:"product_id"`
	Num       int  `gorm:"not null" json:"num"`
	Total     int  `gorm:"not null" json:"total"`
	Status    int  `gorm:"not null" json:"status"`
}

const (
	OrderStatusWait = iota
	OrderStatusSuccess
	OrderStatusFailed
)
