package modules

import "time"

type NotificationRow struct {
	ID           string         `db:"id"            json:"id"`
	UserID       string         `db:"user_id_text"  json:"user_id,omitempty"`
	Channel      string         `db:"channel"       json:"channel"`
	Recipient    string         `db:"recipient"     json:"recipient"`
	TemplateName string         `db:"template_name" json:"template_name,omitempty"`
	Data         map[string]any `db:"-"             json:"data,omitempty"`
	DataRaw      []byte         `db:"data"`
	Status       string         `db:"status"        json:"status"`
	CreatedAt    time.Time      `db:"created_at"    json:"created_at"`
	UpdatedAt    time.Time      `db:"updated_at"    json:"updated_at"`
}

type NotificationListRow struct {
	ID        string    `db:"id"         json:"id"`
	Status    string    `db:"status"     json:"status"`
	Channel   string    `db:"channel"    json:"channel"`
	Recipient string    `db:"recipient"  json:"recipient"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
