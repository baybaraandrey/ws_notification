package main

import (
	"fmt"
	nativeLog "log"
	"net/http"
	"os"
	"sync"
	"time"

	_ "github.com/baybaraandrey/ws_notification/docs"
	httpSwagger "github.com/swaggo/http-swagger"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	config "github.com/baybaraandrey/ws_notification/internal/config"
	monitor "github.com/baybaraandrey/ws_notification/internal/monitoring/delivery/rest"
	notificationsRest "github.com/baybaraandrey/ws_notification/internal/notification/delivery/rest"
	notificationsWs "github.com/baybaraandrey/ws_notification/internal/notification/delivery/ws"
	notificationUsecases "github.com/baybaraandrey/ws_notification/internal/notification/usecases"
)

// Version indicates the current version of the application.
var Version = "v0.0.1"
var (
	version = kingpin.Flag("version", "show version").Short('v').Bool()
)

func createServer(addr string, router *mux.Router) *http.Server {
	server := &http.Server{
		Handler: handlers.CORS(handlers.AllowedOrigins([]string{"*"}))(router),
		Addr:    addr,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	return server
}

// @title Swagger Cyberjin notifications API
// @version 0.0.1
// @description This is a ws notifications app

// @contact.name Andrey Baybara
// @contact.url
// @contact.email baybaraandrey@gmail.com

// @host localhost:8080
// @BasePath /
func main() {
	wg := new(sync.WaitGroup)

	kingpin.Parse()
	if *version {
		fmt.Println(Version)
		os.Exit(0)
	}

	// load application configuration
	cfg, err := config.Load()
	if err != nil {
		nativeLog.Fatal(err)
	}

	restRouter := mux.NewRouter()
	wsRouter := mux.NewRouter()

	restv1 := restRouter.PathPrefix("/api/v1").Subrouter()
	wsv1 := wsRouter.PathPrefix("/ws/v1").Subrouter()

	monitor.NewMonitorHandler(restRouter)

	restRouter.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://0.0.0.0:%d/swagger/doc.json", cfg.RESTServerPort)), //The url pointing to API definition
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("#swagger-ui"),
	))

	restRouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/", http.StatusMovedPermanently)
	})

	notification := notificationUsecases.NewWebsocketNotification()

	jwtsecretkey := cfg.JWTSecretKey

	// API
	notificationsRest.NewNotificationHandler(restv1, notification)
	notificationsWs.NewNotificationHandler(wsv1, notification, jwtsecretkey)

	restAddr := fmt.Sprintf("0.0.0.0:%d", cfg.RESTServerPort)
	wsAddr := fmt.Sprintf("0.0.0.0:%d", cfg.WsServerPort)

	wg.Add(2)
	go func() {
		nativeLog.Println("REST API server serving on", restAddr)
		server := createServer(restAddr, restRouter)
		nativeLog.Fatal(server.ListenAndServe())
		wg.Done()
	}()
	go func() {
		nativeLog.Println("WS server serving on", wsAddr)
		server := createServer(wsAddr, wsRouter)
		nativeLog.Fatal(server.ListenAndServe())
		wg.Done()
	}()

	wg.Wait()
}
