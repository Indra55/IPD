package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/Omkardalvi01/IPD/networking"
	"github.com/pion/webrtc/v3"
)

func main(){


	conn , err := networking.Createconnection()
	if err != nil{
		return
	} 
	
	var uid string
	fmt.Print("Give the unique_id: ")
	fmt.Scan(&uid)

	err = networking.Forward(conn, uid)
	if err != nil{
		log.Fatal("Error at line 23",err)
	}

	pc , err := webrtc.NewPeerConnection(networking.Webconfig)
	if err != nil{
		log.Fatal("Error at peer connection at line 15", err)
	}
	
	var file_name string

	pc.OnDataChannel(func(dc *webrtc.DataChannel) {
		fmt.Printf("New DataChannel %s\n", dc.Label())

		dc.OnOpen(func() {
			fmt.Println("Connected to peer. Type messages:")

			go func() {
				var msg string
				for {
					fmt.Scan(&msg)
					dc.SendText(msg)
				}
			}()
		})

		dc.OnMessage(func(msg webrtc.DataChannelMessage) {

			if msg.IsString{
				file_name = string(msg.Data)
			}else{
				f , err:= os.Create(file_name)

				if err != nil{
					log.Fatal("Error while creating file", err)
				}

				io.Copy(f, bytes.NewBuffer(msg.Data))

				f.Close()
			}
		
		})
	})

	offer , err := networking.Recieve(conn)
	if err != nil{
		log.Fatal("Error at line 39",err)
	}

	offer_SDP := webrtc.SessionDescription{
		SDP: offer,
		Type: webrtc.SDPTypeOffer,
	}

	err = pc.SetRemoteDescription(offer_SDP)
	if err != nil{
		log.Fatal("Error at setting remote description", err)
	}

	
	answer , err := pc.CreateAnswer(nil)
	if err != nil{
		log.Fatal("Error at creating answer")
	}

	err = pc.SetLocalDescription(answer)
	if err != nil{
		log.Fatal("Error at setting local description")
	}

	fmt.Print(answer.SDP)
	err = networking.Forward(conn, answer.SDP)
	if err != nil{
		log.Fatal("Error at line 63",err)
	}

	select{}

}