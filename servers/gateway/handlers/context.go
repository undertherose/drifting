package handlers

import (
	"github.com/final-project-kool-kids/servers/finalgateway/models/users"
	"github.com/final-project-kool-kids/servers/finalgateway/sessions"
)

//HandlerContext contains information about the key used to sign and validate sessionIDs, session store and user store
type HandlerContext struct {
	Key          string
	UserStore    users.Store
	SessionStore sessions.Store
	Notifier     *Notifier
}
