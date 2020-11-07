package ws

import (
	"fmt"
	"net/http"
	_ "path"
	"strconv"

	notificationRepositories "github.com/baybaraandrey/ws_notification/internal/notification/repositories"
	notificationUsecases "github.com/baybaraandrey/ws_notification/internal/notification/usecases"

	log "github.com/baybaraandrey/ws_notification/pkg/log"
	"github.com/baybaraandrey/ws_notification/pkg/utils"

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
	userRepository notificationRepositories.UserRepository,
) {
	handler := notificationHandler{wsNotificator, userRepository}

	r.HandleFunc("/notifications/", handler.handle).Methods("GET")
}

// notificationHandler represent API service for notifications
type notificationHandler struct {
	wsNotificator  notificationUsecases.NotificatorUsecase
	userRepository notificationRepositories.UserRepository
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

	username, ok := msg["username"].(string)
	if !ok {
		fmt.Println("ws.handle error when casting username type")
		c.Close()
		return
	}
	password, ok := msg["password"].(string)
	if !ok {
		fmt.Println("ws.handle error when casting password type")
		c.Close()
		return
	}

	user, err := h.userRepository.GetUserByUsername(username)
	if err != nil {
		fmt.Println(err)
		c.Close()
		return
	}

	if ok := utils.CheckDjangoPassword(password, user.Password); !ok {
		fmt.Println("Wrong password")
		c.Close()
		return
	}

	userID := strconv.Itoa(user.ID)
	notificationUsecases.ServeWs(c, userID)
}
