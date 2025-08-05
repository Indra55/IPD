package networking

import (
	"log"

	"github.com/gorilla/websocket"
)


func Createconnection() (*websocket.Conn, error){
	conn, _ , err := websocket.DefaultDialer.Dial("wss://sdp-server-poak.onrender.com/", nil)
	if err != nil{
		log.Println("Error at line 11 in signaling.go")
		return nil , err
	}

	return conn, nil
}

func Forward(conn *websocket.Conn, msg string) error{
	err := conn.WriteMessage(websocket.TextMessage, []byte(msg))
	if err != nil{
		log.Println("Error at line 20 in signaling.go")
		return err
	}

	return nil
}

func Recieve(conn *websocket.Conn) (string,error){

	_ , resp , err := conn.ReadMessage()
	if err != nil {
		log.Println("Error at line 30 in signaling.go")
		return "",err
	}
	return string(resp),nil
}