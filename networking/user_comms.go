package networking

import (
	"fmt"
	"log"

	"github.com/pion/webrtc/v3"
)


func Peerconnection(uid	string) (*webrtc.PeerConnection, *webrtc.DataChannel, error) {

	conn, err := Createconnection()
	if err != nil{
		return nil, nil, err
	}
	defer conn.Close()

	err = Forward(conn, uid)
	if err != nil{
		return nil, nil,  err
	}
	
	peer_conn, err := webrtc.NewPeerConnection(Webconfig)
	if err != nil {
		log.Println("error at line 26 of user_comms.go")
		return nil, nil, err 
	}
	
	dc, err := peer_conn.CreateDataChannel("data", nil)
	if err != nil{
		log.Println("error at line 30 of user_comms.go")
		return nil, nil, err 
	}

	offer , err := peer_conn.CreateOffer(nil)
	if err != nil {
		log.Println("error at line 38 of user_comms.go")
		return nil, nil, err 
	}


	err = peer_conn.SetLocalDescription(offer)
	if err != nil {
		log.Println("error at line 45 of user_comms.go")
		return nil, nil, err 
	}

	<-webrtc.GatheringCompletePromise(peer_conn)

	err = Forward(conn, peer_conn.LocalDescription().SDP)
	if err != nil{
		return nil, nil, err
	}
	
	resp , err := Recieve(conn)
	if err != nil{
		return nil, nil, err
	}

	answer_sdp := webrtc.SessionDescription{
		SDP: resp,
		Type: webrtc.SDPTypeAnswer,
	}

	err = peer_conn.SetRemoteDescription(answer_sdp)
	if err != nil {
		log.Println("error at line 69 of user_comms.go")
		return nil, nil, err  
	}
	fmt.Print(peer_conn.RemoteDescription().SDP)

	return peer_conn , dc, nil

}