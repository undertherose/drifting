package sessions

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"strings"
)

//InvalidSessionID represents an empty, invalid session ID
const InvalidSessionID SessionID = ""

//idLength is the length of the ID portion
const idLength = 32

//signedLength is the full length of the signed session ID
//(ID portion plus signature)
const signedLength = idLength + sha256.Size

//SessionID represents a valid, digitally-signed session ID.
//This is a base64 URL encoded string created from a byte slice
//where the first `idLength` bytes are crytographically random
//bytes representing the unique session ID, and the remaining bytes
//are an HMAC hash of those ID bytes (i.e., a digital signature).
//The byte slice layout is like so:
//+-----------------------------------------------------+
//|...32 crypto random bytes...|HMAC hash of those bytes|
//+-----------------------------------------------------+
type SessionID string

//ErrInvalidID is returned when an invalid session id is passed to ValidateID()
var ErrInvalidID = errors.New("Invalid Session ID")

//NewSessionID creates and returns a new digitally-signed session ID,
//using `signingKey` as the HMAC signing key. An error is returned only
//if there was an error generating random bytes for the session ID
func NewSessionID(signingKey string) (SessionID, error) {

	if len(signingKey) == 0 {
		return InvalidSessionID, errors.New("signing key maybe empty")
	}

	// create a byte slice where the first `idLength` of bytes
	// are cryptographically random bytes for the new session ID,
	// and the remaining bytes are an HMAC hash of those ID bytes,
	// using the provided `signingKey` as the HMAC key.
	// encode that byte slice using base64 URL Encoding and return
	// the result as a SessionID type

	byteSlice, _ := GenerateRandomBytes(idLength)
	mac := hmac.New(sha256.New, []byte(signingKey))
	mac.Write(byteSlice)
	hashSlice := mac.Sum(nil)
	byteSlice = append(byteSlice, hashSlice...)
	encodedSlice := SessionID(base64.URLEncoding.EncodeToString(byteSlice))

	return encodedSlice, nil
}

//GenerateRandomBytes generates a slice of `n` random bytes
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)

	if err != nil {
		return nil, err
	}

	return b, nil
}

//ValidateID validates the string in the `id` parameter
//using the `signingKey` as the HMAC signing key
//and returns an error if invalid, or a SessionID if valid
func ValidateID(id string, signingKey string) (SessionID, error) {

	if strings.Contains(id, ":") {
		id = strings.Split(id, ":")[1]
	}

	decodedSlice, _ := base64.URLEncoding.DecodeString(id)
	if len(decodedSlice) < idLength {
		return InvalidSessionID, ErrInvalidID
	}

	idByteSlice := decodedSlice[0:idLength]
	mac := hmac.New(sha256.New, []byte(signingKey))
	mac.Write(idByteSlice)
	sliceToCompare := mac.Sum(nil)
	if bytes.Equal(decodedSlice[idLength:], sliceToCompare) {
		return SessionID(id), nil
	}

	return InvalidSessionID, ErrInvalidID
}

//String returns a string representation of the sessionID
func (sid SessionID) String() string {
	return string(sid)
}
