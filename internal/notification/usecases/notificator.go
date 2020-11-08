package usecases

import (
	"runtime"
	"time"

	log "github.com/baybaraandrey/ws_notification/pkg/log"
	"github.com/gorilla/websocket"
)

var logger = log.New()

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 10 * time.Second

	// send pings to peer with this period. Must be less than pongWait.
	pingPeriod = 5 * time.Second

	// Maximum message size allowed from peer.
	maxMessageSize = 4096 * 4
)

type NotificatorUsecase interface {
	SendAll(b []byte)
	SendDirect(directMessage *DirectMessage)
}

var wshub = &hub{
	clients:    make(map[string][]*Client, 100),
	broadcast:  make(chan []byte),
	register:   make(chan *Client, 256),
	unregister: make(chan *Client, 256),
	direct:     make(chan *DirectMessage, 256),
}

func getHub() *hub {
	return wshub
}
func NewWebsocketNotificator() NotificatorUsecase {
	return wshub
}

//Websocket Connections
type Client struct {
	id string

	hub  *hub
	send chan []byte
	conn *websocket.Conn
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		msgType, _, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Debugf("IsUnexpectedCloseError: %v\n", err)
			}
			return
		}

		if msgType == websocket.CloseMessage {
			logger.Debug("readPump close message receive")
			return
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
		logger.Debug("@Client.writePump: unregister before")
		c.hub.unregister <- c
		logger.Debug("@Client.writePump: unregister after")
	}()
	for {
		select {
		case msgBytes, ok := <-c.send:
			logger.Debug("@Client.writePump: receive message from c.send: ", ok)
			// log.Debug("writePump receive message: ", string(msgBytes))

			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel
				logger.Debug("@Client.writePump: !ok the hub closed the channel")
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			_, err = w.Write(msgBytes)
			if err != nil {
				return
			}

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				m, ok := <-c.send
				if !ok {
					logger.Debug("@Client.writePump: inside loop hub closed channel")
					return
				}
				_, err := w.Write(m)
				if err != nil {
					logger.Debug("@Client.writePump: error w.Write(<-c.send) drain", err)
					for j := i; j < n; j++ {
						<-c.send
						return
					}
				}
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

type DirectMessage struct {
	UserIDs []string
	Data    []byte
}

type hub struct {
	// Registered clients.
	clients map[string][]*Client

	// All Inbound messages for the clients.
	broadcast chan []byte

	// Direct Inbound messages for the clients
	direct chan *DirectMessage

	// register requests from the clients.
	register chan *Client

	// unregister requests from clients.
	unregister chan *Client
}

func (h *hub) SendAll(b []byte) {
	h.broadcast <- b
}

func (h *hub) SendDirect(directMessage *DirectMessage) {
	h.direct <- directMessage
}

func (h *hub) run() {
	for {
		select {
		case client := <-h.register:
			if _, ok := h.clients[client.id]; !ok {
				h.clients[client.id] = make([]*Client, 0)
			}
			h.clients[client.id] = append(h.clients[client.id], client)
			logger.Debug("@hub.run: clients add len", len(h.clients[client.id]))
		case client := <-h.unregister:
			logger.Debug("@hub.run: unregister client")
			logger.Debug("@hub.run: unregister len(client.send)", len(client.send))
			logger.Debug("@hub.run: unregister len(h.unregister)", len(h.unregister))
			clients, ok := h.clients[client.id]
			if ok {
			Loop:
				for index, c := range clients {
					if c.conn == client.conn {
						h.clients[client.id] = append(clients[:index], clients[index+1:]...)
						close(client.send)
						break Loop
					}
				}
				logger.Debug("@hub.run: clients delete len", len(h.clients[client.id]))
				if len(h.clients[client.id]) == 0 {
					delete(h.clients, client.id)
				}
				logger.Debug("@hub.run: len(@hub.clients):", len(h.clients))
			}
			n := len(h.unregister)
			for i := 0; i < n; i++ {
				client := <-h.unregister
				clients, ok := h.clients[client.id]
				if ok {
				Loop1:
					for index, c := range clients {
						if c.conn == client.conn {
							h.clients[client.id] = append(clients[:index], clients[index+1:]...)
							close(client.send)
							break Loop1
						}
					}
					logger.Debug("@hub.run: clients delete len", len(h.clients[client.id]))
					if len(h.clients[client.id]) == 0 {
						delete(h.clients, client.id)
					}
					logger.Debug("@hub.run: len(@hub.clients):", len(h.clients))
				}
			}
		case message := <-h.broadcast:
			for _, connList := range h.clients {
				for index, c := range connList {
					_ = index
					c.send <- message
				}
			}
		case directMessage := <-h.direct:
			for _, userID := range directMessage.UserIDs {
				clients, ok := h.clients[userID]

				if ok {
					logger.Debug("@hub.direct: sendto", userID)
					for _, c := range clients {
						select {
						case c.send <- directMessage.Data:
						default:
							logger.Debug("@hub.direct: missing message len(c.send)", len(c.send))
						}
					}
				}
			}
		}
	}
}

func ServeWs(conn *websocket.Conn, userID string) {
	client := &Client{
		id:   userID,
		hub:  getHub(),
		conn: conn,
		send: make(chan []byte, 256),
	}

	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}

func init() {
	go wshub.run()
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for range ticker.C {
			logger.Debug("Num goroutine:", runtime.NumGoroutine())
		}
	}()
}
