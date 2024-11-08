package models

import "gorm.io/gorm"

type (
	TbankAccount struct {
		gorm.Model
		TerminalKey string
		InviteToken string
		Password    string
	}

	TbankOrder struct {
		gorm.Model
		Amount         int
		Token          string
		Description    string
		InviteToken    string
		PaymentId      int
		DocsSalesId    int
		TbankPaymentId string
		PaymentURL     string
	}
)
