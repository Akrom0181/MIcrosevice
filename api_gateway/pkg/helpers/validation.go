package helpers

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v3"
)

func ValidatePhone(phone string) error {
	re := regexp.MustCompile(`^[+][9][9][8]\d{9}$`)
	if !re.MatchString(phone) {
		return errors.New("invalid phone number: " + phone)
	}

	return nil
}

func ValidateDates(startDate, endDate string) error {
	var layout = "2006-01-02 15:04:05"

	from, err := time.Parse(layout, startDate)
	if err != nil {
		return errors.New("start_date is invalid")
	}

	to, err := time.Parse(layout, endDate)
	if err != nil {
		return errors.New("end_time is invalid")
	}

	if !from.Before(to) {
		return errors.New("start_time can not be greater than end_time")
	}

	return nil
}

func ValidatePassword(password string) error {
	if password == "" {
		return errors.New("password cannot be blank")
	}
	if len(password) < 8 || len(password) > 30 {
		return errors.New("password length should be 8 to 30 characters")
	}
	if validation.Validate(password, validation.Match(regexp.MustCompile("^[A-Za-z0-9$_@.#]+$"))) != nil {
		return errors.New("password should contain only alphabetic characters, numbers and special characters(@, $, _, ., #)")
	}
	if validation.Validate(password, validation.Match(regexp.MustCompile("[0-9]"))) != nil {
		return errors.New("password should contain at least one number")
	}
	if validation.Validate(password, validation.Match(regexp.MustCompile("[A-Za-z]"))) != nil {
		return errors.New("password should contain at least one alphabetic character")
	}
	return nil
}

func ValidateUsername(username string) error {
	if username == "" {
		return errors.New("username cannot be blank")
	}
	if len(username) < 5 || len(username) > 30 {
		return errors.New("username length should be 6 to 30 characters")
	}
	if validation.Validate(username, validation.Match(regexp.MustCompile("^[A-Za-z0-9$@_.#]+$"))) != nil {
		return errors.New("username should contain only alphabetic characters, numbers and special characters(@, $, _, ., #)")
	}
	return nil
}

func ValidateEmailAddress(email string) error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	if !emailRegex.MatchString(email) {
		return fmt.Errorf("email address %s is not valid", email)
	}

	return nil
}
