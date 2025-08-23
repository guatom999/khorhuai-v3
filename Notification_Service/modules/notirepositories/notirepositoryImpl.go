package notirepositories

import "github.com/jmoiron/sqlx"

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
