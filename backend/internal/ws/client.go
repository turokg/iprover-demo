package ws

import (
	"backend/internal"
	"bytes"
	"context"
	"log"
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

//type Client interface {
//	Start()
//}

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type Client struct {
	conn   *websocket.Conn
	cancel context.CancelFunc
	wg     *sync.WaitGroup
	launch chan []byte
	logger internal.Logger
}

func checkOrigin(r *http.Request) bool {
	// TODO check origin properly
	return true
}
func NewClient(
	w http.ResponseWriter,
	r *http.Request,
	launch chan []byte,
	logger internal.Logger,
) (*Client, error) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	upgrader.CheckOrigin = checkOrigin
	conn, err := upgrader.Upgrade(w, r, nil)

	return &Client{
		conn:   conn,
		wg:     &sync.WaitGroup{},
		launch: launch,
		logger: logger,
	}, err
}

func (c *Client) Close() {
	c.cancel()
	c.wg.Wait()
	close(c.launch)
	c.conn.Close()
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
				log.Printf("error: %v", err)
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
		ticker.Stop()
	}()
	for {
		select {
		case message, ok := <-c.launch:
			//c.logger.WithField("message", string(message)).Info(c.ctx, "sending message")
			err := c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				c.logger.Error(ctx, "error while setting deadline", err)
			}
			if !ok {
				// The hub closed the channel.
				err := c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					c.logger.Error(ctx, "error while writing message", err)
				}
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				c.logger.Error(ctx, "error while getting writer", err)
				return
			}
			_, err = w.Write(message)
			if err != nil {
				c.logger.Error(ctx, "error while writing message", err)
			}

			// Add queued chat messages to the current websocket message.
			n := len(c.launch)
			for i := 0; i < n; i++ {
				_, err = w.Write(newline)
				if err != nil {
					c.logger.Error(ctx, "error while writing newline", err)
				}
				_, err = w.Write(<-c.launch)
				if err != nil {
					c.logger.Error(ctx, "error while writing message", err)
				}
			}

			if err := w.Close(); err != nil {
				c.logger.Error(ctx, "error closing websocket connection", err)
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
