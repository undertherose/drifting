package handlers

import (
	"time"

	"github.com/final-project-kool-kids/servers/finalgateway/models/users"
)

//SessionState represents the session state for web server
type SessionState struct {
	//Tells when the session was started
	StartTime time.Time `json:"startTime"`
	//User who started the session
	User *users.User `json:"user"`
}
