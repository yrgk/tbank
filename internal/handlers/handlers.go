package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

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

	amount := int(body.Amount * 100)

	newOrder := models.TbankOrder{
		Amount:      amount,
		Description: body.Description,
		InviteToken: inviteToken,
		PaymentId:   body.PaymentId,
		DocsSalesId: body.DocsSalesId,
	}

	if err := repository.CreateOrder(newOrder); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	order := repository.GetDataForPayment(body.PaymentId, inviteToken)

	expiration := time.Now().AddDate(0, 0, 90).Format("2006-01-02T15:04:05-07:00")

	tokenPayload := models.CreateTokenRequest{
		Id:          order.Id,
		TerminalKey: order.TerminalKey,
		Amount:      amount,
		Description: order.Description,
		Password:    order.Password,
		RedirectDueDate: expiration,
	}

	utils.ModifyItems(&body.Items)

	token := fmt.Sprintf("%x", utils.MakeToken(tokenPayload))

	orderId := strconv.Itoa(order.Id)

	receiptData := models.ReceiptData{
		FfdVersion: body.FfdVersion,
		Taxation:   body.Taxation,
		Email:      body.Email,
		Items:      body.Items,
	}

	tbankRequestBody := models.CreatePaymentPayload{
		TerminalKey:     order.TerminalKey,
		Amount:          order.Amount,
		OrderId:         orderId,
		Description:     order.Description,
		RedirectDueDate: expiration,
		Token:           token,
		Receipt:         receiptData,
	}
	fmt.Println(tbankRequestBody)

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
	paymentId := c.QueryInt("payment_id")
	docsSalesId := c.QueryInt("docs_sales_id")
	userToken := c.Query("token")

	inviteToken, err := utils.GetIviteTokenByUserToken(userToken)
	if err != nil {
		return c.SendStatus(fiber.StatusForbidden)
	}

	if paymentId == 0 && docsSalesId == 0 {
		return c.SendStatus(fiber.StatusNotFound)
	}

	if paymentId == 0 && docsSalesId != 0 {
		paymentURL := repository.GetPaymentUrlByDocsSalesId(docsSalesId, inviteToken)
		if paymentURL != "" {
			return c.JSON(fiber.Map{"payment_url": paymentURL})
		}
	}

	if paymentId != 0 && docsSalesId == 0 {
		paymentURL := repository.GetPaymentUrlByPaymentId(paymentId, inviteToken)
		if paymentURL != "" {
			return c.JSON(fiber.Map{"payment_url": paymentURL})
		}
	}

	if paymentId != 0 && docsSalesId != 0 {
		paymentURL := repository.GetPaymentUrlByPaymentId(paymentId, inviteToken)
		if paymentURL != "" {
			return c.JSON(fiber.Map{"payment_url": paymentURL})
		}
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
