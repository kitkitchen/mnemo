package mnemo

import (
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type (
	Conn struct {
		websocket *websocket.Conn
		Pool      *Pool
		Key       interface{}
		Messages  chan interface{}
	}
)

func NewConnection(ctx echo.Context) (*Conn, error) {
	upgrader := websocket.Upgrader{}
	websocket, err := upgrader.Upgrade(ctx.Response(), ctx.Request(), nil)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError)
	}

	c := &Conn{
		websocket: websocket,
		Key:       uuid.New(),
		Messages:  make(chan interface{}, 16),
	}
	return c, nil
}

func (c *Conn) Close() error {
	if c == nil {
		return fmt.Errorf("attemped to close nil Connection")
	}
	c.Pool.removeConnection(c)
	c.websocket.Close()
	return nil
}

func (c *Conn) Listen() {
	go func(c *Conn) {
		for {
			if _, _, err := c.websocket.ReadMessage(); err != nil {
				if websocket.IsUnexpectedCloseError(
					err,
					websocket.CloseGoingAway,
					websocket.CloseAbnormalClosure,
					websocket.CloseNormalClosure,
				) {
					log.Println(err)
				}
				close(c.Messages)
				break
			}
		}
	}(c)

	for {
		msg, ok := <-c.Messages
		if !ok {
			c.Close()
			break
		}

		if err := c.websocket.WriteJSON(msg); err != nil {
			log.Println(err)
			c.Close()
		}
	}
}

func (c *Conn) Publish(msg interface{}) {
	c.Messages <- msg
}