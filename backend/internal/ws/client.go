package ws

import (
	"backend/internal"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
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

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type Client struct {
	conn   *websocket.Conn
	cancel context.CancelFunc
	//wg       *sync.WaitGroup
	messages chan internal.LogMessage
	logger   internal.Logger
}

func checkOrigin(r *http.Request) bool {
	// TODO check origin properly
	return true
}
func NewClient(
	w http.ResponseWriter,
	r *http.Request,
	messages chan internal.LogMessage,
	logger internal.Logger,
) (*Client, error) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	upgrader.CheckOrigin = checkOrigin
	conn, err := upgrader.Upgrade(w, r, nil)

	return &Client{
		conn:     conn,
		messages: messages,
		logger:   logger,
	}, err
}

func (c *Client) Read(ctx context.Context, wg *sync.WaitGroup, cancelFunc context.CancelFunc) {
	wg.Add(1)
	defer wg.Done()

	c.conn.SetReadLimit(maxMessageSize)
	err := c.conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		c.logger.Error(ctx, "error while setting deadline", err)
	}
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.logger.Error(ctx, "error while reading from websocket", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		if bytes.Equal(message, []byte(internal.StopWord)) {
			cancelFunc()
			return
		}
	}
}

func (c *Client) Write(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	ticker := time.NewTicker(pingPeriod)
	defer func() {
		if err := c.conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
			c.logger.Error(ctx, "error while writing message", err)
		}
		if err := c.conn.Close(); err != nil {
			c.logger.Error(ctx, "error closing connection", err)
		}
		ticker.Stop()
	}()

	for {
		select {
		case message, ok := <-c.messages:
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				c.logger.Error(ctx, "error while setting deadline", err)
			}
			if !ok {
				// channel is closed
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				c.logger.Error(ctx, "error while getting writer", err)
				return
			}

			bytesMsg, _ := json.Marshal(message)
			if _, err = w.Write(bytesMsg); err != nil {
				c.logger.Error(ctx, "error while writing message", err)
			}

			// Add queued chat messages to the current websocket message.
			n := len(c.messages)
			for i := 0; i < n; i++ {
				bytesMsg, _ = json.Marshal(c.messages)
				if _, err = w.Write(bytesMsg); err != nil {
					c.logger.Error(ctx, "error while writing message", err)
				}
			}

			if err := w.Close(); err != nil {
				c.logger.Error(ctx, "error closing websocket connection", err)
				return
			}
		case <-ticker.C:
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				c.logger.Error(ctx, "error while setting deadline", err)
			}
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
