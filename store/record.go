package store

import (
	"path/filepath"
	"time"
)

type Record struct {
	CreatedAt    string  `json:"created_at"`
	Currency     string  `json:"currency"`
	TotalBalance float64 `json:"total_balance"`
}

func recordsPath() string { return filepath.Join(dataDir, "records.json") }

func InsertRecords(balances map[string]float64) error {
	var records []Record
	load(recordsPath(), &records)
	now := time.Now().UTC().Format(time.RFC3339)
	for currency, balance := range balances {
		records = append(records, Record{
			CreatedAt:    now,
			Currency:     currency,
			TotalBalance: balance,
		})
	}
	return save(recordsPath(), records)
}

func LastRecords() ([]Record, string, error) {
	var records []Record
	if err := load(recordsPath(), &records); err != nil {
		return nil, "", err
	}
	if len(records) == 0 {
		return nil, "", nil
	}

	var latest string
	for _, r := range records {
		if r.CreatedAt > latest {
			latest = r.CreatedAt
		}
	}

	var result []Record
	for _, r := range records {
		if r.CreatedAt == latest {
			result = append(result, r)
		}
	}
	return result, latest, nil
}

func RecordsSince(d time.Duration) ([]Record, error) {
	var all []Record
	if err := load(recordsPath(), &all); err != nil {
		return nil, err
	}
	since := time.Now().UTC().Add(-d)
	var result []Record
	for _, r := range all {
		t, err := time.Parse(time.RFC3339, r.CreatedAt)
		if err != nil {
			continue
		}
		if t.After(since) || t.Equal(since) {
			result = append(result, r)
		}
	}
	return result, nil
}
