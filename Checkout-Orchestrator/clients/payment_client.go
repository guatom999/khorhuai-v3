package clients

import (
	"context"

	"github.com/go-resty/resty/v2"
)

type (
	PaymentClient struct {
		baseUrl string
		htpp    *resty.Client
	}
)

func NewPaymentClient(baseUrl string) *PaymentClient {
	return &PaymentClient{
		baseUrl: baseUrl,
		htpp:    resty.New(),
	}
}

func (c *PaymentClient) Create(ctx context.Context, orderId, userId string, amount int64, currency string) (string, error) {

	type paymentOutput struct {
		PaymentId string `json:"payment_id"`
	}

	paymentOut := new(paymentOutput)

	if _, err := c.htpp.R().SetContext(ctx).SetBody(map[string]any{
		"order_id": orderId,
		"user_id":  userId,
		"amount":   amount,
		"currency": currency,
	}).SetResult(paymentOut).Post(c.baseUrl + "/app/v1/payments"); err != nil {
		return "", err
	}

	return paymentOut.PaymentId, nil
}
