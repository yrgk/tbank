package repository

import (
	"payment-microservice/internal/models"
	"payment-microservice/pkg/postgres"
)

func CreateUser(body models.CreateAccountPayload) error {
	user := models.TbankAccount{
		TerminalKey: body.TerminalKey,
		InviteToken: body.InviteToken,
		Password:    body.Password,
	}

	if err := postgres.DB.Create(&user).Error; err != nil {
		return err
	}

	return nil
}

func GetPasswordByToken(userToken string) string {
	var password string
	postgres.DB.Model(&models.TbankAccount{}).Where("user_token = ?", userToken).Find(&password)

	return password
}

func GetUser(inviteToken string) models.TbankAccount {
	var user models.TbankAccount

	postgres.DB.Model(models.TbankAccount{}).Where("invite_token = ?", inviteToken).First(&user)

	return user
}
