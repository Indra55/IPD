package networking

import (
	"log"
	"github.com/pion/webrtc/v3"

)


func Peerconnection(uid	string) (*webrtc.PeerConnection, error) {

	conn, err := Createconnection()
	if err != nil{
		return nil, err
	}
	defer conn.Close()

	err = Forward(conn, uid)
	if err != nil{
		return nil, err
	}
	
	peer_conn, err := webrtc.NewPeerConnection(Webconfig)
	if err != nil {
		log.Println("error at line 12 of user_comms.go")
		return nil, err 
	}
	
	offer , err := peer_conn.CreateOffer(nil)
	if err != nil {
		log.Println("error at line 19 of user_comms.go")
		return nil, err 
	}

	err = peer_conn.SetLocalDescription(offer)
	if err != nil {
		log.Println("error at line 25 of user_comms.go")
		return nil, err 
	}

	err = Forward(conn, offer.SDP)
	if err != nil{
		return nil, err
	}
	
	resp , err := Recieve(conn)
	if err != nil{
		return nil, err
	}

	answer_sdp := webrtc.SessionDescription{
		SDP: resp,
		Type: webrtc.SDPTypeAnswer,
	}

	err = peer_conn.SetRemoteDescription(answer_sdp)
	if err != nil {
		log.Println("error at line 25 of user_comms.go")
		return nil, err  
	}

	return peer_conn , nil

}