package handlers

import (
	"fmt"
	"net/http"

	"github.com/final-project-kool-kids/servers/finalgateway/sessions"
	"github.com/gorilla/websocket"
)

//WebSocketsHandler is a handler for WebSocket upgrade requests
//client lets server know that it wants to establish a connection
type WebSocketsHandler struct {
	upgrader websocket.Upgrader
	ctx      HandlerContext
}

//NewWebSocketsHandler constructs a new WebSocketsHandler
func NewWebSocketsHandler(ctx HandlerContext) *WebSocketsHandler {
	return &WebSocketsHandler{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return r.Header.Get("Origin") == "https://iqueue.zubinchopra.me"
			},
		},
		ctx: ctx,
	}
}

//ServeHTTP implements the http.Handler interface for the WebSocketsHandler
func (wsh *WebSocketsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sessionState := &SessionState{}
	_, err := sessions.GetState(r, wsh.ctx.Key, wsh.ctx.SessionStore, sessionState)
	if err != nil {
		http.Error(w, fmt.Sprintf("error not authorized: %v", err), http.StatusUnauthorized)
	}

	// handle the websocket handshake
	conn, err := wsh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to open websocket connection", 401)
		return
	}

	wsh.ctx.Notifier.AddClient(conn, sessionState.User.ID)
}
