package modules

type CreateNotificationRequest struct {
	UserID       string         `json:"user_id,omitempty"`
	Channel      string         `json:"channel" validate:"required"`   // email|sms|push
	Recipient    string         `json:"recipient" validate:"required"` // email/phone/device
	TemplateName string         `json:"template_name,omitempty"`
	Data         map[string]any `json:"data,omitempty"` // payload/render vars
}

type SendAttemptRequest struct {
	Status       string         `json:"status" validate:"required"` // sent|failed
	ErrorMessage string         `json:"error_message,omitempty"`
	ProviderRaw  map[string]any `json:"provider_raw,omitempty"` // ใส่ response/provider id ได้
}

type UpdateNotificationStatusRequest struct {
	Status string `json:"status" validate:"required"` // เช่น cancelled
}

type CreateInput struct {
	UserID       string
	Channel      string
	Recipient    string
	TemplateName string
	Data         map[string]any
}
