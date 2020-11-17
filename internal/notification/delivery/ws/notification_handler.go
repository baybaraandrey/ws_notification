package ws

import (
	"encoding/json"
	"fmt"
	"net/http"
	_ "path"
	"strconv"

	notificationUsecases "github.com/baybaraandrey/ws_notification/internal/notification/usecases"

	auth "github.com/baybaraandrey/ws_notification/pkg/auth"
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
	jwtsecretkey string,
) {
	handler := notificationHandler{notificationWs, jwtsecretkey}

	r.HandleFunc("/notifications/", handler.handle).Methods("GET")
}

// notificationHandler represent API service for notifications
type notificationHandler struct {
	notificationWs notificationUsecases.NotificationUsecase
	jwtsecretkey   string
}

var (
	upgrader = websocket.Upgrader{}
)

type Message struct {
	Msg string `json:"msg"`
}

func MessageResponse(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(&Message{msg})
}

// @Summary notifications api
// @Description notifications api jwt token must be sended after websocet connection established
// @Accept  json
// @Produce  json
// @Param query body auth.JWTAuth true "jwt token"
// @Success 200 {object} Message
// @Failure 400 {object} Message
// @Router /ws/v1/notifications/ [get]
// @Tags ws-notifications
func (h *notificationHandler) handle(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("WTF @handleWs upgrade ws", err)
		return
	}

	msg := auth.JWTAuth{}
	err = c.ReadJSON(&msg)
	if err != nil {
		MessageResponse(w, http.StatusBadRequest, err.Error())
		c.WriteJSON(&Message{err.Error()})
		c.Close()
		return
	}

	uid, err := auth.ValidateGetUIDJWT(h.jwtsecretkey, msg.JWTToken)
	if err != nil {
		c.WriteJSON(&Message{err.Error()})
		c.Close()
		return
	}

	go notificationUsecases.ServeWs(c, strconv.Itoa(uid))
}
