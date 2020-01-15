package providers

import (
	"database/sql"
	"errors"
	"time"
)

// AuthorizationByCodeProvider is a base interface for authorization providers, like UCaller
type AuthorizationByCodeProvider interface {
	GetDBType() string
	Begin(phoneNumber string, code string) error
	Send(phoneNumber string, code string) error
	Verify(phoneNumber string, code string) (bool, error)
}

type BaseAuthorizationProvider struct {
	db *sql.DB
}

func (p *BaseAuthorizationProvider) GetDBType() string {
	return ""
}

func (p *BaseAuthorizationProvider) Send(phoneNumber string, code string) error {
	return errors.New("Not implemented")
}

func (p *BaseAuthorizationProvider) Begin(phoneNumber string, code string) error {
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(
		"INSERT INTO authorization_provider_code VALUES ($1, $2, $3, $4)",
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = p.Send(phoneNumber, code)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(p.GetDBType(), phoneNumber, code, time.Now())
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

func (p *BaseAuthorizationProvider) Verify(phoneNumber string, code string) error {
	var authID string
	err := p.db.QueryRow(`
		SELECT id FROM authorization_provider_code
		WHERE provider = $1
			AND phone_number = $2
			AND code = $3
			AND issued_on >= current_timestamp - interval '1 day'
		ORDER BY issued_on DESC
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
