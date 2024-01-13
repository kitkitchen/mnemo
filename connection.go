package store

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type (
	Connection struct {
		websocket *websocket.Conn
		Pool      *Pool
		Key       interface{}
		Messages  chan interface{}
	}
)

func NewConnection(ctx echo.Context) (*Connection, error) {
	upgrader := websocket.Upgrader{}
	websocket, err := upgrader.Upgrade(ctx.Response(), ctx.Request(), nil)
	if err != nil {
		logger.Error("error upgrading connection: ", err.Error())
		return nil, echo.NewHTTPError(http.StatusInternalServerError)
	}

	c := &Connection{
		websocket: websocket,
		Key:       uuid.New(),
		Messages:  make(chan interface{}, 16),
	}
	return c, nil
}

func (c *Connection) Close() error {
	if c == nil {
		return fmt.Errorf("attemped to close nil Connection")
	}
	c.Pool.removeConnection(c)
	c.websocket.Close()
	return nil
}

func (c *Connection) Listen() {
	go func(c *Connection) {
		for {
			if _, _, err := c.websocket.ReadMessage(); err != nil {
				if websocket.IsUnexpectedCloseError(
					err,
					websocket.CloseGoingAway,
					websocket.CloseAbnormalClosure,
					websocket.CloseNormalClosure,
				) {
					logger.Warn(err)
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
			logger.Error(err)
			c.Close()
		}
	}
}

func (c *Connection) Publish(msg interface{}) {
	c.Messages <- msg
}
