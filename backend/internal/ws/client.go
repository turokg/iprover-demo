package ws

import (
	"backend/internal"
	"bytes"
	"context"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	maxMessageSize = 512

	pingPeriod = (pongWait * 9) / 10
)

type Client interface {
	Start()
}

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type client struct {
	conn   *websocket.Conn
	ctx    context.Context
	cancel context.CancelFunc
	wg     *sync.WaitGroup
	launch chan []byte
	logger internal.Logger
}

func NewClient(ctx context.Context, conn *websocket.Conn, launch chan []byte, logger internal.Logger, cancelFunc context.CancelFunc) Client {
	return &client{
		conn:   conn,
		ctx:    ctx,
		cancel: cancelFunc,
		wg:     &sync.WaitGroup{},
		launch: launch,
		logger: logger,
	}
}

func (c *client) Start() {
	go c.readPump()
	go c.writePump()
}

func (c *client) Close() {
	c.cancel()
	c.wg.Wait()
	close(c.launch)
	c.conn.Close()
}

// readPump pumps messages from the websocket connection to the hub.
func (c *client) readPump() {
	defer c.cancel()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		if bytes.Equal(message, []byte(internal.StopWord)) {
			return
		}
	}
}

func (c *client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		case message, ok := <-c.launch:
			c.logger.WithField("message", string(message)).Info(c.ctx, "sending message")
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.launch)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.launch)
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
