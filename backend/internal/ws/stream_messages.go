package ws

import (
	"backend/internal/launcher"
	"fmt"
)

func (c *client) streamMessages(inputParams string) {
	c.wg.Add(1)
	msgs, _ := launcher.Launch(c.ctx, inputParams)
	for msg := range msgs {
		fmt.Println("recieved in channel", string(msg))
		c.send <- msg
	}
	c.wg.Done()
	c.Close()
}
