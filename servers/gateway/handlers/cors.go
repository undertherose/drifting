package handlers

import (
	"net/http"
)

//Cors defines the structure for the CORS request
type Cors struct {
	Handler http.Handler
}

func (l *Cors) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, PATCH, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Expose-Headers", "Authorization")
	w.Header().Set("Access-Control-Max-Age", "600")

	//if this is preflight request, the method will
	//be OPTIONS, so call the real handler only if
	//the method is something else
	if r.Method != http.MethodOptions {
		l.Handler.ServeHTTP(w, r)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

//NewCors return a CORS wrapped response
func NewCors(handlerToWrap http.Handler) http.Handler {
	return &Cors{handlerToWrap}
}
