package models

type (
	CreatePaymentResponse struct {
		Success     bool
		ErrorCode   string
		TerminalKey string
		Status      string
		PaymentId   string
		OrderId     string
		Amount      int
		PaymentURL  string
	}

	Cboxes struct {
		Name    string `json:"name"`
		Token   string `json:"token"`
		Balance int    `json:"balance"`
	}

	InviteTokenResponse struct {
		InviteToken string `json:"invite_token"`
		Cboxes      Cboxes `json:"cboxess"`
	}
)
