package utils

import (
	"errors"
	"strconv"
)

func IsDigit(text string) bool {
	if _, err := strconv.Atoi(text); err != nil {
		return false
	}
	return true
}

func IsInvalidPhoneNumber(phoneNumber string) error {
	if len(phoneNumber) != 11 {
		return errors.New("Phone number must be a string with 11 digits, e.g.: 79991112233")
	}
	if isDigit := IsDigit(phoneNumber); !isDigit {
		return errors.New("Phone number must be a integer, e.g.: 79991112233")
	}
	return nil
}

func IsInvalidVerificationCode(code string) error {
	if len(code) != 4 {
		return errors.New("Verification code must be a string with 4 digits")
	}
	if isDigit := IsDigit(code); !isDigit {
		return errors.New("Verification code must contain only digits")
	}
	return nil
}
