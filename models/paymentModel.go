package models

import (
	"time"

	"gorm.io/gorm"
)

type Payment struct {
	gorm.Model
	InstitutionID   uint
	TournamentID    uint
	CategoryID      uint
	UserID          uint
	PlayerID        uint
	PaymentType     string  `gorm:"type:varchar(50);not null"` // "torneo","liga", "membresia"
	Amount          float64 `gorm:"type:decimal(10,2);not null"`
	Currency        string  `gorm:"type:char(3);not null"`     // "MEX"
	PaymentMethod   string  `gorm:"type:varchar(50);not null"` // "card", "cash", "transfer"
	PaymentVendorID string  `gorm:"type:varchar(255)"`         // "stripe", "paypal", "mercadopago"
	TransactionID   string  `gorm:"type:varchar(255);unique"`  // transaction id from payment vendor
	Description     string  `gorm:"type:text"`                 // description of the payment i.e. "Torneo en Smash padel", "Membresia 2025 Circuito Moerlia"
	Status          string  `gorm:"type:varchar(20);not null"` // "pending", "processing", "paid", "cancel", "refunded"
	PaymentDate     *time.Time
	RefundStatus    string `gorm:"type:varchar(20)"`
	RefundDate      *time.Time
	MetaData        string `gorm:"type:json"` // additional information about the payment
}
