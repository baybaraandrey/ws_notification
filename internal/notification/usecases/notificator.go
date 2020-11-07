package usecases

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// send pings to peer with this period. Must be less than pongWait.
	pingPeriod = 55 * time.Second

	// Maximum message size allowed from peer.
	maxMessageSize = 4096 * 4
)

type NotificatorUsecase interface {
	SendAll(b []byte)
	SendDirect(directMessage *DirectMessage)
}

var wshub = &hub{
	clients:    make(map[string][]*Client, 0),
	broadcast:  make(chan []byte, 512),
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

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
		c.hub.unregister <- c
		fmt.Println("unregister", c)
	}()
	for {
		select {
		case msgBytes, ok := <-c.send:
			fmt.Println("@Client.writePump: <-c.send")
			// log.Debug("writePump receive message: ", string(msgBytes))

			// c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel
				fmt.Println("@Client.writePump: !ok the hub closed the channel")
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := c.conn.WriteMessage(websocket.TextMessage, msgBytes)
			if err != nil {
				fmt.Println("Error when writing", err)
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
			fmt.Println("@hub.run: clients add", h.clients)
		case client := <-h.unregister:
			fmt.Println("@hub.run: unregister")
			clients, ok := h.clients[client.id]
			if ok {
			Loop:
				for index, c := range clients {
					fmt.Println("Loop", index, c)
					if c.conn == client.conn {
						fmt.Println("Same", index, c)
						h.clients[client.id] = append(clients[:index], clients[index+1:]...)
						break Loop
					}
				}
				if len(h.clients[client.id]) == 0 {
					delete(h.clients, client.id)
				}
				close(client.send)
				fmt.Println("len(@hub.clients):", h.clients)
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
					fmt.Println("@hub.direct: sendto", userID)
					for _, c := range clients {
						c.send <- directMessage.Data
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

	fmt.Println("@ServeWs client:", userID)

	client.hub.register <- client

	go client.writePump()
}

func init() {
	fmt.Println("@notificator.init: run wshub")
	go wshub.run()
}
