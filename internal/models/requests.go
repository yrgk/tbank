package models

type (

	ReceiptData struct {
		FfdVersion string
		Taxation   string
		Email      string
		Items      []Item
	}

	CreatePaymentPayload struct {
		OrderId     string
		TerminalKey string
		Amount      int
		Description string
		Token       string
		Receipt     ReceiptData
	}
	Item struct {
		Name            string
		Price           int
		Quantity        int
		Amount          int
		Tax             string
		PaymentMethod   string
		PaymentObject   string
		MeasurementUnit string
	}

	CreatePaymentRequest struct {
		// Для создания ссылки на оплату
		Amount      int
		Description string
		UserToken   string
		PaymentId   int
		// Для выдачи чеков
		FfdVersion string
		Email      string
		Taxation   string
		Items      []Item
	}

	CreateTokenRequest struct {
		Id          int
		TerminalKey string
		Amount      int
		Description string
		Password    string
	}

	CreateAccountRequest struct {
		TerminalKey string
		UserToken   string
		Password    string
	}

	CreateAccountPayload struct {
		TerminalKey string
		InviteToken string
		Password    string
	}

	PaymentDoneData struct {
		PaymentId    string
		Order_number string
	}

	PaymentDoneRequest struct {
		PaymentId int
	}
)
