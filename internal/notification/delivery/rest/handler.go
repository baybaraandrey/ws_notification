package rest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	nativeLog "log"
	"net/http"
	_ "path"

	NotificationUsecases "github.com/baybaraandrey/ws_notification/internal/notification/usecases"
	log "github.com/baybaraandrey/ws_notification/pkg/log"

	"github.com/gorilla/mux"
)

var (
	logger = log.New()
)

// NewNotificationHandler register server routes
func NewNotificationHandler(
	r *mux.Router,
	wsNotificator NotificationUsecases.NotificatorUsecase,
) {
	handler := notificationHandler{wsNotificator}

	r.HandleFunc("/notifications/", handler.handle).Methods("GET", "POST")
}

type message struct {
	TransMap struct {
		UserIDs []string `json:"user_ids"`
	} `json:"trans_map"`
}

// notificationHandler represent API service for notifications
type notificationHandler struct {
	wsNotificator NotificationUsecases.NotificatorUsecase
}

var i = 0

// @Summary notify client
// @Description notify client
// @Router /api/v1/notifications/ [post]
// @Tags rest-notifications
func (h *notificationHandler) handle(w http.ResponseWriter, r *http.Request) {
	msg := &message{}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		nativeLog.Fatal(err)
	}

	err = json.Unmarshal(b, &msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Printf("@rest.handle %+v\n", *msg)

	if len(msg.TransMap.UserIDs) >= 0 {
		h.wsNotificator.SendDirect(&NotificationUsecases.DirectMessage{
			UserIDs: msg.TransMap.UserIDs,
			Data:    b,
		})
	} else {
		h.wsNotificator.SendAll(b)
	}
	i++
	fmt.Println("@rest.handle sended:", i)
	// Do something with the Person struct...
	w.Write(b)
}
