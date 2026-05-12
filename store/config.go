package store

import "database/sql"

func SetConfig(key, value string) error {
	_, err := db.Exec(
		"INSERT OR REPLACE INTO config (key, value) VALUES (?, ?)",
		key, value,
	)
	return err
}

func GetConfig(key string) (string, error) {
	var value string
	err := db.QueryRow(
		"SELECT value FROM config WHERE key = ?", key,
	).Scan(&value)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return value, err
}
