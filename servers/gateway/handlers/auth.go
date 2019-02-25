package handlers

import (
	"encoding/json"
	"net/http"
	"path"
	"strconv"
	"time"

	"github.com/drifting/servers/gateway/models/users"
	"github.com/drifting/servers/gateway/sessions"
	"golang.org/x/crypto/bcrypt"
)

//NewHandlerContext constructs a new HandlerContext
func NewHandlerContext(key string, userStore users.Store, sessionStore sessions.Store) *HandlerContext {
	if userStore == nil {
		panic("UserStore is nil")
	} else if sessionStore == nil {
		panic("SessionStore is nil")
	}
	return &HandlerContext{
		Key:          key,
		UserStore:    userStore,
		SessionStore: sessionStore,
	}
}

//UsersHandler handles requests for the "users" resource (sign up)
func (ctx *HandlerContext) UsersHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Request body must be in JSON", http.StatusUnsupportedMediaType)
		return
	}

	newUser := users.NewUser{}
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, "New user could not be decoded", http.StatusInternalServerError)
		return
	}

	user, err := newUser.ToUser()
	if err != nil {
		http.Error(w, "New user could not be encoded to a user", http.StatusInternalServerError)
		return
	}

	dbUser, err := ctx.UserStore.Insert(user)
	if err != nil {
		http.Error(w, "User could not be inserted into database", http.StatusInternalServerError)
		return
	}

	_, err = sessions.BeginSession(ctx.Key, ctx.SessionStore, SessionState{StartTime: time.Now(), User: user}, w)
	if err != nil {
		http.Error(w, "Could not begin new session", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(dbUser)
	if err != nil {
		http.Error(w, "Could not encode user", http.StatusInternalServerError)
		return
	}

}

//SpecificUserHandler handles requests for a specific user (sign in)
func (ctx *HandlerContext) SpecificUserHandler(w http.ResponseWriter, r *http.Request) {

	sessionState := &SessionState{}
	_, err := sessions.GetState(r, ctx.Key, ctx.SessionStore, sessionState)
	if err != nil {
		http.Error(w, "Unauthorized user", http.StatusUnauthorized)
		return
	}

	//TODO: me path: gets personal stuff
	//TODO: Get rid of ability to see other users
	id := path.Base(r.URL.Path)
	var numID int64
	if id == "me" {
		numID = sessionState.User.ID
	} else {
		numID, _ = strconv.ParseInt(id, 10, 64)
	}

	if r.Method == http.MethodGet {

		user, err := ctx.UserStore.GetByID(numID)
		if err != nil {
			http.Error(w, "User not found in DB", http.StatusNotFound)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		err = json.NewEncoder(w).Encode(user)
		if err != nil {
			http.Error(w, "User could not be encoded", http.StatusInternalServerError)
			return
		}

	} else if r.Method == http.MethodPatch {
		if numID != sessionState.User.ID {
			http.Error(w, "IDs do not match", http.StatusUnauthorized)
			return
		}

		dbUser, err := ctx.UserStore.GetByID(numID)
		if err != nil {
			http.Error(w, "User not found in DB", http.StatusNotFound)
			return
		}

		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "Request body must be in JSON", http.StatusInternalServerError)
			return
		}

		userUpdate := &users.Updates{}
		err = json.NewDecoder(r.Body).Decode(userUpdate)
		if err != nil {
			http.Error(w, "An error occured while decoding the user updates", http.StatusInternalServerError)
			return
		}

		if userUpdate.FirstName == "" || userUpdate.LastName == "" || userUpdate.FirstName == sessionState.User.Email || userUpdate.FirstName == sessionState.User.UserName || userUpdate.LastName == sessionState.User.Email || userUpdate.LastName == sessionState.User.UserName {
			http.Error(w, "Invalid update values provided", http.StatusBadRequest)
			return
		}

		err = dbUser.ApplyUpdates(userUpdate)
		if err != nil {
			http.Error(w, "Updates could not be applied to the user", http.StatusInternalServerError)
			return
		}

		updatedUser, err := ctx.UserStore.Update(dbUser.ID, userUpdate)
		if err != nil {
			http.Error(w, "Updates could not be applied to the user", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		err = json.NewEncoder(w).Encode(updatedUser)
		if err != nil {
			http.Error(w, "Updated user could not be encoded", http.StatusInternalServerError)
			return
		}

	} else {
		http.Error(w, "Method type not supported", http.StatusMethodNotAllowed)
		return
	}

}

//SessionsHandler handles requests for the sessions resource (sign in)
func (ctx *HandlerContext) SessionsHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method type not supported", http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Request body must be in JSON", http.StatusUnsupportedMediaType)
		return
	}

	userCred := &users.Credentials{}
	err := json.NewDecoder(r.Body).Decode(userCred)
	if err != nil {
		http.Error(w, "Could not decode user credentials", http.StatusBadRequest)
		return
	}

	user, err := ctx.UserStore.GetByEmail(userCred.Email)
	if err != nil {
		bcrypt.GenerateFromPassword(user.PassHash, 13) // Ensure that process takes time
		http.Error(w, "Unauthorized user", http.StatusUnauthorized)
		return
	}

	err = user.Authenticate(userCred.Password)
	if err != nil {
		http.Error(w, "Unauthorized user", http.StatusUnauthorized)
		return
	}

	_, err = sessions.BeginSession(ctx.Key, ctx.SessionStore, SessionState{StartTime: time.Now(), User: user}, w)
	if err != nil {
		http.Error(w, "Could not begin a new session", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, "Could not encode user", http.StatusInternalServerError)
		return
	}

}

//SpecificSessionHandler handles requests related to specific authenticated sessions (sign out)
func (ctx *HandlerContext) SpecificSessionHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodDelete {
		http.Error(w, "Method type not supported", http.StatusMethodNotAllowed)
		return
	}

	if path.Base(r.URL.Path) != "mine" {
		http.Error(w, "Users can only delete their own sessions", http.StatusForbidden)
		return
	}

	_, err := sessions.EndSession(r, ctx.Key, ctx.SessionStore)
	if err != nil {
		http.Error(w, "An error occured while ending the session", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "text/plain")
	w.Write([]byte("signed out"))
	w.WriteHeader(http.StatusOK)

}

//GetAllUsersHandler gets all the users from the db
func (ctx *HandlerContext) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method type not supported", http.StatusMethodNotAllowed)
		return
	}

	users, err := ctx.UserStore.GetAll()

	if err != nil {
		http.Error(w, "Couldn't retrieve users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		http.Error(w, "Users list could not be encoded", http.StatusInternalServerError)
		return
	}

}
