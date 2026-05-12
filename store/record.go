package store

import (
	"database/sql"
	"time"
)

type Record struct {
	ID           int
	CreatedAt    string
	Currency     string
	TotalBalance float64
}

func InsertRecords(balances map[string]float64) error {
	now := time.Now().UTC().Format(time.RFC3339)
	for currency, balance := range balances {
		_, err := db.Exec(
			"INSERT INTO records (created_at, currency, total_balance) VALUES (?, ?, ?)",
			now, currency, balance,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func LastRecords() ([]Record, string, error) {
	var ts string
	err := db.QueryRow(
		"SELECT created_at FROM records ORDER BY created_at DESC LIMIT 1",
	).Scan(&ts)
	if err == sql.ErrNoRows {
		return nil, "", nil
	}
	if err != nil {
		return nil, "", err
	}

	rows, err := db.Query(
		"SELECT id, created_at, currency, total_balance FROM records WHERE created_at = ? ORDER BY currency",
		ts,
	)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	var records []Record
	for rows.Next() {
		var r Record
		if err := rows.Scan(&r.ID, &r.CreatedAt, &r.Currency, &r.TotalBalance); err != nil {
			return nil, "", err
		}
		records = append(records, r)
	}
	return records, ts, rows.Err()
}
