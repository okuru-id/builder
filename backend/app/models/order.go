package models

import (
	"time"

	"github.com/goravel/framework/database/orm"
)

type Order struct {
	orm.Model
	Reference     string     `gorm:"uniqueIndex" json:"reference"`
	ProductID     uint       `json:"product_id"`
	BuyerEmail    string     `json:"buyer_email"`
	BuyerName     string     `json:"buyer_name"`
	Amount        int        `json:"amount"`
	Status        string     `json:"status"`
	PaymentRef    *string    `json:"payment_ref"`
	DownloadToken *string    `json:"download_token"`
	DownloadCount int        `json:"download_count"`
	MaxDownloads  int        `json:"max_downloads"`
	ExpiresAt     *time.Time `json:"expires_at"`
	PaidAt        *time.Time `json:"paid_at"`
}
