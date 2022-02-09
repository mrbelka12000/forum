package database

import (
	"database/sql"
	"errors"
	forum "forum/internal"
	"strings"
)

//CanRegister ..
func CanRegister(db *sql.DB, user forum.User) error {
	if !checkForSpace(user.Password) || !checkForSpace(user.Name) {
		return errors.New("Please, do not use space")
	}
	if user.Password != user.Confirm {
		return errors.New("Password confirmation must match Password")
	}
	if err := ValidatingName(user.Name); err != nil {
		return err
	}

	if err := ValidatingLogin(user.Login); err != nil {
		return err
	}

	if err := ValidatingPassword(user.Password); err != nil {
		return err
	}

	return InsertNewUserIntoDB(db, user)
}

func checkForSpace(s string) bool {
	if s == "" {
		return false
	}
	for _, v := range s {
		if v == 32 {
			return false
		}
	}
	return true
}

//ValidatingName ..
func ValidatingName(name string) error {
	if len(name) > 30 {
		return errors.New("Name must contain less than 30 characters")
	}
	return nil
}

//ValidatingPassword ..
func ValidatingPassword(password string) error {
	var checked bool
	for _, l := range password {
		if l < 32 || l > 126 {
			checked = true
		}
	}
	if checked {
		return errors.New("Please, use only latin alphabet")
	}
	if len(password) < 8 {
		return errors.New("Password must contain more than 8 characters")
	}
	var numbers, bigletters, smallletters, specialsymbols bool
	for _, v := range password {
		if v >= 'a' && v <= 'z' {
			smallletters = true
		} else if v >= 'A' && v <= 'Z' {
			bigletters = true
		} else if v >= '0' && v <= '9' {
			numbers = true
		} else {
			specialsymbols = true
		}
	}
	if !numbers {
		return errors.New("Password must contains numbers")
	}
	if !bigletters {
		return errors.New("Password must contains uppercase")
	}
	if !smallletters {
		return errors.New("Password must contains lowercase")
	}
	if !specialsymbols {
		return errors.New("Password must contains symbols")
	}
	return nil
}

//ValidatingLogin ..
func ValidatingLogin(login string) error {
	email := strings.Split(login, "@")

	if len(email) != 2 {
		return errors.New("Invalid email")
	}

	if len(email[0]) < 5 {
		return errors.New("Not enough charaterters in domain")
	}

	for _, v := range email[0] {
		if (v >= 32 && v <= 47) || (v < 32 || v > 122) || (v >= 58 && v <= 63) || (v >= 92 && v <= 96) {
			if v == 46 {
				continue
			}
			return errors.New("Invalid domain addres")
		}
	}

	if len(email[1]) < 5 {
		return errors.New("Not enough charaterters in domain")
	}

	count := 0
	for i, v := range email[1] {
		if (v >= 32 && v <= 47) || (v < 32 || v > 122) || (v >= 58 && v <= 63) || (v >= 92 && v <= 96) || (v >= 65 && v <= 90) {
			if v == 46 {
				if i == len(email[1])-1 {
					return errors.New("Invalid addres")
				}
				count++
				continue
			}
			return errors.New("Invalid domain addres")
		}
		if count == 2 {
			return errors.New("Invalid addres")
		}
	}
	if count == 0 {
		return errors.New("Invalid addres")
	}
	return nil
}

/*

	CanRegister проверяем на уникальность


*/
