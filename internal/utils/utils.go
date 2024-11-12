package utils

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"payment-microservice/config"
	"payment-microservice/internal/models"
)

func MakeToken(body models.CreateTokenRequest) [32]byte {
	orderId := strconv.Itoa(body.Id)
	amount := fmt.Sprintf("%d", body.Amount)

	strings := []string{amount, body.Description, orderId, body.Password, body.RedirectDueDate, body.TerminalKey}

	var bs string

	for _, val := range strings {
		bs += val
	}

	res := sha256.Sum256([]byte(bs))

	return res
}

func GetPaymentById(id int, token string) []byte {
	res, err := http.Get(fmt.Sprintf("%s/api/v1/payment/%d/?token=%s", config.Config.CrmApiUrl, id, token))
	if err != nil {
		return []byte{}
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return []byte{}
	}

	return body
}

func ModifyItems(items *[]models.Item) {
	for i := range *items {
		(*items)[i].Price *= 100
		(*items)[i].Amount *= 100
	}
}

func GetIviteTokenByUserToken(userToken string) (string, error) {
	url := fmt.Sprintf("%s/api/v1/cashboxes_meta/?token=%s", config.Config.CrmApiUrl, userToken)

	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	if res.StatusCode != 200 {
		return "", errors.New("bad request")
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	var result models.InviteTokenResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", errors.New("bad request")
	}

	return result.InviteToken, nil
}