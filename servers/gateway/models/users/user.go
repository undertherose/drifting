package users

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/mail"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

//gravatarBasePhotoURL is the base URL for Gravatar image requests.
//See https://id.gravatar.com/site/implement/images/ for details
const gravatarBasePhotoURL = "https://www.gravatar.com/avatar/"

//bcryptCost is the default bcrypt cost to use when hashing passwords
var bcryptCost = 13

//User represents a user account in the database
type User struct {
	ID        int64  `json:"id"`
	Email     string `json:"-"` //never JSON encoded/decoded
	PassHash  []byte `json:"-"` //never JSON encoded/decoded
	UserName  string `json:"userName"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	PhotoURL  string `json:"photoURL"`
}

//Credentials represents user sign-in credentials
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

//NewUser represents a new user signing up for an account
type NewUser struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	PasswordConf string `json:"passwordConf"`
	UserName     string `json:"userName"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
}

//Updates represents allowed updates to a user profile
type Updates struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

//Validate validates the new user and returns an error if
//any of the validation rules fail, or nil if its valid
func (nu *NewUser) Validate() error {

	_, err := mail.ParseAddress(nu.Email)

	if err != nil {
		return fmt.Errorf("invalid email")
	} else if len(nu.Password) < 6 {
		return fmt.Errorf("password smaller than 6 characters")
	} else if nu.Password != nu.PasswordConf {
		return fmt.Errorf("password do not match")
	} else if len(nu.UserName) == 0 {
		return fmt.Errorf("username cannot be empty")
	} else if strings.Contains(nu.UserName, " ") {
		return fmt.Errorf("username cannot contain any spaces")
	}

	if !strings.Contains(nu.Email, "uw.edu") {
		return fmt.Errorf("Not a uw.edu")
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
		Email:     nu.Email,
		UserName:  nu.UserName,
		FirstName: nu.FirstName,
		LastName:  nu.LastName,
	}
	user.SetPassword(nu.Password)
	emailTrimAndLower := strings.ToLower(strings.Trim(nu.Email, " "))
	hasher := md5.New()
	hasher.Write([]byte(emailTrimAndLower))
	emailHash := hex.EncodeToString(hasher.Sum(nil))
	user.PhotoURL = gravatarBasePhotoURL + emailHash

	//TODO: also call .SetPassword() to set the PassHash
	//field of the User to a hash of the NewUser.Password

	return user, nil
}

//FullName returns the user's full name, in the form:
// "<FirstName> <LastName>"
//If either first or last name is an empty string, no
//space is put between the names. If both are missing,
//this returns an empty string
func (u *User) FullName() string {

	if u.FirstName == "" && u.LastName == "" {
		return ""
	} else if u.FirstName == "" {
		return u.LastName
	} else if u.LastName == "" {
		return u.FirstName
	} else {
		return u.FirstName + " " + u.LastName
	}
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

//ApplyUpdates applies the updates to the user. An error
//is returned if the updates are invalid
func (u *User) ApplyUpdates(updates *Updates) error {
	u.FirstName = updates.FirstName
	u.LastName = updates.LastName

	return nil
}
