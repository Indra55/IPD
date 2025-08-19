package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"github.com/Omkardalvi01/IPD/networking"
	"github.com/pion/webrtc/v3"
)

const(
	Role string = "E"
)

func main(){
	var dir_name string
	fmt.Println("Name of dir you want to copy into:")
	fmt.Scan(&dir_name)
	
	err := os.MkdirAll(dir_name, 0755)
	if err != nil{
		log.Fatal("Error while make dir")
	}

	conn , err := networking.Createconnection()
	if err != nil{
		return
	} 
	
	var uid string
	fmt.Print("Give the unique_id: ")
	fmt.Scan(&uid)
	
	err = networking.Forward(conn, Role)
	if err != nil{
		log.Fatal("Error while forwarding role",err)
	}

	err = networking.Forward(conn, uid)
	if err != nil{
		log.Fatal("Error while forwarding uid",err)
	}

	pc , err := webrtc.NewPeerConnection(networking.Webconfig)
	if err != nil{
		log.Fatal("Error while intializing peer connectrion", err)
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
				file_path := filepath.Join(dir_name,file_name)
				f , err:= os.Create(file_path)

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
		log.Fatal("Error while recieveing answer",err)
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

	<-webrtc.GatheringCompletePromise(pc)

    finalAnswer := pc.LocalDescription()

    fmt.Print(finalAnswer.SDP)
	err = networking.Forward(conn, finalAnswer.SDP)
	if err != nil{
		log.Fatal("Error while forwarding answer",err)
	}

	select{}

}