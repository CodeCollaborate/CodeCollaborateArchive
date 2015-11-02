package models
import "github.com/gorilla/websocket"

type WSConnection struct {
	Username string
	WSConn   *websocket.Conn
}