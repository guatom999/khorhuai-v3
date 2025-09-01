package workflows

import (
	"fmt"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type (
	Item struct {
		ProductID string `json:"product_id"`
		Quantity  int    `json:"quantity"`
	}

	CheckoutInput struct {
		OrderID     string
		UserID      string
		Items       []Item
		AmountCents int64
		Currency    string
		TTLSeconds  int
	}

	PaymentSignal struct {
		Status    string // "success" or "failed"
		PaymentId string
	}
)

func CheckoutWorkflow(ctx workflow.Context, in CheckoutInput) (string, error) {

	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
		HeartbeatTimeout:    10 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2,
			MaximumAttempts:    5,
		},
	}

	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	var reservationId string
	if err := workflow.ExecuteActivity(ctx, "Activities.ReserveStock", in.OrderID, in.Items, in.TTLSeconds).Get(ctx, &reservationId); err != nil {
		return "", fmt.Errorf("reserve stock: %w", err)
	}
	defer workflow.ExecuteActivity(ctx, "Activities.ReleaseStock", reservationId)

	if err := workflow.ExecuteActivity(ctx, "Activities.CreateOrder", in.OrderID, in.UserID, in.Items, in.Currency, in.AmountCents).Get(ctx, nil); err != nil {
		return "", fmt.Errorf("create order: %w", err)
	}
	defer workflow.ExecuteActivity(ctx, "Activities.CancelOrder", in.OrderID)

	var paymentID string
	if err := workflow.ExecuteActivity(ctx, "Activities.CreatePayment", in.OrderID, in.UserID, in.AmountCents, in.Currency).Get(ctx, &paymentID); err != nil {
		return "", fmt.Errorf("create payment: %w", err)
	}

	sigCh := workflow.GetSignalChannel(ctx, "payment_signal")
	var sig PaymentSignal
	received, _ := workflow.AwaitWithTimeout(ctx, 15*time.Minute, func() bool {
		return sigCh.ReceiveAsync(&sig)
	})
	if sig.Status != "" {
		received = true
	}

	if !received || sig.Status != "succeeded" {
		return "", fmt.Errorf("payment not successful")
	}

	if err := workflow.ExecuteActivity(ctx, "Activities.ConfirmOrder", in.OrderID).Get(ctx, nil); err != nil {
		_ = workflow.ExecuteActivity(ctx, "Activities.RefundPayment", paymentID).Get(ctx, nil)
		return "", fmt.Errorf("confirm order: %w", err)
	}

	if err := workflow.ExecuteActivity(ctx, "Activities.CommitStock", reservationId).Get(ctx, nil); err != nil {
		_ = workflow.ExecuteActivity(ctx, "Activities.RefundPayment", paymentID).Get(ctx, nil)
		return "", fmt.Errorf("commit stock: %w", err)
	}

	return in.OrderID, nil
}
