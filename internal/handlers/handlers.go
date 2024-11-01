package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"fmt"
	"payment-microservice/config"
	"payment-microservice/internal/models"
	"payment-microservice/internal/utils"
	"payment-microservice/pkg/postgres"
	"payment-microservice/repository"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func RegisterTbankHandler(c *fiber.Ctx) error {
	body := new(models.CreateAccountRequest)

	if err := c.BodyParser(body); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	inviteToken, err := utils.GetIviteTokenByUserToken(body.UserToken)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	payload := models.CreateAccountPayload{
		TerminalKey: body.TerminalKey,
		InviteToken: inviteToken,
		Password:    body.Password,
	}

	if err := repository.CreateUser(payload); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	return c.SendStatus(fiber.StatusOK)
}

func CreatePaymentHandler(c *fiber.Ctx) error {
	body := new(models.CreatePaymentRequest)

	if err := c.BodyParser(body); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	inviteToken, err := utils.GetIviteTokenByUserToken(body.UserToken)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	newOrder := models.TbankOrder{
		Amount:      body.Amount,
		Description: body.Description,
		InviteToken: inviteToken,
		PaymentId:   body.PaymentId,
	}

	if err := repository.CreateOrder(newOrder); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	order := repository.GetDataForPayment(body.PaymentId, inviteToken)

	tokenPayload := models.CreateTokenRequest{
		Id:          order.Id,
		TerminalKey: order.TerminalKey,
		Amount:      order.Amount,
		Description: order.Description,
		Password:    order.Password,
	}

	token := fmt.Sprintf("%x", utils.MakeToken(tokenPayload))

	orderId := strconv.Itoa(order.Id)

	receiptData := models.ReceiptData{
		FfdVersion: body.FfdVersion,
		Taxation:   body.Taxation,
		Email:      body.Email,
		Items:      body.Items,
	}

	tbankRequestBody := models.CreatePaymentPayload{
		TerminalKey: order.TerminalKey,
		Amount:      order.Amount,
		OrderId:     orderId,
		Description: order.Description,
		Token:       token,
		Receipt:     receiptData,
	}

	marshalled, err := json.Marshal(tbankRequestBody)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	res, err := http.Post(fmt.Sprintf("%s/Init", config.Config.TbankApiUrl), "application/json", bytes.NewReader(marshalled))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	defer res.Body.Close()

	resBody, _ := io.ReadAll(res.Body)

	var response models.CreatePaymentResponse
	json.Unmarshal(resBody, &response)

	newOrder.ID = uint(order.Id)
	newOrder.TbankPaymentId = response.PaymentId
	newOrder.PaymentURL = response.PaymentURL
	newOrder.Token = token

	postgres.DB.Save(&newOrder)

	return c.JSON(response)
}

// Polling for changing payment status to done
func PaymentStatusDoneHandler(c *fiber.Ctx) error {
	body := new(models.PaymentDoneRequest)

	if err := c.BodyParser(body); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if err := repository.PaymentStatusDone(body.PaymentId); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	return c.SendStatus(fiber.StatusOK)
}

// Getting payment url by payment id
func GetPaymentURLHandler(c *fiber.Ctx) error {
	id := c.QueryInt("payment_id")
	userToken := c.Query("token")

	inviteToken, err := utils.GetIviteTokenByUserToken(userToken)
	if err != nil {
		return c.SendStatus(fiber.StatusForbidden)
	}

	paymentURL := repository.GetPaymentUrlByPaymentId(id, inviteToken)
	if paymentURL != "" {
		return c.JSON(fiber.Map{"payment_url": paymentURL})
	}

	return c.SendStatus(fiber.StatusNotFound)
}

// Checking if user is logged
func CheckIsLoggedHandler(c *fiber.Ctx) error {
	userToken := c.Query("token")

	inviteToken, err := utils.GetIviteTokenByUserToken(userToken)
	if err != nil {
		return c.JSON(fiber.Map{"is_registered": false})
	}

	user := repository.GetUser(inviteToken)
	if user.ID != 0 {
		return c.JSON(fiber.Map{"is_registered": true})
	}

	return c.JSON(fiber.Map{"is_registered": false})
}
