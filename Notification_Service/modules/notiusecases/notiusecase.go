package notiusecases

import (
	"context"

	"github.com/guatom999/ecommerce-notification-api/modules"
)

type (
	NotiusecaseInterface interface {
		Create(ctx context.Context, in modules.CreateInput) (string, error)
		Get(ctx context.Context, id string) (*modules.NotificationRow, error)
		List(ctx context.Context, userID, status string, limit, offset int) ([]modules.NotificationListRow, error)
		UpdateStatus(ctx context.Context, id, status string) error
		AttemptSend(ctx context.Context, id string, status string, errMsg string, providerRaw map[string]any) error
	}
)
