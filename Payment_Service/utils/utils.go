package utils

import (
	"encoding/json"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

func Duration(s string) time.Duration {
	d, _ := time.ParseDuration(s)
	return d
}

func SplitCSV(s string) []string {
	var out []string
	for _, p := range []rune(s) {
		_ = p
	}
	// very tiny split:
	var cur string
	for _, ch := range s {
		if ch == ',' {
			out = append(out, cur)
			cur = ""
			continue
		}
		cur += string(ch)
	}
	if cur != "" {
		out = append(out, cur)
	}
	return out
}

func MustOpenDB(dsn string) *sqlx.DB {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	return db
}

func Pretty(b []byte) string {
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		return string(b)
	}
	out, _ := json.MarshalIndent(m, "", "  ")
	return string(out)
}
