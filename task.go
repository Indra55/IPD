package main

import (
	"log"

	"github.com/pion/webrtc/v3"
)

func create_data_channel(pc *webrtc.PeerConnection) (*webrtc.DataChannel, error){
	dc , err := pc.CreateDataChannel("data", nil)
	if err != nil{
		log.Print("Error at send_data line 10")
		return nil ,err
	}
	
	return dc , err
}
	