package ws

import (
	"backend/internal/conf"
	"context"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type client struct {
	conn   *websocket.Conn
	send   chan []byte
	ctx    context.Context
	cancel context.CancelFunc
	wg     *sync.WaitGroup
}

type Client interface {
	Start(inputParams string)
}

func NewClient(ctx context.Context, conn *websocket.Conn) Client {
	ctx, cancel := context.WithTimeout(ctx, conf.RunTimeout)
	return &client{
		conn:   conn,
		send:   make(chan []byte, 256),
		ctx:    ctx,
		cancel: cancel,
		wg:     &sync.WaitGroup{},
	}
}

func (c *client) Start(inputParams string) {
	go c.readPump()
	go c.writePump()
	go c.streamMessages(inputParams)
}

func (c *client) Close() {
	c.cancel()
	c.wg.Wait()
	close(c.send)
	c.conn.Close()
}
