package ws

import (
	"fmt"
	"net/http"
	_ "path"

	notificationUsecases "github.com/baybaraandrey/ws_notification/internal/notification/usecases"

	log "github.com/baybaraandrey/ws_notification/pkg/log"

	// "github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var (
	logger = log.New()
)

// NewNotificationHandler register server routes
func NewNotificationHandler(
	r *mux.Router,
	notificationWs notificationUsecases.NotificationUsecase,
) {
	handler := notificationHandler{notificationWs}

	r.HandleFunc("/notifications/", handler.handle).Methods("GET")
}

// notificationHandler represent API service for notifications
type notificationHandler struct {
	notificationWs notificationUsecases.NotificationUsecase
}

var (
	upgrader = websocket.Upgrader{}
)

// @Summary notifications api
// @Description notifications api
// @Router /ws/v1/notifications/ [get]
// @Tags ws-notifications
func (h *notificationHandler) handle(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("WTF @handleWs upgrade ws", err)
		return
	}

	msg := make(map[string]interface{})
	err = c.ReadJSON(&msg)
	if err != nil {
		fmt.Println("Websocket connection closed", err)
		c.Close()
		return
	}

	if id, ok := msg["id"].(string); ok {
		go notificationUsecases.ServeWs(c, id)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("cannot cast id"))
		c.Close()
		return
	}
}
