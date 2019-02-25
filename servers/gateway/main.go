package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/final-project-kool-kids/servers/finalgateway/handlers"
	"github.com/final-project-kool-kids/servers/finalgateway/models/users"
	"github.com/final-project-kool-kids/servers/finalgateway/sessions"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
)

const headerUser = "X-User"

//NewServiceProxy returns a new ReverseProxy
//for a microservice given a comma-delimited
//list of network addresses
func NewServiceProxy(addrs string, ctx handlers.HandlerContext) *httputil.ReverseProxy {
	splitAddrs := strings.Split(addrs, ",")
	nextAddr := 0
	mx := sync.Mutex{}

	return &httputil.ReverseProxy{
		Director: func(r *http.Request) {
			r.URL.Scheme = "http"
			mx.Lock()
			r.URL.Host = splitAddrs[nextAddr]
			nextAddr = (nextAddr + 1) % len(splitAddrs)
			mx.Unlock()

			r.Header.Del(headerUser)
			currentState := &handlers.SessionState{}
			_, err := sessions.GetState(r, ctx.Key, ctx.SessionStore, currentState)

			if err != nil {
				return
			}
			userJSON, _ := json.Marshal(currentState.User)
			r.Header.Set(headerUser, string(userJSON))
		},
	}
}

func main() {

	sessionKey := "sessionkey"

	// gateway should listen on port 443
	addr := os.Getenv("ADDR")
	if len(addr) == 0 {
		addr = ":443"
	}

	// establishing connection with redis
	redisAddr := os.Getenv("REDISADDR")
	if len(redisAddr) == 0 {
		redisAddr = "redisServer:6379"
	}

	redisDb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	myRedisDB := sessions.NewRedisStore(redisDb, time.Hour)

	// establishing connection with sqldb
	dsn := fmt.Sprintf(os.Getenv("DSN"), os.Getenv("MYSQL_ROOT_PASSWORD"))
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("error opening database: %v\n", err)
		os.Exit(1)
	}
	mySQLDB := users.NewMySQLStore(db)
	defer db.Close()

	mqAddr := os.Getenv("MQADDR")
	if len(mqAddr) == 0 {
		mqAddr = "rabbitmq:5672"
	}

	mqName := os.Getenv("MQNAME")
	if len(mqName) == 0 {
		mqName = "rabbitmq"
	}

	// get the value of TLSCERT and TLSKEY from environment
	tlsCertPath := os.Getenv("TLSCERT")
	tlsKeyPath := os.Getenv("TLSKEY")

	mux := http.NewServeMux()

	notifier := handlers.NewNotifier()

	handlersContext := handlers.HandlerContext{
		Key:          sessionKey,
		SessionStore: myRedisDB,
		UserStore:    mySQLDB,
		Notifier:     notifier,
	}

	courseServiceAddrs := os.Getenv("COURSESADDR")
	if len(courseServiceAddrs) == 0 {
		courseServiceAddrs = ":80"
	}

	faqServiceAddrs := os.Getenv("FAQADDR")
	if len(faqServiceAddrs) == 0 {
		faqServiceAddrs = ":80"
	}

	courseProxy := NewServiceProxy(courseServiceAddrs, handlersContext)
	faqProxy := NewServiceProxy(faqServiceAddrs, handlersContext)

	handlersContext.StartMQ(mqAddr, mqName)

	mux.HandleFunc("/final/users", handlersContext.UsersHandler)
	mux.HandleFunc("/final/users/", handlersContext.SpecificUserHandler)
	mux.HandleFunc("/final/sessions", handlersContext.SessionsHandler)
	mux.HandleFunc("/final/sessions/", handlersContext.SpecificSessionHandler)
	mux.HandleFunc("/final/allusers", handlersContext.GetAllUsersHandler)

	mux.Handle("/final/courses", courseProxy)
	mux.Handle("/final/courses/", courseProxy)
	mux.Handle("/final/faq", faqProxy)
	mux.Handle("/final/faq/", faqProxy)
	mux.Handle("/final/ws", handlers.NewWebSocketsHandler(handlersContext))

	wrappedMux := handlers.NewCors(mux)

	// initiating the server
	log.Printf("server is listening on %s", addr)
	log.Fatal(http.ListenAndServeTLS(addr, tlsCertPath, tlsKeyPath, wrappedMux))

}
