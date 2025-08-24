package notirepositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"strconv"

	"github.com/guatom999/ecommerce-notification-api/modules"
	"github.com/jmoiron/sqlx"
)

type (
	notirepository struct {
		db *sqlx.DB
	}
)

func NewNotirepository(db *sqlx.DB) NotirepositoryInterface {
	return &notirepository{
		db: db,
	}
}

func (r *notirepository) Create(ctx context.Context, in modules.CreateInput) (string, error) {

	var userId *string
	if in.UserID != "" {
		userId = &in.UserID
	}

	data, err := json.Marshal(in.Data)
	if err != nil {
		log.Printf("Error: marshal data failed: %v", err)
		return "", errors.New("error: marshal data failed")
	}

	var resultId string

	if err := r.db.GetContext(ctx, &resultId, `	  INSERT INTO notifications (user_id, channel, recipient, template_name, data, status)
	  VALUES ($1::uuid, $2, $3, $4, $5, 'queued')
	  RETURNING id;
	`, userId, in.Channel, in.Recipient, in.TemplateName, data); err != nil {
		log.Printf("Error: created notification failed: %v", err)
		return "", errors.New("error: create failed")
	}

	return resultId, nil
}
func (r *notirepository) Get(ctx context.Context, id string) (*modules.NotificationRow, error) {

	result := new(modules.NotificationRow)

	if err := r.db.GetContext(ctx, result, `
		  SELECT id,
	         COALESCE(user_id::text,'') AS user_id_text,
	         channel, recipient, template_name, data, status,
	         created_at, updated_at
	  FROM notifications
	  WHERE id = $1
	`, id); err != nil {
		log.Printf("Error: get notification failed: %v", err)
		return nil, err
	}

	err := json.Unmarshal(result.DataRaw, &result.Data)
	if err != nil {
		log.Printf("Error: unmarshal data failed: %v", err)
		return nil, errors.New("error: unmarshal data failed")
	}

	return result, nil
}
func (r *notirepository) List(ctx context.Context, userID, status string, limit, offset int) ([]modules.NotificationListRow, error) {

	queryString := `
	SELECT id , status , recipient, created_at 
	FROM notifications
	WHERE 1 = 1
	`
	args := []any{}
	i := 1
	if userID != "" {
		queryString += ` AND user_id = $` + strconv.Itoa(i)
		args = append(args, userID)
		i++
	}
	if status != "" {
		queryString += ` AND status = $` + strconv.Itoa(i)
		args = append(args, status)
		i++
	}
	queryString += `ORDER BY created_at DESC LIMIT $` + strconv.Itoa(i) + ` OFFSET $` + strconv.Itoa(i+1)
	args = append(args, limit, offset)

	result := make([]modules.NotificationListRow, 0)

	if err := r.db.SelectContext(ctx, result, queryString, args...); err != nil {
		return nil, err
	}

	return result, nil
}
func (r *notirepository) UpdateStatus(ctx context.Context, id, status string) error {

	_, err := r.db.ExecContext(ctx,
		`
	  UPDATE notifications
	  SET status=$2, updated_at=CURRENT_TIMESTAMP
	  WHERE id=$1
	`, id, status)

	return err
}

func nullIfEmpty(s string) any {
	if s == "" {
		return sql.NullString{}
	}
	return s
}

func (r *notirepository) AttemptSendTx(ctx context.Context, id string, status string, errMsg string, providerRaw map[string]any) error {

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Printf("Error: begin transaction failed: %v", err)
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	rawJSON, _ := json.Marshal(providerRaw)

	if _, err := tx.ExecContext(ctx, `
		INSERT INTO delivery_attempts (notification_id, status, error_message, provider_raw)
		VALUES ($1, $2, $3, $4)
	`, id, status, nullIfEmpty(errMsg), rawJSON,
	); err != nil {
		return err
	}

	if _, err := r.db.ExecContext(ctx,
		`
		UPDATE notifications 
		SET status=$2 , updated_at=CURRENT_TIMESTAMP
		WHERE id=$1 
		`, id, status,
	); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
