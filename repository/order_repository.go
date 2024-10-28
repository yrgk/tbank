package repository

import (
	"payment-microservice/internal/models"
	"payment-microservice/pkg/postgres"
	"strconv"
)

func CreateOrder(body models.TbankOrder) error {
	if err := postgres.DB.Create(&body).Error; err != nil {
		return err
	}

	return nil
}

func GetDataForPayment(paymentId int, inviteToken string) models.CreateTokenRequest {
	var data models.CreateTokenRequest

	err := postgres.DB.Table("tbank_orders").
		Select("tbank_orders.id, tbank_accounts.terminal_key, tbank_orders.amount, tbank_orders.description, tbank_orders.amount, tbank_accounts.password").
		Joins("JOIN tbank_accounts ON tbank_orders.invite_token = tbank_accounts.invite_token").
		Where(&models.TbankOrder{PaymentId: paymentId, InviteToken: inviteToken}).
		Last(&data).Error

	if err != nil {
		return data
	}

	return data
}

func GetPaymentUrlByPaymentId(paymentId int, inviteToken string) string {
	var order models.TbankOrder
	postgres.DB.Where("payment_id = ?", paymentId).Where("invite_token = ?", inviteToken).First(&order)

	return order.PaymentURL
}

func GetPaymentByTbankPaymentId(paymentId string) models.TbankOrder {
	pId, _ := strconv.Atoi(paymentId)

	var data = models.TbankOrder{PaymentId: pId}

	postgres.DB.First(&data)

	return data
}

func PaymentStatusDone(tbankPaymentId int) error {
	var order models.TbankOrder
	postgres.DB.Where("tbank_payment_id = ?", tbankPaymentId).First(&order)
	if err := postgres.DB.Exec("UPDATE payments SET status = true WHERE id = ?", order.PaymentId).Error; err != nil {
		return err
	}

	return nil
}