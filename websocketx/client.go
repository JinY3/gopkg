package websocketx

import (
	"net/url"

	"github.com/gorilla/websocket"
)

func Client(url url.URL) (*websocket.Conn, error) {
	c, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		return nil, err
	}
	return c, nil
}
