package providers

import (
	"database/sql"
	"errors"
)

// AuthenticationByCodeProvider is a base interface for authentication providers, like UCaller
type AuthenticationByCodeProvider interface {
	GetDBType() string
	Begin(phoneNumber string, code string) error
	Send(phoneNumber string, code string) error
	Verify(phoneNumber string, code string) (bool, error)
}

type BaseAuthenticationProvider struct {
	db *sql.DB
}

func (p *BaseAuthenticationProvider) GetDBType() string {
	return ""
}

func (p *BaseAuthenticationProvider) Send(phoneNumber string, code string) error {
	return errors.New("Not implemented")
}

func (p *BaseAuthenticationProvider) Begin(phoneNumber string, code string) error {
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(
		"INSERT INTO authentication_provider_code VALUES ($1, $2, $3)",
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = p.Send(phoneNumber, code)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(p.GetDBType(), phoneNumber, code)
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

func (p *BaseAuthenticationProvider) Verify(phoneNumber string, code string) error {
	var authID string
	err := p.db.QueryRow(`
		SELECT id FROM authentication_provider_code
		WHERE provider = $1
			AND phone_number = $2
			AND code = $3
			AND created_at >= current_timestamp - interval '1 day'
		ORDER BY created_at DESC
		LIMIT 1
	`, p.GetDBType(), phoneNumber, code).Scan(&authID)
	if err != nil {
		return err
	}
	if authID == "" {
		return errors.New("Invalid or expired code")
	}
	return nil
}
