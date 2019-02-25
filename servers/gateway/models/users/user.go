package users

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

//bcryptCost is the default bcrypt cost to use when hashing passwords
var bcryptCost = 13

//User represents a user account in the database
type User struct {
	ID       int64  `json:"id"`
	PassHash []byte `json:"-"` //never JSON encoded/decoded
	UserName string `json:"userName"`
	Type     string `json:"-"` //never JSON encoded/decoded
}

//Credentials represents user sign-in credentials
type Credentials struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

//NewUser represents a new user signing up for an account
type NewUser struct {
	UserName     string `json:"userName"`
	Password     string `json:"password"`
	PasswordConf string `json:"passwordConf"`
	Type         string `json:"type"`
}

//Updates for changing password
//type Updates struct {
//	UserName        string `json:"firstName"`
//	CurrPassword    string `json:"currPassword"`
//	NewPassword     string `json:"newPassword"`
//	NewPasswordConf string `json:"newPasswordConf"`
//}

//Validate validates the new user and returns an error if
//any of the validation rules fail, or nil if its valid
func (nu *NewUser) Validate() error {

	if len(nu.UserName) == 0 {
		return fmt.Errorf("Username cannot be empty")
	} else if strings.Contains(nu.UserName, " ") {
		return fmt.Errorf("Username cannot contain any spaces")
	} else if len(nu.Password) < 6 {
		return fmt.Errorf("Password cannot have less than 6 characters")
	} else if nu.Password != nu.PasswordConf {
		return fmt.Errorf("Passwords do not match")
	}
	return nil
}

//ToUser converts the NewUser to a User, setting the
//PhotoURL and PassHash fields appropriately
func (nu *NewUser) ToUser() (*User, error) {

	if err := nu.Validate(); err != nil {
		return nil, err
	}
	user := &User{
		UserName: nu.UserName,
	}
	user.SetPassword(nu.Password)
	return user, nil
}

//SetPassword hashes the password and stores it in the PassHash field
func (u *User) SetPassword(password string) error {
	hashPass, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)

	if err == nil {
		u.PassHash = hashPass
	}

	return err
}

//Authenticate compares the plaintext password against the stored hash
//and returns an error if they don't match, or nil if they do
func (u *User) Authenticate(password string) error {
	if err := bcrypt.CompareHashAndPassword(u.PassHash, []byte(password)); err != nil {
		return err
	}
	return nil
}

// TODO:
// NOT MVP
// Have Admin be able to change user types
// Have mods be able to change user types
