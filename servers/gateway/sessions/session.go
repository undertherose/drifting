package sessions

import (
	"errors"
	"net/http"
	"strings"
)

const headerAuthorization = "Authorization"
const paramAuthorization = "auth"
const schemeBearer = "Bearer "

//ErrNoSessionID is used when no session ID was found in the Authorization header
var ErrNoSessionID = errors.New("no session ID found in " + headerAuthorization + " header")

//ErrInvalidScheme is used when the authorization scheme is not supported
var ErrInvalidScheme = errors.New("authorization scheme not supported")

//BeginSession creates a new SessionID, saves the `sessionState` to the store, adds an
//Authorization header to the response with the SessionID, and returns the new SessionID
func BeginSession(signingKey string, store Store, sessionState interface{}, w http.ResponseWriter) (SessionID, error) {

	if len(signingKey) == 0 {
		return InvalidSessionID, errors.New("signing key length cannot be 0")
	}

	sessionID, sessionIDerr := NewSessionID(signingKey)
	if sessionIDerr != nil {
		return InvalidSessionID, sessionIDerr
	}

	if err := store.Save(sessionID, sessionState); err != nil {
		return InvalidSessionID, err
	}

	w.Header().Add(headerAuthorization, schemeBearer+sessionID.String())
	return sessionID, nil
}

//GetSessionID extracts and validates the SessionID from the request headers
func GetSessionID(r *http.Request, signingKey string) (SessionID, error) {

	authHeaderVal := r.Header.Get(headerAuthorization)
	if len(authHeaderVal) == 0 {
		authHeaderVal = r.URL.Query().Get("auth")
	}

	if !strings.Contains(authHeaderVal, schemeBearer) {
		return InvalidSessionID, ErrInvalidScheme
	}

	authHeaderVal = strings.SplitN(authHeaderVal, " ", 2)[1]
	sessionID, err := ValidateID(authHeaderVal, signingKey)
	if err != nil {
		return InvalidSessionID, ErrInvalidID
	}

	return sessionID, nil
}

//GetState extracts the SessionID from the request,
//gets the associated state from the provided store into
//the `sessionState` parameter, and returns the SessionID
func GetState(r *http.Request, signingKey string, store Store, sessionState interface{}) (SessionID, error) {

	sessionID, err := GetSessionID(r, signingKey)
	if err != nil {
		return InvalidSessionID, err
	}

	getErr := store.Get(sessionID, sessionState)
	if getErr != nil {
		return InvalidSessionID, getErr
	}

	return sessionID, nil
}

//EndSession extracts the SessionID from the request,
//and deletes the associated data in the provided store, returning
//the extracted SessionID.
func EndSession(r *http.Request, signingKey string, store Store) (SessionID, error) {

	sessionID, err := GetSessionID(r, signingKey)
	if err != nil {
		return InvalidSessionID, err
	}

	delErr := store.Delete(sessionID)
	if delErr != nil {
		return InvalidSessionID, delErr
	}

	return sessionID, nil
}
