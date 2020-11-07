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
	wsNotificator notificationUsecases.NotificatorUsecase,
) {
	handler := notificationHandler{wsNotificator}

	r.HandleFunc("/notifications/", handler.handle).Methods("GET")
}

// notificationHandler represent API service for notifications
type notificationHandler struct {
	wsNotificator notificationUsecases.NotificatorUsecase
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
	fmt.Println("Websocket connection established")

	msg := make(map[string]interface{})
	err = c.ReadJSON(&msg)
	if err != nil {
		fmt.Println("Websocket connection closed", err)
		c.Close()
		return
	}

	userID, ok := msg["id"].(string)
	if !ok {
		fmt.Println("ws.handle error when casting userID type")
		c.Close()
		return
	}

	notificationUsecases.ServeWs(c, userID)
}
