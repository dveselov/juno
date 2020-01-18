package providers

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

// AuthenticationByCodeProvider is a base interface for authentication providers, like UCaller
type AuthenticationByCodeProvider interface {
	GetDB() *sqlx.DB
	GetDBType() string
	Send(phoneNumber string, code string) error
}

func Begin(provider AuthenticationByCodeProvider, phoneNumber string, code string) error {
	tx, err := provider.GetDB().Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(
		"INSERT INTO authentication_provider_code(provider, phone_number, code) VALUES ($1, $2, $3)",
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	err = provider.Send(phoneNumber, code)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = stmt.Exec(provider.GetDBType(), phoneNumber, code)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func Verify(provider AuthenticationByCodeProvider, phoneNumber string, code string) error {
	var authID string
	err := provider.GetDB().QueryRow(`
		SELECT id FROM authentication_provider_code
		WHERE provider = $1
			AND phone_number = $2
			AND code = $3
			AND created_at >= current_timestamp - interval '10 days'
		ORDER BY created_at DESC
		LIMIT 1
	`, provider.GetDBType(), phoneNumber, code).Scan(&authID)
	if err != nil {
		return err
	}
	if authID == "" {
		return errors.New("Invalid or expired code")
	}
	return nil
}
